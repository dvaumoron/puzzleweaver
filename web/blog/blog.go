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

package blog

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dvaumoron/puzzleweaver/web"
	blogservice "github.com/dvaumoron/puzzleweaver/web/blog/service"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/config"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"golang.org/x/exp/slog"
)

const emptyTitle = "EmptyPostTitle"
const emptyContent = "EmptyPostContent"

const postIdName = "postId"
const commentMsgName = "CommentMsg"

const parsingPostIdErrorMsg = "Failed to parse postId"

var errEmptyComment = errors.New("EmptyComment")
var errFeedFormat = errors.New("unrecognized feed format")

// TODO draft with modify until publish ?
type blogWidget struct {
	listHandler          gin.HandlerFunc
	viewHandler          gin.HandlerFunc
	saveCommentHandler   gin.HandlerFunc
	deleteCommentHandler gin.HandlerFunc
	createHandler        gin.HandlerFunc
	previewHandler       gin.HandlerFunc
	saveHandler          gin.HandlerFunc
	deleteHandler        gin.HandlerFunc
	rssHandler           gin.HandlerFunc
}

func (w blogWidget) LoadInto(router gin.IRouter) {
	router.GET("/", w.listHandler)
	router.GET("/view/:postId", w.viewHandler)
	router.POST("/comment/save/:postId", w.saveCommentHandler)
	router.GET("/comment/delete/:postId/:commentId", w.deleteCommentHandler)
	router.GET("/create", w.createHandler)
	router.POST("/preview", w.previewHandler)
	router.POST("/save", w.saveHandler)
	router.GET("/delete/:postId", w.deleteHandler)
	router.GET("/rss", w.rssHandler)
}

