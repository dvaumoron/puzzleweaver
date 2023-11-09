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
	"errors"
	"log"
	"strings"

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
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/dvaumoron/puzzleweb/common/build"
	web "github.com/dvaumoron/puzzleweb/core"
	"go.uber.org/zap"
)

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
	// TODO wrapper
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

	site, ok := web.BuildDefaultSite(globalConfig)
	if !ok {
		return errors.New("TODO")
	}

	for _, pageGroup := range globalConfig.StaticPages {
		site.AddStaticPages(pageGroup)
	}

	widgets := globalConfig.Widgets
	for _, widgetPageConfig := range globalConfig.WidgetPages {
		name := widgetPageConfig.Path
		nested := false
		var parentPage web.Page
		if index := strings.LastIndex(name, "/"); index != -1 {
			emplacement := name[:index]
			name = name[index+1:]
			parentPage, nested = site.GetPageWithPath(emplacement)
			if !nested {
				logger.Error("Failed to retrieve parentPage", zap.String("emplacement", emplacement))
				continue
			}
		}

		widgetPage, add := build.MakeWidgetPage(name, ctx, globalConfig, widgets[widgetPageConfig.WidgetRef])
		if add {
			if nested {
				parentPage.AddSubPage(widgetPage)
			} else {
				site.AddPage(widgetPage)
			}
		}
	}
	return site.RunListener(globalConfig, app.web)
}
