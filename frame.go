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
	"log/slog"

	"github.com/ServiceWeaver/weaver"
	adminimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/admin"
	blogimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/blog"
	forumimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/forum"
	loginimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/login"
	markdownimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/markdown"
	passwordstrengthimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/passwordstrength"
	profileimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/profile"
	remotewidgetimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget"
	saltimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/salt"
	sessionimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/session"
	settingsimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/settings"
	templatesimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/templates"
	wikiimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/wiki"
	"github.com/dvaumoron/puzzleweaver/web"
	"github.com/dvaumoron/puzzleweaver/web/blog"
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/dvaumoron/puzzleweaver/web/forum"
	"github.com/dvaumoron/puzzleweaver/web/remotewidget"
	"github.com/dvaumoron/puzzleweaver/web/wiki"
)

const notFound = "notFound"
const castMsg = "Failed to cast value"
const valueName = "valueName"

//go:embed version.txt
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
	web                     weaver.Listener
	sessionService          weaver.Ref[sessionimpl.SessionService]
	templateService         weaver.Ref[templatesimpl.TemplateService]
	settingsService         weaver.Ref[settingsimpl.SettingsService]
	passwordStrengthService weaver.Ref[passwordstrengthimpl.PasswordStrengthService]
	saltService             weaver.Ref[saltimpl.SaltService]
	loginService            weaver.Ref[loginimpl.RemoteLoginService]
	adminService            weaver.Ref[adminimpl.AdminService]
	profileService          weaver.Ref[profileimpl.RemoteProfileService]
	forumService            weaver.Ref[forumimpl.RemoteForumService]
	markdownService         weaver.Ref[markdownimpl.MarkdownService]
	blogService             weaver.Ref[blogimpl.RemoteBlogService]
	wikiService             weaver.Ref[wikiimpl.RemoteWikiService]
	widgetService           weaver.Ref[remotewidgetimpl.RemoteWidgetService]
}

// frameServe is called by weaver.Run and contains the body of the application.
func frameServe(ctx context.Context, app *frameApp) error {
	logger := app.Logger(ctx)

	globalConfig, err := config.New(
		app.Config(), app, logger, version, app.sessionService.Get(), app.templateService.Get(), app.settingsService.Get(),
		app.passwordStrengthService.Get(), app.saltService.Get(), app.loginService.Get(), app.adminService.Get(),
		app.profileService.Get(), app.forumService.Get(), app.markdownService.Get(), app.blogService.Get(),
		app.wikiService.Get(), app.widgetService.Get(),
	)
	if err != nil {
		return err
	}

	site := web.BuildDefaultSite(logger, globalConfig)

	site.AddPage(web.MakeHiddenStaticPage(app, notFound, adminimpl.PublicGroupId, notFound))

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
				logger.Error("Failed to retrive parentPage", "emplacement", emplacement)
			}
		}

		widgetPage, add := makeWidgetPage(app, widgetPageConfig.Name, globalConfig, ctx, logger, widgets[widgetPageConfig.WidgetRef])
		if add {
			if ok {
				parentPage.AddSubPage(widgetPage)
			} else {
				site.AddPage(widgetPage)
			}
		}
	}
	return site.Run(globalConfig, app.web)
}

func makeWidgetPage(app *frameApp, pageName string, globalConfig *config.GlobalServiceConfig, ctx context.Context, logger *slog.Logger, widgetConfig config.WidgetConfig) (web.Page, bool) {
	switch widgetConfig.Kind {
	case "forum":
		return forum.MakeForumPage(pageName, logger, globalConfig.CreateForumConfig(
			widgetConfig.ObjectId, widgetConfig.GroupId, widgetConfig.Templates...,
		)), true
	case "blog":
		return blog.MakeBlogPage(pageName, logger, globalConfig.CreateBlogConfig(
			widgetConfig.ObjectId, widgetConfig.GroupId, widgetConfig.Templates...,
		)), true
	case "wiki":
		return wiki.MakeWikiPage(pageName, logger, globalConfig.CreateWikiConfig(
			widgetConfig.ObjectId, widgetConfig.GroupId, widgetConfig.Templates...,
		)), true
	case "remote":
		return remotewidget.MakeRemotePage(pageName, ctx, logger, globalConfig.CreateWidgetConfig(
			widgetConfig.WidgetName, widgetConfig.ObjectId, widgetConfig.GroupId,
		))
	default:
		logger.Error("Widget kind unknown", "kind", widgetConfig.Kind)
		return web.Page{}, false
	}
}