func MakeBlogPage(blogName string, logger *slog.Logger, blogService blogservice.BlogService, commentService forumservice.CommentService, markdownService service.MarkdownService, blogConfig config.BlogConfig) web.Page {
	host := blogConfig.Domain
	if port := common.CheckPort(blogConfig.Port); port != ":80" {
		host += port
	}
	dateFormat := blogConfig.DateFormat
	defaultPageSize := blogConfig.PageSize
	extractSize := blogConfig.ExtractSize
	feedFormat := blogConfig.FeedFormat
	feedSize := blogConfig.FeedSize

	listTmpl := "blog/list"
	viewTmpl := "blog/view"
	createTmpl := "blog/create"
	previewTmpl := "blog/preview"
	switch args := blogConfig.Args; len(args) {
	default:
		logger.Info("MakeBlogPage should be called with 0 to 4 optional arguments.")
		fallthrough
	case 4:
		if args[3] != "" {
			previewTmpl = args[3]
		}
		fallthrough
	case 3:
		if args[2] != "" {
			createTmpl = args[2]
		}
		fallthrough
	case 2:
		if args[1] != "" {
			viewTmpl = args[1]
		}
		fallthrough
	case 1:
		if args[0] != "" {
			listTmpl = args[0]
		}
	case 0:
	}

	p := web.MakePage(blogName)
	p.Widget = blogWidget{
		listHandler: web.CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			ctx := c.Request.Context()
			userId, _ := data[common.IdName].(uint64)

			pageNumber, start, end, filter := common.GetPagination(defaultPageSize, c)

			total, posts, err := blogService.GetPosts(ctx, userId, start, end, filter)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			filterPostsExtract(posts, extractSize)

			common.InitPagination(data, filter, pageNumber, end, total)
			data["Posts"] = posts
			data[common.AllowedToCreateName] = blogService.CreateRight(ctx, userId) == nil
			data[common.AllowedToDeleteName] = blogService.DeleteRight(ctx, userId) == nil
			web.InitNoELementMsg(data, len(posts), c)
			return listTmpl, ""
		}),
		viewHandler: web.CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			ctx := c.Request.Context()
			logger := web.GetLogger(c)
			userId, _ := data[common.IdName].(uint64)

			pageNumber, start, end, _ := common.GetPagination(defaultPageSize, c)

			postId, err := strconv.ParseUint(c.Param(postIdName), 10, 64)
			if err != nil {
				logger.Warn(parsingPostIdErrorMsg, common.ErrorKey, err)
				return "", common.DefaultErrorRedirect(common.ErrorTechnicalKey)
			}

			post, err := blogService.GetPost(ctx, userId, postId)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			total, comments, err := commentService.GetCommentThread(ctx, userId, post.Title, start, end)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			common.InitPagination(data, "", pageNumber, end, total)
			data[common.BaseUrlName] = common.GetBaseUrl(2, c)
			data["Post"] = post
			data["Comments"] = comments
			data[common.AllowedToCreateName] = commentService.CreateMessageRight(ctx, userId) == nil
			data[common.AllowedToDeleteName] = commentService.DeleteRight(ctx, userId) == nil
			if len(comments) == 0 {
				if err == nil {
					data[commentMsgName] = "NoComment"
				} else {
					data[commentMsgName] = "CommentDisplayError"
				}
			}
			return viewTmpl, ""
		}),
		saveCommentHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			logger := web.GetLogger(c)
			userId := web.GetSessionUserId(c)

			postId, err := strconv.ParseUint(c.Param(postIdName), 10, 64)
			if err != nil {
				logger.Warn(parsingPostIdErrorMsg, common.ErrorKey, err)
				return common.DefaultErrorRedirect(common.ErrorTechnicalKey)
			}
			comment := c.PostForm("comment")

			err = errEmptyComment
			if comment != "" {
				var post blogservice.BlogPost
				post, err = blogService.GetPost(ctx, userId, postId)
				if err != nil {
					return common.DefaultErrorRedirect(err.Error())
				}

				err = commentService.CreateComment(ctx, userId, post.Title, comment)
			}

			targetBuilder := postUrlBuilder(common.GetBaseUrl(3, c), postId)
			if err != nil {
				common.WriteError(targetBuilder, err.Error())
			}
			return targetBuilder.String()
		}),
		deleteCommentHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			logger := web.GetLogger(c)
			userId := web.GetSessionUserId(c)

			postId, err := strconv.ParseUint(c.Param(postIdName), 10, 64)
			if err != nil {
				logger.Warn(parsingPostIdErrorMsg, common.ErrorKey, err)
				return common.DefaultErrorRedirect(common.ErrorTechnicalKey)
			}
			commentId, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
			if err != nil {
				logger.Warn("Failed to parse commentId", common.ErrorKey, err)
				return common.DefaultErrorRedirect(common.ErrorTechnicalKey)
			}

			post, err := blogService.GetPost(ctx, userId, postId)
			if err != nil {
				return common.DefaultErrorRedirect(err.Error())
			}

			err = commentService.DeleteComment(ctx, userId, post.Title, commentId)
			targetBuilder := postUrlBuilder(common.GetBaseUrl(4, c), postId)
			if err != nil {
				common.WriteError(targetBuilder, err.Error())
			}
			return targetBuilder.String()
		}),
		createHandler: web.CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			data[common.BaseUrlName] = common.GetBaseUrl(1, c)
			return createTmpl, ""
		}),
		previewHandler: web.CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			title := c.PostForm("title")
			markdown := c.PostForm("markdown")

			if title == "" {
				return "", common.DefaultErrorRedirect(emptyTitle)
			}
			if markdown == "" {
				return "", common.DefaultErrorRedirect(emptyContent)
			}

			html, err := markdownService.Apply(c.Request.Context(), markdown)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			data[common.BaseUrlName] = common.GetBaseUrl(1, c)
			data["PreviewTitle"] = title
			data["Markdown"] = markdown
			data["PreviewHTML"] = html
			return previewTmpl, ""
		}),
		saveHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			title := c.PostForm("title")
			userId := web.GetSessionUserId(c)
			markdown := c.PostForm("markdown")

			if title == "" {
				return common.DefaultErrorRedirect(emptyTitle)
			}
			if markdown == "" {
				return common.DefaultErrorRedirect(emptyContent)
			}

			html, err := markdownService.Apply(ctx, markdown)
			if err != nil {
				return common.DefaultErrorRedirect(err.Error())
			}

			postId, err := blogService.CreatePost(ctx, userId, title, string(html))
			if err != nil {
				return common.DefaultErrorRedirect(err.Error())
			}

			err = commentService.CreateCommentThread(ctx, userId, title)
			if err != nil {
				return common.DefaultErrorRedirect(err.Error())
			}
			return postUrlBuilder(common.GetBaseUrl(1, c), postId).String()
		}),
		deleteHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			logger := web.GetLogger(c)
			var targetBuilder strings.Builder
			targetBuilder.WriteString(common.GetBaseUrl(2, c))

			postId, err := strconv.ParseUint(c.Param(postIdName), 10, 64)
			if err != nil {
				logger.Warn(parsingPostIdErrorMsg, common.ErrorKey, err)
				common.WriteError(&targetBuilder, common.ErrorTechnicalKey)
				return targetBuilder.String()
			}
			userId := web.GetSessionUserId(c)

			post, err := blogService.GetPost(ctx, userId, postId)
			if err != nil {
				common.WriteError(&targetBuilder, err.Error())
				return targetBuilder.String()
			}

			if err = blogService.DeletePost(ctx, userId, postId); err != nil {
				common.WriteError(&targetBuilder, err.Error())
				return targetBuilder.String()
			}

			if err = commentService.DeleteCommentThread(ctx, userId, post.Title); err != nil {
				common.WriteError(&targetBuilder, err.Error())
			}
			return targetBuilder.String()
		}),
		rssHandler: func(c *gin.Context) {
			ctx := c.Request.Context()
			logger := web.GetLogger(c)
			userId := web.GetSessionUserId(c)

			_, posts, err := blogService.GetPosts(ctx, userId, 0, feedSize, "")
			if err != nil {
				status := http.StatusInternalServerError
				if err == common.ErrNotAuthorized {
					status = http.StatusForbidden
				}
				c.AbortWithStatus(status)
				return
			}

			baseUrl := host + common.GetBaseUrl(1, c)
			// TODO improve blog title ?
			data, err := buildFeed(posts, blogName, baseUrl, dateFormat, extractSize, feedFormat)
			if err != nil {
				common.LogOriginalError(logger, err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Data(http.StatusOK, http.DetectContentType(data), data)
		},
	}
	return p
}

