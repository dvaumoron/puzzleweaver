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
	"slices"
	"strings"
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
	forumclient "github.com/dvaumoron/puzzleweaver/web/forum/client"
	"github.com/dvaumoron/puzzleweaver/web/loginclient"
	"github.com/dvaumoron/puzzleweaver/web/profileclient"
	remotewidgetclient "github.com/dvaumoron/puzzleweaver/web/remotewidget/client"
	wikiclient "github.com/dvaumoron/puzzleweaver/web/wiki/client"
	"github.com/dvaumoron/puzzleweb/common/config"
	"github.com/dvaumoron/puzzleweb/common/config/parser"
	"github.com/dvaumoron/puzzleweb/common/log"
	forumservice "github.com/dvaumoron/puzzleweb/forum/service"
	loginservice "github.com/dvaumoron/puzzleweb/login/service"
	profileservice "github.com/dvaumoron/puzzleweb/profile/service"
	"github.com/spf13/afero"
)

type GlobalConfig struct {
	Domain string
	Port   string

	DefaultLang        string
	SessionTimeOut     int
	ServiceTimeOut     time.Duration
	MaxMultipartMemory int64
	DateFormat         string
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

	StaticPages []parser.StaticPagesConfig
	Widgets     map[string]parser.WidgetConfig
	WidgetPages []parser.WidgetPageConfig

	GinReleaseMode bool
}

type GlobalServiceConfig struct {
	*GlobalConfig
	LoggerGetter            log.LoggerGetter
	VersionName             string
	AllLang                 []string
	StaticFileSystem        http.FileSystem
	SessionService          sessionimpl.SessionService
	TemplateService         templatesimpl.TemplateService
	SettingsService         settingsimpl.SettingsService
	PasswordStrengthService passwordstrengthimpl.PasswordStrengthService
	SaltService             saltimpl.SaltService
	LoginService            loginservice.FullLoginService
	AdminService            adminimpl.AdminService
	ProfileService          profileservice.ProfileService
	ForumService            forumimpl.RemoteForumService
	MarkdownService         markdownimpl.MarkdownService
	BlogService             blogimpl.RemoteBlogService
	WikiService             wikiimpl.RemoteWikiService
	WidgetService           remotewidgetimpl.RemoteWidgetService
}

func New(conf *GlobalConfig, loggerGetter log.LoggerGetter, logger *slog.Logger, version string, sessionService sessionimpl.SessionService, templateService templatesimpl.TemplateService, settingsService settingsimpl.SettingsService, passwordStrengthService passwordstrengthimpl.PasswordStrengthService, saltService saltimpl.SaltService, loginService loginimpl.RemoteLoginService, adminService adminimpl.AdminService, profileService profileimpl.RemoteProfileService, forumService forumimpl.RemoteForumService, markdownService markdownimpl.MarkdownService, blogService blogimpl.RemoteBlogService, wikiService wikiimpl.RemoteWikiService, widgetService remotewidgetimpl.RemoteWidgetService) (*GlobalServiceConfig, error) {
	allLang := make([]string, 0, len(conf.LangPicturePaths))
	for lang := range conf.LangPicturePaths {
		allLang = append(allLang, lang)
	}
	slices.Sort(allLang)

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
		VersionName:             "PuzzleWeaver" + version,
		AllLang:                 allLang,
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

func (c *GlobalServiceConfig) GetLogger() log.Logger {
	// TODO
	return nil
}

func (c *GlobalServiceConfig) CreateBlogConfig(widgetConfig parser.WidgetConfig) (config.BlogConfig, bool) {
	blogService := blogclient.MakeBlogServiceWrapper(
		c.BlogService, c.AdminService, c.ProfileService, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	commentService := forumclient.MakeForumServiceWrapper(
		c.ForumService, c.AdminService, c.ProfileService, c.LoggerGetter, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	return config.BlogConfig{
		ServiceConfig: config.MakeServiceConfig(c, blogService), CommentService: commentService, MarkdownService: c.MarkdownService,
		Domain: c.Domain, Port: c.Port, DateFormat: c.DateFormat, PageSize: c.PageSize, ExtractSize: c.ExtractSize,
		FeedFormat: c.FeedFormat, FeedSize: c.FeedSize, Args: widgetConfig.Templates,
	}, true
}

func (c *GlobalServiceConfig) CreateForumConfig(widgetConfig parser.WidgetConfig) (config.ForumConfig, bool) {
	forumService := forumclient.MakeForumServiceWrapper(
		c.ForumService, c.AdminService, c.ProfileService, c.LoggerGetter, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	return config.ForumConfig{
		ServiceConfig: config.MakeServiceConfig[forumservice.ForumService](c, forumService),
		PageSize:      c.PageSize, Args: widgetConfig.Templates,
	}, true
}

func (c *GlobalServiceConfig) CreateWidgetConfig(widgetConfig parser.WidgetConfig) (config.RemoteWidgetConfig, bool) {
	widgetName, remoteKind := strings.CutPrefix(widgetConfig.Kind, "remote/")
	return config.MakeServiceConfig(c, remotewidgetclient.MakeWidgetServiceWrapper(
		c.WidgetService, c.LoggerGetter, widgetName, widgetConfig.ObjectId, widgetConfig.GroupId,
	)), remoteKind
}

func (c *GlobalServiceConfig) CreateWikiConfig(wikiId uint64, groupId uint64, args ...string) (config.WikiConfig, bool) {
	wikiService := wikiclient.MakeWikiServiceWrapper(
		c.WikiService, c.AdminService, c.ProfileService, c.LoggerGetter, wikiId, groupId, c.DateFormat,
	)
	return config.WikiConfig{
		ServiceConfig:   config.MakeServiceConfig(c, wikiService),
		MarkdownService: c.MarkdownService, Args: args,
	}, true
}
