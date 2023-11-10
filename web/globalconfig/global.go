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

package globalconfig

import (
	"context"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strings"
	"time"

	fsclient "github.com/dvaumoron/puzzleweaver/client/fs"
	adminimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/admin"
	blogimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/blog"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
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
	"github.com/dvaumoron/puzzleweaver/web/adminclient"
	blogclient "github.com/dvaumoron/puzzleweaver/web/blogclient"
	forumclient "github.com/dvaumoron/puzzleweaver/web/forumclient"
	"github.com/dvaumoron/puzzleweaver/web/loginclient"
	"github.com/dvaumoron/puzzleweaver/web/profileclient"
	remotewidgetclient "github.com/dvaumoron/puzzleweaver/web/remotewidgetclient"
	"github.com/dvaumoron/puzzleweaver/web/templateclient"
	wikiclient "github.com/dvaumoron/puzzleweaver/web/wikiclient"
	adminservice "github.com/dvaumoron/puzzleweb/admin/service"
	"github.com/dvaumoron/puzzleweb/common/config"
	"github.com/dvaumoron/puzzleweb/common/config/parser"
	"github.com/dvaumoron/puzzleweb/common/log"
	forumservice "github.com/dvaumoron/puzzleweb/forum/service"
	loginservice "github.com/dvaumoron/puzzleweb/login/service"
	profileservice "github.com/dvaumoron/puzzleweb/profile/service"
	sessionservice "github.com/dvaumoron/puzzleweb/session/service"
	templateservice "github.com/dvaumoron/puzzleweb/templates/service"
	"github.com/spf13/afero"
	"go.uber.org/zap/zapcore"
)

type loggerWrapper struct {
	inner *slog.Logger
}

func (lw loggerWrapper) Debug(msg string, args ...zapcore.Field) {
	lw.inner.Debug(msg, convertLogArgs(args)...)
}

func (lw loggerWrapper) Info(msg string, args ...zapcore.Field) {
	lw.inner.Info(msg, convertLogArgs(args)...)
}

func (lw loggerWrapper) Warn(msg string, args ...zapcore.Field) {
	lw.inner.Warn(msg, convertLogArgs(args)...)
}

func (lw loggerWrapper) Error(msg string, args ...zapcore.Field) {
	lw.inner.Error(msg, convertLogArgs(args)...)
}

func convertLogArgs(args []zapcore.Field) []any {
	resArgs := make([]any, 0, 2*len(args))
	for _, arg := range args {
		resArgs = append(resArgs, arg.Key, extractZapValue(arg))
	}
	return resArgs
}

func extractZapValue(arg zapcore.Field) any {
	switch arg.Type {
	case zapcore.BoolType:
		return arg.Integer == 1
	case zapcore.DurationType:
		return time.Duration(arg.Integer)
	case zapcore.Float64Type:
		return math.Float64frombits(uint64(arg.Integer))
	case zapcore.Float32Type:
		return math.Float32frombits(uint32(arg.Integer))
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return arg.Integer
	case zapcore.StringType:
		return arg.String
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
		return uint64(arg.Integer)
	}
	return arg.Interface
}

type loggerGetterWrapper struct {
	inner servicecommon.LoggerGetter
}

func (lgw loggerGetterWrapper) Logger(ctx context.Context) log.Logger {
	return loggerWrapper{inner: lgw.inner.Logger(ctx)}
}

type settingsServiceWrapper struct {
	settingsimpl.SettingsService
}

func (_ settingsServiceWrapper) Generate(_ context.Context) (uint64, error) {
	return 0, nil
}

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
	Logger                  log.Logger
	VersionName             string
	AllLang                 []string
	StaticFileSystem        http.FileSystem
	SessionService          sessionservice.SessionService
	TemplateService         templateservice.TemplateService
	SettingsService         sessionservice.SessionService
	PasswordStrengthService passwordstrengthimpl.PasswordStrengthService
	LoginService            loginservice.FullLoginService
	AdminService            adminservice.AdminService
	ProfileService          profileservice.AdvancedProfileService
	ForumImpl               forumimpl.RemoteForumService
	MarkdownImpl            markdownimpl.MarkdownService
	BlogImpl                blogimpl.RemoteBlogService
	WikiImpl                wikiimpl.RemoteWikiService
	WidgetImpl              remotewidgetimpl.RemoteWidgetService
}

func New(conf *GlobalConfig, loggerGetter servicecommon.LoggerGetter, logger *slog.Logger, version string, sessionService sessionimpl.SessionService, templateService templatesimpl.TemplateService, settingsService settingsimpl.SettingsService, passwordStrengthService passwordstrengthimpl.PasswordStrengthService, saltService saltimpl.SaltService, loginService loginimpl.RemoteLoginService, adminService adminimpl.AdminService, profileService profileimpl.RemoteProfileService, forumService forumimpl.RemoteForumService, markdownService markdownimpl.MarkdownService, blogService blogimpl.RemoteBlogService, wikiService wikiimpl.RemoteWikiService, widgetService remotewidgetimpl.RemoteWidgetService) (*GlobalServiceConfig, error) {
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

	wrappedLoggerGetter := loggerGetterWrapper{inner: loggerGetter}

	loginServiceWrapper := loginclient.MakeLoginServiceWrapper(
		loginService, saltService, passwordStrengthService, conf.DateFormat,
	)
	profileServiceWrapper := profileclient.MakeProfileServiceWrapper(
		profileService, loginServiceWrapper, adminService, wrappedLoggerGetter, conf.ProfileGroupId, defaultPicture,
	)

	return &GlobalServiceConfig{
		GlobalConfig:            conf,
		LoggerGetter:            wrappedLoggerGetter,
		Logger:                  loggerWrapper{inner: logger},
		VersionName:             "PuzzleWeaver" + version,
		AllLang:                 allLang,
		StaticFileSystem:        afero.NewHttpFs(afero.NewBasePathFs(baseFS, conf.StaticPath)),
		SessionService:          sessionService,
		TemplateService:         templateclient.MakeTemplateServiceWrapper(templateService, wrappedLoggerGetter),
		SettingsService:         settingsServiceWrapper{SettingsService: settingsService},
		PasswordStrengthService: passwordStrengthService,
		LoginService:            loginServiceWrapper,
		AdminService:            adminclient.MakeAdminServiceWrapper(adminService),
		ProfileService:          profileServiceWrapper,
		ForumImpl:               forumService,
		MarkdownImpl:            markdownService,
		BlogImpl:                blogService,
		WikiImpl:                wikiService,
		WidgetImpl:              widgetService,
	}, nil
}