func postUrlBuilder(base string, postId uint64) *strings.Builder {
	targetBuilder := new(strings.Builder)
	targetBuilder.WriteString(base)
	targetBuilder.WriteString("view/")
	targetBuilder.WriteString(strconv.FormatUint(postId, 10))
	return targetBuilder
}

func filterPostsExtract(posts []blogservice.BlogPost, extractSize uint64) {
	for index := range posts {
		posts[index].Content = common.FilterExtractHtml(string(posts[index].Content), extractSize)
	}
}

func buildFeed(posts []blogservice.BlogPost, blogTitle string, baseUrl string, dateFormat string, extractSize uint64, feedFormat string) ([]byte, error) {
	feedData := feeds.Feed{
		Title:   blogTitle,
		Link:    &feeds.Link{Href: baseUrl},
		Created: time.Now(),
	}

	for _, post := range posts {
		date, err := time.Parse(dateFormat, post.Date)
		if err != nil {
			return nil, err
		}

		feedData.Items = append(feedData.Items, &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: postUrlBuilder(baseUrl, post.PostId).String()},
			Description: common.FilterExtractHtml(string(post.Content), extractSize),
			Author:      &feeds.Author{Name: post.Creator.Login},
			Created:     date,
		})
	}

	data := ""
	var err error
	switch feedFormat {
	case "atom":
		data, err = feedData.ToAtom()
	case "json":
		data, err = feedData.ToJSON()
	case "rss":
		data, err = feedData.ToRss()
	default:
		return nil, errFeedFormat
	}
	return []byte(data), err
}
