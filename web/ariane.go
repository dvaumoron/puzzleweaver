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
	"net/url"
	"strings"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/locale"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

const errorMsgName = "ErrorMsg"

type PageDesc struct {
	Name string
	Url  string
}

func makePageDesc(name string, url string) PageDesc {
	return PageDesc{Name: getPageTitleKey(name), Url: url}
}

func getPageTitleKey(name string) string {
	return "PageTitle" + locale.CamelCase(name)
}

func buildAriane(splittedPath []string) []PageDesc {
	pageDescs := make([]PageDesc, 0, len(splittedPath))
	var urlBuilder strings.Builder
	for _, name := range splittedPath {
		urlBuilder.WriteByte('/')
		urlBuilder.WriteString(name)
		pageDescs = append(pageDescs, makePageDesc(name, urlBuilder.String()))
	}
	return pageDescs
}

func getSite(c *gin.Context) *Site {
	siteUntyped, _ := c.Get(siteName)
	return siteUntyped.(*Site)
}

func GetLogger(c *gin.Context) *slog.Logger {
	return getSite(c).loggerGetter.Logger(c.Request.Context())
}

func GetLocalesManager(c *gin.Context) locale.Manager {
	return getSite(c).localesManager
}

func InitNoELementMsg(data gin.H, size int) {
	if size == 0 {
		data[errorMsgName] = "NoElement"
	}
}

func (site *Site) extractArianeInfoFromUrl(url string) (Page, []string) {
	current := site.root
	splitted := strings.Split(url, "/")[1:]
	names := make([]string, 0, len(splitted))
	for _, name := range splitted {
		subPage, ok := current.GetSubPage(name)
		if !ok {
			break
		}
		current = subPage
		names = append(names, name)
	}
	return current, names
}

func (p Page) extractSubPageNames(url string, c *gin.Context) []PageDesc {
	sw, ok := p.Widget.(*staticWidget)
	if !ok {
		return nil
	}

	pages := sw.subPages
	size := len(pages)
	if size == 0 {
		return nil
	}

	pageDescs := make([]PageDesc, 0, size)
	for _, page := range pages {
		if page.visible {
			name := page.name
			pageDescs = append(pageDescs, makePageDesc(name, url+name))
		}
	}
	return pageDescs
}

func initData(c *gin.Context) gin.H {
	site := getSite(c)
	ctx := c.Request.Context()
	localesManager := site.localesManager
	currentUrl := common.GetCurrentUrl(c)
	page, path := site.extractArianeInfoFromUrl(currentUrl)
	data := gin.H{
		locale.LangName: localesManager.GetLang(site.loggerGetter.Logger(ctx), c),
		"PageTitle":     getPageTitleKey(page.name),
		common.UrlName:  currentUrl,
		"Ariane":        buildAriane(path),
		"SubPages":      page.extractSubPageNames(currentUrl, c),
		errorMsgName:    c.Query(common.ErrorKey),
	}
	escapedUrl := url.QueryEscape(c.Request.URL.Path)
	if localesManager.GetMultipleLang() {
		data["LangSelectorUrl"] = "/changeLang?Redirect=" + escapedUrl
		data["AllLang"] = localesManager.GetAllLang()
	}
	session := GetSession(c)
	var currentUserId uint64
	if login := session.Load(loginName); login == "" {
		data[loginUrlName] = "/login?Redirect=" + escapedUrl
	} else {
		currentUserId = GetSessionUserId(c)
		data[loginName] = login
		data[common.UserIdName] = currentUserId
		data[loginUrlName] = "/login/logout?Redirect=" + escapedUrl
	}
	data[viewAdminName] = site.authService.AuthQuery(
		ctx, currentUserId, service.AdminGroupId, service.ActionAccess,
	) == nil
	for _, adder := range site.adders {
		adder(data, c)
	}
	return data
}
