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

package blogclient

import (
	"cmp"
	"context"
	"slices"
	"time"

	adminimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/admin"
	blogimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/blog"
	blogservice "github.com/dvaumoron/puzzleweb/blog/service"
	profileservice "github.com/dvaumoron/puzzleweb/profile/service"
)

type blogServiceWrapper struct {
	blogService    blogimpl.RemoteBlogService
	authService    adminimpl.AuthService
	profileService profileservice.ProfileService
	blogId         uint64
	groupId        uint64
	dateFormat     string
}

func MakeBlogServiceWrapper(blogService blogimpl.RemoteBlogService, authService adminimpl.AuthService, profileService profileservice.ProfileService, blogId uint64, groupId uint64, dateFormat string) blogservice.BlogService {
	return blogServiceWrapper{
		blogService: blogService, authService: authService, profileService: profileService,
		blogId: blogId, groupId: groupId, dateFormat: dateFormat,
	}
}

func cmpDesc(a blogimpl.RawBlogPost, b blogimpl.RawBlogPost) int {
	return cmp.Compare(b.CreatedAt, a.CreatedAt)
}

func (client blogServiceWrapper) CreatePost(ctx context.Context, userId uint64, title string, content string) (uint64, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionCreate)
	if err != nil {
		return 0, err
	}
	return client.blogService.CreatePost(ctx, client.blogId, userId, title, content)
}

func (client blogServiceWrapper) GetPost(ctx context.Context, userId uint64, postId uint64) (blogservice.BlogPost, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionAccess)
	if err != nil {
		return blogservice.BlogPost{}, err
	}

	post, err := client.blogService.GetPost(ctx, client.blogId, postId)
	if err != nil {
		return blogservice.BlogPost{}, err
	}

	creatorId := post.CreatorId
	users, err := client.profileService.GetProfiles(ctx, []uint64{creatorId})
	if err != nil {
		return blogservice.BlogPost{}, err
	}
	return convertPost(post, users[creatorId], client.dateFormat), nil
}

func (client blogServiceWrapper) GetPosts(ctx context.Context, userId uint64, start uint64, end uint64, filter string) (uint64, []blogservice.BlogPost, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionAccess)
	if err != nil {
		return 0, nil, err
	}

	total, list, err := client.blogService.GetPosts(ctx, client.blogId, start, end, filter)
	if len(list) == 0 {
		return total, nil, nil
	}

	size := len(list)
	slices.SortFunc(list, cmpDesc)
	// no duplicate check, there is one in GetProfiles
	userIds := make([]uint64, 0, size)
	for _, content := range list {
		userIds = append(userIds, content.CreatorId)
	}

	users, err := client.profileService.GetProfiles(ctx, userIds)
	if err != nil {
		return 0, nil, err
	}

	posts := make([]blogservice.BlogPost, 0, size)
	for _, content := range list {
		posts = append(posts, convertPost(content, users[content.CreatorId], client.dateFormat))
	}
	return total, posts, nil
}

func (client blogServiceWrapper) DeletePost(ctx context.Context, userId uint64, postId uint64) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionDelete)
	if err != nil {
		return err
	}
	return client.blogService.Delete(ctx, client.blogId, postId)
}

func (client blogServiceWrapper) CreateRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionCreate) == nil
}

func (client blogServiceWrapper) DeleteRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionDelete) == nil
}

func convertPost(post blogimpl.RawBlogPost, creator profileservice.UserProfile, dateFormat string) blogservice.BlogPost {
	createdAt := time.Unix(post.CreatedAt, 0)
	return blogservice.BlogPost{
		PostId: post.Id, Creator: creator, Date: createdAt.Format(dateFormat), Title: post.Title, Content: post.Content,
	}
}
