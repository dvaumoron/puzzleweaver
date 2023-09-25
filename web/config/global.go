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
	"log/slog"
	"net/http"
	"time"

	fsclient "github.com/dvaumoron/puzzleweaver/client/fs"
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
	blogclient "github.com/dvaumoron/puzzleweaver/web/blog/client"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumclient "github.com/dvaumoron/puzzleweaver/web/forum/client"
	"github.com/dvaumoron/puzzleweaver/web/loginclient"
	"github.com/dvaumoron/puzzleweaver/web/profileclient"
	remotewidgetclient "github.com/dvaumoron/puzzleweaver/web/remotewidget/client"
	remotewidgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	wikiclient "github.com/dvaumoron/puzzleweaver/web/wiki/client"
	"github.com/spf13/afero"
)

const WebKey = "puzzleWeaver"

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

	FsConf      fsclient.FsConf
	StaticPath  string
	FaviconPath string
	Page404Url  string

	LangPicturePaths map[string]string

	ProfileGroupId            uint64
	ProfileDefaultPicturePath string

	PageGroups  []PageGroup
	Widgets     map[string]WidgetConfig
	WidgetPages []WidgetPageConfig

	GinReleaseMode bool
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
	StaticFileSystem        http.FileSystem
	SessionService          sessionimpl.SessionService
	TemplateService         templatesimpl.TemplateService
	SettingsService         settingsimpl.SettingsService
	PasswordStrengthService passwordstrengthimpl.PasswordStrengthService
	SaltService             saltimpl.SaltService
	LoginService            service.LoginService
	AdminService            adminimpl.AdminService
	ProfileService          service.ProfileService
	ForumService            forumimpl.RemoteForumService
	MarkdownService         markdownimpl.MarkdownService
	BlogService             blogimpl.RemoteBlogService
	WikiService             wikiimpl.RemoteWikiService
	WidgetService           remotewidgetimpl.RemoteWidgetService
}

func New(conf *GlobalConfig, loggerGetter common.LoggerGetter, logger *slog.Logger, sessionService sessionimpl.SessionService, templateService templatesimpl.TemplateService, settingsService settingsimpl.SettingsService, passwordStrengthService passwordstrengthimpl.PasswordStrengthService, saltService saltimpl.SaltService, loginService loginimpl.RemoteLoginService, adminService adminimpl.AdminService, profileService profileimpl.RemoteProfileService, forumService forumimpl.RemoteForumService, markdownService markdownimpl.MarkdownService, blogService blogimpl.RemoteBlogService, wikiService wikiimpl.RemoteWikiService, widgetService remotewidgetimpl.RemoteWidgetService) (*GlobalServiceConfig, error) {
	baseFS, err := fsclient.New(conf.FsConf)
	if err != nil {
		return nil, err
	}

	// read default picture file
	defaultPicture, err := afero.ReadFile(baseFS, conf.ProfileDefaultPicturePath)
	if err != nil {
		return nil, err
	}

	loginServiceWrapper := loginclient.MakeLoginServiceWrapper(
		loginService, saltService, passwordStrengthService, conf.DateFormat,
	)
	profileServiceWrapper := profileclient.MakeProfileServiceWrapper(
		profileService, loginServiceWrapper, adminService, loggerGetter, conf.ProfileGroupId, defaultPicture,
	)

	return &GlobalServiceConfig{
		GlobalConfig:            conf,
		LoggerGetter:            loggerGetter,
		StaticFileSystem:        afero.NewHttpFs(afero.NewBasePathFs(baseFS, conf.StaticPath)),
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
	}, nil
}

func (c *GlobalServiceConfig) CreateBlogConfig(blogId uint64, groupId uint64, args ...string) BlogConfig {
	blogService := blogclient.MakeBlogServiceWrapper(
		c.BlogService, c.AdminService, c.ProfileService, c.LoggerGetter, blogId, groupId, c.DateFormat,
	)
	commentService := forumclient.MakeForumServiceWrapper(
		c.ForumService, c.AdminService, c.ProfileService, c.LoggerGetter, blogId, groupId, c.DateFormat,
	)
	return BlogConfig{
		BlogService: blogService, CommentService: commentService, MarkdownService: c.MarkdownService,
		Domain: c.Domain, Port: c.Port, DateFormat: c.DateFormat, PageSize: c.PageSize, ExtractSize: c.ExtractSize,
		FeedFormat: c.FeedFormat, FeedSize: c.FeedSize, Args: args,
	}
}

func (c *GlobalServiceConfig) CreateForumConfig(forumId uint64, groupId uint64, args ...string) ForumConfig {
	forumService := forumclient.MakeForumServiceWrapper(
		c.ForumService, c.AdminService, c.ProfileService, c.LoggerGetter, forumId, groupId, c.DateFormat,
	)
	return ForumConfig{ForumService: forumService, PageSize: c.PageSize, Args: args}
}

func (c *GlobalServiceConfig) CreateWidgetConfig(widgetName string, objectId uint64, groupId uint64) remotewidgetservice.WidgetService {
	return remotewidgetclient.MakeWidgetServiceWrapper(c.WidgetService, c.LoggerGetter, widgetName, objectId, groupId)
}

func (c *GlobalServiceConfig) CreateWikiConfig(wikiId uint64, groupId uint64, args ...string) WikiConfig {
	wikiService := wikiclient.MakeWikiServiceWrapper(
		c.WikiService, c.AdminService, c.ProfileService, c.LoggerGetter, wikiId, groupId, c.DateFormat,
	)
	return WikiConfig{WikiService: wikiService, MarkdownService: c.MarkdownService, Args: args}
}
