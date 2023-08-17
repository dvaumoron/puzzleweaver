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
	"os"
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
	"github.com/dvaumoron/puzzleweaver/web/wiki"
	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
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
// weaver.Run creates it and passes it to serve.
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
	forumService            weaver.Ref[forumservice.FullForumService]
	markdownService         weaver.Ref[service.MarkdownService]
	blogService             weaver.Ref[blogservice.BlogService]
	wikiService             weaver.Ref[wikiservice.WikiService]
}

// frameServe is called by weaver.Run and contains the body of the application.
func frameServe(ctx context.Context, app *frameApp) error {
	globalConfig := app.Config()
	site := web.BuildDefaultSite(globalConfig)
	ctxLogger := app.Logger(ctx)

	// TODO read config from weaver toml
	frameConfigBody, err := os.ReadFile(os.Getenv("FRAME_CONFIG_PATH"))
	if err != nil {
		ctxLogger.Error("Failed to read frame configuration file", common.ErrorKey, err)
	}

	var frameConfig map[string]any
	if err = yaml.Unmarshal(frameConfigBody, &frameConfig); err != nil {
		ctxLogger.Error("Failed to parse frame configuration", common.ErrorKey, err)
	}

	site.AddPage(web.MakeHiddenStaticPage(app, notFound, service.PublicGroupId, notFound))

	for _, pageGroup := range asSlice("pageGroups", frameConfig["pageGroups"], ctxLogger) {
		castedPageGroup := asMap("pageGroup", pageGroup, ctxLogger)
		site.AddStaticPages(
			app,
			asUint64("pageGroup.id", castedPageGroup["id"], ctxLogger),
			asStringSlice("pageGroup.pages", castedPageGroup["pages"], ctxLogger),
		)
	}

	widgets := asMap("widgets", frameConfig["widgets"], ctxLogger)
	for _, widgetPageConfig := range asSlice("widgetPages", frameConfig["widgetPages"], ctxLogger) {
		castedWidgetPage := asMap("widgetPage", widgetPageConfig, ctxLogger)
		emplacement := asString("widgetPage.emplacement", castedWidgetPage["emplacement"], ctxLogger)
		ok := false
		var parentPage web.Page
		if emplacement != "" {
			parentPage, ok = site.GetPageWithPath(emplacement)
			if !ok {
				ctxLogger.Error("Failed to retrive parentPage", "emplacement", emplacement)
			}
		}

		widgetPage := makeWidgetPage(
			app, asString("widgetPage.name", castedWidgetPage["name"], ctxLogger), app, ctx, globalConfig,
			widgets[asString("widgetPage.widgetRef", castedWidgetPage["widgetRef"], ctxLogger)],
		)

		if ok {
			parentPage.AddSubPage(widgetPage)
		} else {
			site.AddPage(widgetPage)
		}
	}

	siteConfig := globalConfig.ExtractSiteConfig()
	// emptying data no longer useful for GC cleaning
	globalConfig = nil

	return site.Run(siteConfig)
}

func makeWidgetPage(app *frameApp, pageName string, loggerGetter common.LoggerGetter, ctx context.Context, globalConfig *config.GlobalConfig, widgetConfig any) web.Page {
	ctxLogger := loggerGetter.Logger(ctx)
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
			pageName, ctxLogger, app.blogService.Get(), app.forumService.Get(), app.markdownService.Get(),
			globalConfig.CreateBlogConfig(blogId, groupId, args...),
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
		serviceAddr := asString("widget.serviceAddr", castedConfig["serviceAddr"], ctxLogger)
		widgetName := asString("widget.widgetName", castedConfig["widgetName"], ctxLogger)
		objectId := asUint64("widget.objectId", castedConfig["objectId"], ctxLogger)
		groupId := asUint64("widget.groupId", castedConfig["groupId"], ctxLogger)
		return remotewidget.MakeRemotePage(
			pageName, loggerGetter, ctx, widgetName, globalConfig.CreateWidgetConfig(serviceAddr, objectId, groupId),
		)
	default:
		ctxLogger.Error("Widget kind unknown ", "widgetKind", kind)
	}
	return web.Page{} // unreachable
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
