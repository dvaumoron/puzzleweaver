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

package blogimpl

import (
	"context"
	"sort"
	"time"

	"github.com/ServiceWeaver/weaver"
	pb "github.com/dvaumoron/puzzleblogservice"
	blogservice "github.com/dvaumoron/puzzleweaver/web/blog/service"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweb/common"
)

// check matching with interface
var _ blogservice.BlogService = &blogImpl{}

type blogImpl struct {
	weaver.Implements[blogservice.BlogService]
	authService    weaver.Ref[service.AdminService]
	profileService weaver.Ref[service.AdvancedProfileService]
	blogId         uint64
	groupId        uint64
	dateFormat     string
}

type sortableContents []*pb.Content

func (s sortableContents) Len() int {
	return len(s)
}

func (s sortableContents) Less(i, j int) bool {
	return s[i].CreatedAt > s[j].CreatedAt
}

func (s sortableContents) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (impl blogImpl) CreatePost(ctx context.Context, userId uint64, title string, content string) (uint64, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionCreate)
	if err != nil {
		return 0, err
	}

	success := true
	// TODO
	if !success {
		return 0, common.ErrUpdate
	}
	return 0, nil
}

func (impl blogImpl) GetPost(ctx context.Context, userId uint64, postId uint64) (blogservice.BlogPost, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionAccess)
	if err != nil {
		return blogservice.BlogPost{}, err
	}

	// TODO
	creatorId := uint64(0)
	users, err := impl.profileService.Get().GetProfiles(ctx, []uint64{creatorId})
	if err != nil {
		return blogservice.BlogPost{}, err
	}
	return convertPost(nil, users[creatorId], impl.dateFormat), nil
}

func (impl blogImpl) GetPosts(ctx context.Context, userId uint64, start uint64, end uint64, filter string) (uint64, []blogservice.BlogPost, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionAccess)
	if err != nil {
		return 0, nil, err
	}

	// TODO
	total := uint64(0)
	list := []*pb.Content{}
	if len(list) == 0 {
		return total, nil, nil
	}

	posts, err := impl.sortConvertPosts(ctx, list)
	if err != nil {
		return 0, nil, err
	}
	return total, posts, nil
}

func (impl blogImpl) DeletePost(ctx context.Context, userId uint64, postId uint64) error {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionDelete)
	if err != nil {
		return err
	}

	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl blogImpl) CreateRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionCreate)
}

func (impl blogImpl) DeleteRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionDelete)
}

func (impl blogImpl) sortConvertPosts(ctx context.Context, list []*pb.Content) ([]blogservice.BlogPost, error) {
	sort.Sort(sortableContents(list))

	size := len(list)
	// no duplicate check, there is one in GetProfiles
	userIds := make([]uint64, 0, size)
	for _, content := range list {
		userIds = append(userIds, content.UserId)
	}

	users, err := impl.profileService.Get().GetProfiles(ctx, userIds)
	if err != nil {
		return nil, err
	}

	contents := make([]blogservice.BlogPost, 0, size)
	for _, content := range list {
		contents = append(contents, convertPost(content, users[content.UserId], impl.dateFormat))
	}
	return contents, nil
}

func convertPost(post *pb.Content, creator service.UserProfile, dateFormat string) blogservice.BlogPost {
	createdAt := time.Unix(post.CreatedAt, 0)
	return blogservice.BlogPost{
		PostId: post.PostId, Creator: creator, Date: createdAt.Format(dateFormat), Title: post.Title, Content: post.Text,
	}
}
