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

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web"
	"github.com/dvaumoron/puzzleweaver/web/blog"
	blogservice "github.com/dvaumoron/puzzleweaver/web/blog/service"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/dvaumoron/puzzleweaver/web/forum"
	"github.com/dvaumoron/puzzleweaver/web/remotewidget"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
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
	forumService            weaver.Ref[remoteservice.RemoteForumService]
	markdownService         weaver.Ref[service.MarkdownService]
	blogService             weaver.Ref[blogservice.BlogService]
	wikiService             weaver.Ref[wikiservice.WikiService]
	widgetService           weaver.Ref[widgetservice.WidgetService]
}

// frameServe is called by weaver.Run and contains the body of the application.
func frameServe(ctx context.Context, app *frameApp) error {
	globalConfig := &config.GlobalServiceConfig{
		LoggerGetter:            app,
		GlobalConfig:            app.Config(),
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
	for _, widgetPageConfig := range globalConfig.WidgetPages {
		ok := false
		var parentPage web.Page
		if emplacement := widgetPageConfig.Emplacement; emplacement != "" {
			parentPage, ok = site.GetPageWithPath(emplacement)
			if !ok {
				ctxLogger.Error("Failed to retrive parentPage", "emplacement", emplacement)
			}
		}

		widgetPage, add := makeWidgetPage(app, widgetPageConfig.Name, globalConfig, ctx, ctxLogger, widgets[widgetPageConfig.WidgetRef])
		if add {
			if ok {
				parentPage.AddSubPage(widgetPage)
			} else {
				site.AddPage(widgetPage)
			}
		}
	}

	return site.Run(globalConfig)
}

func makeWidgetPage(app *frameApp, pageName string, globalConfig *config.GlobalServiceConfig, ctx context.Context, ctxLogger *slog.Logger, widgetConfig config.WidgetConfig) (web.Page, bool) {
	switch widgetConfig.Kind {
	case "forum":
		return forum.MakeForumPage(pageName, ctxLogger, globalConfig.CreateForumConfig(
			widgetConfig.ObjectId, widgetConfig.GroupId, widgetConfig.Templates...,
		)), true
	case "blog":
		return blog.MakeBlogPage(pageName, ctxLogger, globalConfig.CreateBlogConfig(
			widgetConfig.ObjectId, widgetConfig.GroupId, widgetConfig.Templates...,
		)), true
	case "wiki":
		return wiki.MakeWikiPage(pageName, ctxLogger, globalConfig.CreateWikiConfig(
			widgetConfig.ObjectId, widgetConfig.GroupId, widgetConfig.Templates...,
		)), true
	case "remote":
		return remotewidget.MakeRemotePage(pageName, ctx, widgetConfig.WidgetName,
			globalConfig.CreateWidgetConfig(widgetConfig.ObjectId, widgetConfig.GroupId),
		)
	default:
		ctxLogger.Error("Widget kind unknown ", "kind", widgetConfig.Kind)
		return web.Page{}, false
	}
}
