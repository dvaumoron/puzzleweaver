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

	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
)

const WebKey = "puzzleWeaver"

const defaultName = "default"
const defaultSessionTimeOut = 1200
const defaultServiceTimeOut = 5 * time.Second

const DefaultFavicon = "/favicon.ico"

type BaseConfigExtracter interface {
	GetServiceTimeOut() time.Duration
	ExtractAdminConfig() AdminConfig
}

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
}

func (c *GlobalConfig) GetServiceTimeOut() time.Duration {
	return c.ServiceTimeOut
}

func (c *GlobalConfig) ExtractLocalesConfig() LocalesConfig {
	return LocalesConfig{Domain: c.Domain, SessionTimeOut: c.SessionTimeOut, AllLang: c.AllLang}
}

func (c *GlobalConfig) ExtractSiteConfig() SiteConfig {
	return SiteConfig{
		Domain: c.Domain, Port: c.Port, SessionTimeOut: c.SessionTimeOut, MaxMultipartMemory: c.MaxMultipartMemory,
		StaticPath: c.StaticPath, FaviconPath: c.FaviconPath, LangPicturePaths: c.LangPicturePaths, Page404Url: c.Page404Url,
	}
}

func (c *GlobalConfig) ExtractAdminConfig() AdminConfig {
	return AdminConfig{PageSize: c.PageSize}
}

func (c *GlobalConfig) CreateWikiConfig(wikiId uint64, groupId uint64, args ...string) WikiConfig {
	return WikiConfig{Args: args}
}

func (c *GlobalConfig) CreateForumConfig(forumId uint64, groupId uint64, args ...string) ForumConfig {
	return ForumConfig{PageSize: c.PageSize, Args: args}
}

func (c *GlobalConfig) CreateBlogConfig(blogId uint64, groupId uint64, args ...string) BlogConfig {
	return BlogConfig{
		Domain: c.Domain, Port: c.Port, DateFormat: c.DateFormat, PageSize: c.PageSize, ExtractSize: c.ExtractSize,
		FeedFormat: c.FeedFormat, FeedSize: c.FeedSize, Args: args,
	}
}

func (c *GlobalConfig) CreateWidgetConfig(serviceAddr string, objectId uint64, groupId uint64) widgetservice.WidgetService {
	return nil
}