func (c *GlobalServiceConfig) GetLogger() log.Logger {
	return c.Logger
}

func (c *GlobalServiceConfig) GetLoggerGetter() log.LoggerGetter {
	return c.LoggerGetter
}

func (c *GlobalConfig) GetServiceTimeOut() time.Duration {
	return c.ServiceTimeOut
}

func (c *GlobalServiceConfig) ExtractAuthConfig() config.AuthConfig {
	return config.MakeServiceConfig[adminservice.AuthService](c, c.AdminService)
}

func (c *GlobalServiceConfig) ExtractLocalesConfig() config.LocalesConfig {
	return config.LocalesConfig{Logger: c.GetLogger(), Domain: c.Domain, SessionTimeOut: c.SessionTimeOut, AllLang: c.AllLang}
}

func (c *GlobalServiceConfig) ExtractSiteConfig() config.SiteConfig {
	return config.SiteConfig{
		ServiceConfig: config.MakeServiceConfig[sessionservice.SessionService](c, c.SessionService), TemplateService: c.TemplateService,
		Domain: c.Domain, Port: c.Port, SessionTimeOut: c.SessionTimeOut,
		MaxMultipartMemory: c.MaxMultipartMemory, StaticPath: c.StaticPath, FaviconPath: c.FaviconPath,
		LangPicturePaths: c.LangPicturePaths, Page404Url: c.Page404Url,
	}
}

func (c *GlobalServiceConfig) ExtractLoginConfig() config.LoginConfig {
	return config.MakeServiceConfig[loginservice.LoginService](c, c.LoginService)
}

func (c *GlobalServiceConfig) ExtractAdminConfig() config.AdminConfig {
	return config.AdminConfig{
		ServiceConfig: config.MakeServiceConfig[adminservice.AdminService](c, c.AdminService),
		UserService:   c.LoginService, ProfileService: c.ProfileService, PageSize: c.PageSize,
	}
}

func (c *GlobalServiceConfig) ExtractProfileConfig() config.ProfileConfig {
	return config.ProfileConfig{
		ServiceConfig: config.MakeServiceConfig[profileservice.AdvancedProfileService](c, c.ProfileService),
		AdminService:  c.AdminService, LoginService: c.LoginService,
	}
}

func (c *GlobalServiceConfig) ExtractSettingsConfig() config.SettingsConfig {
	return config.MakeServiceConfig(c, c.SettingsService)
}

func (c *GlobalServiceConfig) MakeBlogConfig(widgetConfig parser.WidgetConfig) (config.BlogConfig, bool) {
	blogService := blogclient.MakeBlogServiceWrapper(
		c.BlogImpl, c.AdminService, c.ProfileService, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	commentService := forumclient.MakeForumServiceWrapper(
		c.ForumImpl, c.AdminService, c.ProfileService, c.LoggerGetter, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	return config.BlogConfig{
		ServiceConfig: config.MakeServiceConfig(c, blogService), CommentService: commentService, MarkdownService: c.MarkdownImpl,
		Domain: c.Domain, Port: c.Port, DateFormat: c.DateFormat, PageSize: c.PageSize, ExtractSize: c.ExtractSize,
		FeedFormat: c.FeedFormat, FeedSize: c.FeedSize, Args: widgetConfig.Templates,
	}, true
}

func (c *GlobalServiceConfig) MakeForumConfig(widgetConfig parser.WidgetConfig) (config.ForumConfig, bool) {
	forumService := forumclient.MakeForumServiceWrapper(
		c.ForumImpl, c.AdminService, c.ProfileService, c.LoggerGetter, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	return config.ForumConfig{
		ServiceConfig: config.MakeServiceConfig[forumservice.ForumService](c, forumService),
		PageSize:      c.PageSize, Args: widgetConfig.Templates,
	}, true
}

func (c *GlobalServiceConfig) MakeWidgetConfig(widgetConfig parser.WidgetConfig) (config.RemoteWidgetConfig, bool) {
	widgetName, remoteKind := strings.CutPrefix(widgetConfig.Kind, "remote/")
	return config.MakeServiceConfig(c, remotewidgetclient.MakeWidgetServiceWrapper(
		c.WidgetImpl, c.LoggerGetter, widgetName, widgetConfig.ObjectId, widgetConfig.GroupId,
	)), remoteKind
}

func (c *GlobalServiceConfig) MakeWikiConfig(widgetConfig parser.WidgetConfig) (config.WikiConfig, bool) {
	wikiService := wikiclient.MakeWikiServiceWrapper(
		c.WikiImpl, c.AdminService, c.ProfileService, c.LoggerGetter, widgetConfig.ObjectId, widgetConfig.GroupId, c.DateFormat,
	)
	return config.WikiConfig{
		ServiceConfig:   config.MakeServiceConfig(c, wikiService),
		MarkdownService: c.MarkdownImpl, Args: widgetConfig.Templates,
	}, true
}
