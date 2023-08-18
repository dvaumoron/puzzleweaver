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

package web

import (
	"context"
	"net/http"
	"time"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/dvaumoron/puzzleweaver/web/locale"
	"github.com/dvaumoron/puzzleweaver/web/templates"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const siteName = "Site"
const unknownUserKey = "ErrorUnknownUser"

type Site struct {
	loggerGetter   common.LoggerGetter
	localesManager locale.Manager
	authService    service.AuthService
	timeOut        time.Duration
	root           Page
	adders         []common.DataAdder
}

func NewSite(globalConfig config.GlobalServiceConfig, localesManager locale.Manager, settingsManager *SettingsManager) *Site {
	loggerGetter := globalConfig.LoggerGetter
	root := MakeStaticPage(loggerGetter, "root", service.PublicGroupId, "index")
	root.AddSubPage(newLoginPage(globalConfig.LoginService, settingsManager))
	root.AddSubPage(newAdminPage(globalConfig))
	root.AddSubPage(newSettingsPage(settingsManager))
	root.AddSubPage(newProfilePage(globalConfig))

	return &Site{
		loggerGetter: loggerGetter, localesManager: localesManager,
		authService: globalConfig.AdminService, timeOut: globalConfig.ServiceTimeOut, root: root,
	}
}

func (site *Site) AddPage(page Page) {
	site.root.AddSubPage(page)
}

func (site *Site) AddStaticPages(loggerGetter common.LoggerGetter, groupId uint64, pagePaths []string) {
	site.root.AddStaticPages(loggerGetter, groupId, pagePaths)
}

func (site *Site) GetPage(name string) (Page, bool) {
	return site.root.GetSubPage(name)
}

func (site *Site) GetPageWithPath(path string) (Page, bool) {
	return site.root.GetSubPageWithPath(path)
}

func (site *Site) AddDefaultData(adder common.DataAdder) {
	site.adders = append(site.adders, adder)
}

func (site *Site) manageTimeOut(c *gin.Context) {
	newCtx, cancel := context.WithTimeout(c.Request.Context(), site.timeOut)
	defer cancel()

	c.Request = c.Request.WithContext(newCtx)
	c.Next()
}

func (site *Site) initEngine(siteConfig config.SiteConfig) *gin.Engine {
	engine := gin.New()
	engine.Use(site.manageTimeOut, otelgin.Middleware(config.WebKey), gin.Recovery())

	if memorySize := siteConfig.MaxMultipartMemory; memorySize != 0 {
		engine.MaxMultipartMemory = memorySize
	}

	engine.HTMLRender = templates.NewServiceRender(siteConfig.ExtractTemplateConfig())

	engine.Static("/static", siteConfig.StaticPath)
	engine.StaticFile(config.DefaultFavicon, siteConfig.FaviconPath)

	engine.Use(func(c *gin.Context) {
		c.Set(siteName, site)
	}, makeSessionManager(siteConfig.ExtractSessionConfig()).manage)

	if localesManager := site.localesManager; localesManager.GetMultipleLang() {
		engine.GET("/changeLang", common.CreateRedirect(changeLangRedirecter))

		langPicturePaths := siteConfig.LangPicturePaths
		for _, lang := range localesManager.GetAllLang() {
			if langPicturePath, ok := langPicturePaths[lang]; ok {
				// allow modified time check (instead of always sending same data)
				engine.StaticFile("/langPicture/"+lang, langPicturePath)
			}
		}
	}

	site.root.Widget.LoadInto(engine)
	engine.NoRoute(common.CreateRedirectString(siteConfig.Page404Url))
	return engine
}

func (site *Site) Run(siteConfig config.SiteConfig) error {
	return site.initEngine(siteConfig).Run(common.CheckPort(siteConfig.Port))
}

type SiteAndConfig struct {
	Site   *Site
	Config config.SiteConfig
}

func Run(ginLogger *zap.Logger, sites ...SiteAndConfig) error {
	var g errgroup.Group
	for _, siteAndConfig := range sites {
		port := common.CheckPort(siteAndConfig.Config.Port)
		handler := siteAndConfig.Site.initEngine(siteAndConfig.Config).Handler()
		g.Go(func() error {
			server := &http.Server{Addr: port, Handler: handler}
			return server.ListenAndServe()
		})
	}
	return g.Wait()
}

func changeLangRedirecter(c *gin.Context) string {
	getSite(c).localesManager.SetLangCookie(c.Query(locale.LangName), c)
	return c.Query(common.RedirectName)
}

func BuildDefaultSite(globalConfig *config.GlobalConfig) *Site {
	localesManager := locale.NewManager(globalConfig.ExtractLocalesConfig())
	settingsManager := NewSettingsManager(globalConfig.ExtractSettingsConfig())

	site := NewSite(globalConfig, localesManager, settingsManager)

	return site
}
