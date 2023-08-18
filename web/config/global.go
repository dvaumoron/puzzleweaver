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

package config

import (
	"time"

	blogservice "github.com/dvaumoron/puzzleweaver/web/blog/service"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
	remotewidgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
)

const WebKey = "puzzleWeaver"

const defaultName = "default"
const defaultSessionTimeOut = 1200
const defaultServiceTimeOut = 5 * time.Second

const DefaultFavicon = "/favicon.ico"

type GlobalConfig struct {
	Domain string
	Port   string

	AllLang            []string
	SessionTimeOut     int
	ServiceTimeOut     time.Duration
	MaxMultipartMemory int64
	DateFormat         string // TODO move this to template service
	PageSize           uint64
	ExtractSize        uint64
	FeedFormat         string
	FeedSize           uint64

	StaticPath  string
	FaviconPath string
	Page404Url  string

	LangPicturePaths map[string]string

	PageGroups  []PageGroup
	Widgets     map[string]WidgetConfig
	WidgetPages []WidgetPage
}

type PageGroup struct {
	Id    uint64
	Pages []string
}

type WidgetConfig struct {
	Kind       string
	WidgetName string
	ObjectId   uint64
	GroupId    uint64
	Templates  []string
}

type WidgetPage struct {
	Emplacement string
	Name        string
	WidgetRef   string
}

type GlobalServiceConfig struct {
	*GlobalConfig
	LoggerGetter            common.LoggerGetter
	SessionService          service.SessionService
	TemplateService         service.TemplateService
	SettingsService         service.SettingsService
	PasswordStrengthService service.PasswordStrengthService
	SaltService             service.SaltService
	LoginService            service.FullLoginService
	AdminService            service.AdminService
	ProfileService          service.AdvancedProfileService
	ForumService            forumservice.FullForumService
	MarkdownService         service.MarkdownService
	BlogService             blogservice.BlogService
	WikiService             wikiservice.WikiService
	WidgetService           remotewidgetservice.WidgetService
}

func (c *GlobalServiceConfig) CreateWikiConfig(wikiId uint64, groupId uint64, args ...string) WikiConfig {
	return WikiConfig{Args: args}
}

func (c *GlobalServiceConfig) CreateForumConfig(forumId uint64, groupId uint64, args ...string) ForumConfig {
	return ForumConfig{PageSize: c.PageSize, Args: args}
}

func (c *GlobalServiceConfig) CreateBlogConfig(blogId uint64, groupId uint64, args ...string) BlogConfig {
	return BlogConfig{
		// TODO wrapper with blogId and groupId
		BlogService: c.BlogService, CommentService: c.ForumService, MarkdownService: c.MarkdownService,
		Domain: c.Domain, Port: c.Port, DateFormat: c.DateFormat, PageSize: c.PageSize, ExtractSize: c.ExtractSize,
		FeedFormat: c.FeedFormat, FeedSize: c.FeedSize, Args: args,
	}
}

func (c *GlobalServiceConfig) CreateWidgetConfig(objectId uint64, groupId uint64) widgetservice.WidgetService {
	// TODO wrapper with objectId and groupId
	return c.WidgetService
}
