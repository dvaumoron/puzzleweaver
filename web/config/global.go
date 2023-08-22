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

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	blogclient "github.com/dvaumoron/puzzleweaver/web/blog/client"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumclient "github.com/dvaumoron/puzzleweaver/web/forum/client"
	"github.com/dvaumoron/puzzleweaver/web/loginclient"
	"github.com/dvaumoron/puzzleweaver/web/profileclient"
	remotewidgetclient "github.com/dvaumoron/puzzleweaver/web/remotewidget/client"
	remotewidgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	wikiclient "github.com/dvaumoron/puzzleweaver/web/wiki/client"
)

const WebKey = "puzzleWeaver"

const defaultName = "default"
const defaultSessionTimeOut = 1200
const defaultServiceTimeOut = 5 * time.Second

const DefaultFavicon = "/favicon.ico"

type GlobalConfig struct {
	Domain string

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

	ProfileGroupId            uint64
	ProfileDefaultPicturePath string

	PageGroups  []PageGroup
	Widgets     map[string]WidgetConfig
	WidgetPages []WidgetPageConfig
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

type WidgetPageConfig struct {
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
	LoginService            service.LoginService
	AdminService            service.AdminService
	ProfileService          service.ProfileService
	ForumService            remoteservice.RemoteForumService
	MarkdownService         service.MarkdownService
	BlogService             remoteservice.RemoteBlogService
	WikiService             remoteservice.RemoteWikiService
	WidgetService           remoteservice.RemoteWidgetService
}

func New(globalConfig *GlobalConfig, loggerGetter common.LoggerGetter, sessionService service.SessionService, templateService service.TemplateService, settingsService service.SettingsService, passwordStrengthService service.PasswordStrengthService, saltService service.SaltService, loginService remoteservice.RemoteLoginService, adminService service.AdminService, profileService remoteservice.RemoteProfileService, forumService remoteservice.RemoteForumService, markdownService service.MarkdownService, blogService remoteservice.RemoteBlogService, wikiService remoteservice.RemoteWikiService, widgetService remoteservice.RemoteWidgetService) *GlobalServiceConfig {
	// TODO read default picture file
	var defaultPicture []byte

	loginServiceWrapper := loginclient.MakeLoginServiceWrapper(
		loginService, saltService, passwordStrengthService, globalConfig.DateFormat,
	)
	profileServiceWrapper := profileclient.MakeProfileServiceWrapper(
		profileService, loginServiceWrapper, adminService, loggerGetter, globalConfig.ProfileGroupId, defaultPicture,
	)

	return &GlobalServiceConfig{
		GlobalConfig:            globalConfig,
		LoggerGetter:            loggerGetter,
		SessionService:          sessionService,
		TemplateService:         templateService,
		SettingsService:         settingsService,
		PasswordStrengthService: passwordStrengthService,
		SaltService:             saltService,
		LoginService:            loginServiceWrapper,
		AdminService:            adminService,
		ProfileService:          profileServiceWrapper,
		ForumService:            forumService,
		MarkdownService:         markdownService,
		BlogService:             blogService,
		WikiService:             wikiService,
		WidgetService:           widgetService,
	}
}

func (c *GlobalServiceConfig) CreateBlogConfig(blogId uint64, groupId uint64, args ...string) BlogConfig {
	return BlogConfig{
		BlogService: blogclient.MakeBlogServiceWrapper(
			c.BlogService, c.AdminService, c.ProfileService, c.LoggerGetter, blogId, groupId, c.DateFormat,
		),
		CommentService: forumclient.MakeForumServiceWrapper(
			c.ForumService, c.AdminService, c.ProfileService, c.LoggerGetter, blogId, groupId, c.DateFormat,
		),
		MarkdownService: c.MarkdownService,
		Domain:          c.Domain, DateFormat: c.DateFormat, PageSize: c.PageSize, ExtractSize: c.ExtractSize,
		FeedFormat: c.FeedFormat, FeedSize: c.FeedSize, Args: args,
	}
}

func (c *GlobalServiceConfig) CreateForumConfig(forumId uint64, groupId uint64, args ...string) ForumConfig {
	return ForumConfig{
		ForumService: forumclient.MakeForumServiceWrapper(
			c.ForumService, c.AdminService, c.ProfileService, c.LoggerGetter, forumId, groupId, c.DateFormat,
		),
		PageSize: c.PageSize, Args: args,
	}
}

func (c *GlobalServiceConfig) CreateWidgetConfig(widgetName string, objectId uint64, groupId uint64) remotewidgetservice.WidgetService {
	return remotewidgetclient.MakeWidgetServiceWrapper(c.WidgetService, c.LoggerGetter, widgetName, objectId, groupId)
}

func (c *GlobalServiceConfig) CreateWikiConfig(wikiId uint64, groupId uint64, args ...string) WikiConfig {
	return WikiConfig{
		WikiService: wikiclient.MakeWikiServiceWrapper(
			c.WikiService, c.AdminService, c.ProfileService, c.LoggerGetter, wikiId, groupId, c.DateFormat,
		),
		Args: args,
	}
}
