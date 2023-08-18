/*
 *
 * Copyright 2023 puzzleweaver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	_ "embed"
	"log"
	"strconv"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web"
	"github.com/dvaumoron/puzzleweaver/web/blog"
	blogservice "github.com/dvaumoron/puzzleweaver/web/blog/service"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/dvaumoron/puzzleweaver/web/forum"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
	"github.com/dvaumoron/puzzleweaver/web/remotewidget"
	remotewidgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	"github.com/dvaumoron/puzzleweaver/web/wiki"
	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
	"golang.org/x/exp/slog"
)

const notFound = "notFound"
const castMsg = "Failed to cast value"
const valueName = "valueName"

// TODO go:embed version.txt
var version string

func main() {
	if err := weaver.Run(context.Background(), frameServe); err != nil {
		log.Fatal(err)
	}
}

// frameApp is the main component of the application.
// weaver.Run creates it and passes it to frameServe.
type frameApp struct {
	weaver.Implements[weaver.Main]
	weaver.WithConfig[config.GlobalConfig]
	sessionService          weaver.Ref[service.SessionService]
	templateService         weaver.Ref[service.TemplateService]
	settingsService         weaver.Ref[service.SettingsService]
	passwordStrengthService weaver.Ref[service.PasswordStrengthService]
	saltService             weaver.Ref[service.SaltService]
	loginService            weaver.Ref[service.FullLoginService]
	adminService            weaver.Ref[service.AdminService]
	profileService          weaver.Ref[service.AdvancedProfileService]
	forumService            weaver.Ref[forumservice.FullForumService]
	markdownService         weaver.Ref[service.MarkdownService]
	blogService             weaver.Ref[blogservice.BlogService]
	wikiService             weaver.Ref[wikiservice.WikiService]
	widgetService           weaver.Ref[remotewidgetservice.WidgetService]
}

// frameServe is called by weaver.Run and contains the body of the application.
func frameServe(ctx context.Context, app *frameApp) error {
	globalConfig := &config.GlobalServiceConfig{
		GlobalConfig: app.Config(), LoggerGetter: app,
		SessionService:          app.sessionService.Get(),
		TemplateService:         app.templateService.Get(),
		SettingsService:         app.settingsService.Get(),
		PasswordStrengthService: app.passwordStrengthService.Get(),
		SaltService:             app.saltService.Get(),
		LoginService:            app.loginService.Get(),
		AdminService:            app.adminService.Get(),
		ProfileService:          app.profileService.Get(),
		ForumService:            app.forumService.Get(),
		MarkdownService:         app.markdownService.Get(),
		BlogService:             app.blogService.Get(),
		WikiService:             app.wikiService.Get(),
		WidgetService:           app.widgetService.Get(),
	}

	ctxLogger := app.Logger(ctx)
	site := web.BuildDefaultSite(ctxLogger, globalConfig)

	site.AddPage(web.MakeHiddenStaticPage(app, notFound, service.PublicGroupId, notFound))

	for _, pageGroup := range globalConfig.PageGroups {
		site.AddStaticPages(app, pageGroup.Id, pageGroup.Pages)
	}

	widgets := globalConfig.Widgets
	for _, widgetPage := range globalConfig.WidgetPages {
		ok := false
		var parentPage web.Page
		if emplacement := widgetPage.Emplacement; emplacement != "" {
			parentPage, ok = site.GetPageWithPath(emplacement)
			if !ok {
				ctxLogger.Error("Failed to retrive parentPage", "emplacement", emplacement)
			}
		}

		widgetPage := makeWidgetPage(app, widgetPage.Name, globalConfig, ctx, widgets[widgetPage.WidgetRef])

		if ok {
			parentPage.AddSubPage(widgetPage)
		} else {
			site.AddPage(widgetPage)
		}
	}

	return site.Run(globalConfig)
}

func makeWidgetPage(app *frameApp, pageName string, globalConfig *config.GlobalServiceConfig, ctx context.Context, widgetConfig any) web.Page {
	ctxLogger := globalConfig.LoggerGetter.Logger(ctx)
	castedConfig := asMap("widget", widgetConfig, ctxLogger)

	switch kind := asString("widget.kind", castedConfig["kind"], ctxLogger); kind {
	case "forum":
		forumId := asUint64("widget.forumId", castedConfig["forumId"], ctxLogger)
		groupId := asUint64("widget.groupId", castedConfig["groupId"], ctxLogger)
		args := asStringSlice("widget.templates", castedConfig["templates"], ctxLogger)
		return forum.MakeForumPage(pageName, ctxLogger, app.forumService.Get(), globalConfig.CreateForumConfig(forumId, groupId, args...))
	case "blog":
		blogId := asUint64("widget.blogId", castedConfig["blogId"], ctxLogger)
		groupId := asUint64("widget.groupId", castedConfig["groupId"], ctxLogger)
		args := asStringSlice("widget.templates", castedConfig["templates"], ctxLogger)
		return blog.MakeBlogPage(
			pageName, ctxLogger, globalConfig.CreateBlogConfig(blogId, groupId, args...),
		)
	case "wiki":
		wikiId := asUint64("widget.wikiId", castedConfig["wikiId"], ctxLogger)
		groupId := asUint64("widget.groupId", castedConfig["groupId"], ctxLogger)
		args := asStringSlice("widget.templates", castedConfig["templates"], ctxLogger)
		return wiki.MakeWikiPage(
			pageName, ctxLogger, app.wikiService.Get(), app.markdownService.Get(),
			globalConfig.CreateWikiConfig(wikiId, groupId, args...),
		)
	case "remote":
		widgetName := asString("widget.widgetName", castedConfig["widgetName"], ctxLogger)
		objectId := asUint64("widget.objectId", castedConfig["objectId"], ctxLogger)
		groupId := asUint64("widget.groupId", castedConfig["groupId"], ctxLogger)
		return remotewidget.MakeRemotePage(
			pageName, globalConfig.LoggerGetter, ctx, widgetName, globalConfig.CreateWidgetConfig(objectId, groupId),
		)
	default:
		ctxLogger.Error("Widget kind unknown ", "widgetKind", kind)
	}
	return web.Page{} // TODO redirecter to not found page
}

func asUint64(name string, value any, ctxLogger *slog.Logger) uint64 {
	if value == nil {
		return 0
	}
	switch casted := value.(type) {
	case uint:
		return uint64(casted)
	case uint8:
		return uint64(casted)
	case uint16:
		return uint64(casted)
	case uint32:
		return uint64(casted)
	case uint64:
		return uint64(casted)
	case int:
		return uint64(casted)
	case int8:
		return uint64(casted)
	case int16:
		return uint64(casted)
	case int32:
		return uint64(casted)
	case int64:
		return uint64(casted)
	case float32:
		return uint64(casted)
	case float64:
		return uint64(casted)
	case string:
		i, err := strconv.ParseUint(casted, 10, 64)
		if err != nil {
			ctxLogger.Error("Failed to parse value", valueName, name, common.ErrorKey, err)
		}
		return i
	default:
		ctxLogger.Error(castMsg, valueName, name)
	}
	return 0 // unreachable
}

func asMap(name string, value any, ctxLogger *slog.Logger) map[string]any {
	if value == nil {
		return nil
	}
	m, ok := value.(map[string]any)
	if !ok {
		ctxLogger.Error(castMsg, "valueName", name)
	}
	return m
}

func asSlice(name string, value any, ctxLogger *slog.Logger) []any {
	if value == nil {
		return nil
	}
	s, ok := value.([]any)
	if !ok {
		ctxLogger.Error(castMsg, valueName, name)
	}
	return s
}

func asString(name string, value any, ctxLogger *slog.Logger) string {
	if value == nil {
		return ""
	}
	s, ok := value.(string)
	if !ok {
		ctxLogger.Error(castMsg, valueName, name)
	}
	return s
}

func asStringSlice(name string, value any, ctxLogger *slog.Logger) []string {
	s := asSlice(name, value, ctxLogger)
	s2 := make([]string, 0, len(s))
	for _, innerValue := range s {
		s2 = append(s2, asString(name, innerValue, ctxLogger))
	}
	return s2
}
