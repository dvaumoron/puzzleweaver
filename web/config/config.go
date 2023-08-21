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
	blogservice "github.com/dvaumoron/puzzleweaver/web/blog/service"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
)

type BlogConfig struct {
	BlogService     blogservice.BlogService
	CommentService  forumservice.CommentService
	MarkdownService service.MarkdownService
	Domain          string
	Port            string
	DateFormat      string
	PageSize        uint64
	ExtractSize     uint64
	FeedFormat      string
	FeedSize        uint64
	Args            []string
}

type ForumConfig struct {
	ForumService forumservice.ForumService
	PageSize     uint64
	Args         []string
}

type WikiConfig struct {
	WikiService     wikiservice.WikiService
	MarkdownService service.MarkdownService
	Args            []string
}

type WidgetServiceConfig struct {
	WidgetService widgetservice.WidgetService
	LoggerGetter  common.LoggerGetter
}
