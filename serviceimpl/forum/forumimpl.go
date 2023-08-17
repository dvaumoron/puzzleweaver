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

package forumimpl

import (
	"context"
	"sort"
	"time"

	"github.com/ServiceWeaver/weaver"
	pb "github.com/dvaumoron/puzzleforumservice"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
	"github.com/dvaumoron/puzzleweb/common"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"golang.org/x/exp/slog"
)

// check matching with interface
var _ forumservice.FullForumService = &forumImpl{}

type forumImpl struct {
	weaver.Implements[forumservice.FullForumService]
	authService    weaver.Ref[service.AuthService]
	profileService weaver.Ref[service.ProfileService]
	forumId        uint64
	groupId        uint64
	dateFormat     string
}

type deleteRequestKind func(pb.ForumClient, otelzap.LoggerWithCtx, *pb.IdRequest) (*pb.Response, error)

type sortableContentsDesc []*pb.Content

func (s sortableContentsDesc) Len() int {
	return len(s)
}

func (s sortableContentsDesc) Less(i, j int) bool {
	return s[i].CreatedAt > s[j].CreatedAt
}

func (s sortableContentsDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortableContentsAsc []*pb.Content

func (s sortableContentsAsc) Len() int {
	return len(s)
}

func (s sortableContentsAsc) Less(i, j int) bool {
	return s[i].CreatedAt < s[j].CreatedAt
}

func (s sortableContentsAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (impl forumImpl) CreateThread(ctx context.Context, userId uint64, title string, message string) (uint64, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionCreate)
	if err != nil {
		return 0, err
	}

	id := uint64(0)
	success := true
	// TODO
	if !success {
		return 0, common.ErrUpdate
	}
	return id, nil
}

func (impl forumImpl) CreateCommentThread(ctx context.Context, userId uint64, elemTitle string) error {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionCreate)
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

func (impl forumImpl) CreateMessage(ctx context.Context, userId uint64, threadId uint64, message string) error {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionUpdate)
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

func (impl forumImpl) CreateComment(ctx context.Context, userId uint64, elemTitle string, comment string) error {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionAccess)
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

func (impl forumImpl) GetThread(ctx context.Context, userId uint64, threadId uint64, start uint64, end uint64, filter string) (uint64, forumservice.ForumContent, []forumservice.ForumContent, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionAccess)
	if err != nil {
		return 0, forumservice.ForumContent{}, nil, err
	}

	// TODO
	total := uint64(0)
	list := []*pb.Content{}
	userIds := extractUserIds(list)
	threadCreatorId := uint64(0)
	userIds = append(userIds, threadCreatorId)

	users, err := impl.profileService.Get().GetProfiles(ctx, userIds)
	if err != nil {
		return 0, forumservice.ForumContent{}, nil, err
	}

	thread := convertContent(nil, users[threadCreatorId], impl.dateFormat)
	sort.Sort(sortableContentsAsc(list))
	messages := convertContents(list, users, impl.dateFormat)
	return total, thread, messages, nil
}

func (impl forumImpl) GetThreads(ctx context.Context, userId uint64, start uint64, end uint64, filter string) (uint64, []forumservice.ForumContent, error) {
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

	users, err := impl.profileService.Get().GetProfiles(ctx, extractUserIds(list))
	if err != nil {
		return 0, nil, err
	}
	sort.Sort(sortableContentsDesc(list))
	return total, convertContents(list, users, impl.dateFormat), nil
}

func (impl forumImpl) GetCommentThread(ctx context.Context, userId uint64, elemTitle string, start uint64, end uint64) (uint64, []forumservice.ForumContent, error) {
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

	users, err := impl.profileService.Get().GetProfiles(ctx, extractUserIds(list))
	if err != nil {
		return 0, nil, err
	}
	sort.Sort(sortableContentsAsc(list))
	return total, convertContents(list, users, impl.dateFormat), nil
}

func (impl forumImpl) DeleteThread(ctx context.Context, userId uint64, threadId uint64) error {
	return impl.deleteContent(ctx, userId, deleteThread, &pb.IdRequest{ContainerId: impl.forumId, Id: threadId})
}

func (impl forumImpl) DeleteCommentThread(ctx context.Context, userId uint64, elemTitle string) error {
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

func (impl forumImpl) DeleteMessage(ctx context.Context, userId uint64, threadId uint64, messageId uint64) error {
	return impl.deleteContent(
		ctx, userId, deleteMessage, &pb.IdRequest{ContainerId: threadId, Id: messageId},
	)
}

func (impl forumImpl) DeleteComment(ctx context.Context, userId uint64, elemTitle string, commentId uint64) error {
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

func (impl forumImpl) CreateThreadRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionCreate)
}

func (impl forumImpl) CreateMessageRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionUpdate)
}

func (impl forumImpl) DeleteRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionDelete)
}

func (impl forumImpl) deleteContent(ctx context.Context, userId uint64, kind deleteRequestKind, request *pb.IdRequest) error {
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

func searchCommentThread(forumimpl pb.ForumClient, logger otelzap.LoggerWithCtx, objectId uint64, elemTitle string) (*pb.Contents, error) {
	return forumimpl.GetThreads(logger.Context(), &pb.SearchRequest{
		ContainerId: objectId, Start: 0, End: 1, Filter: elemTitle,
	})
}

func deleteThread(forumimpl pb.ForumClient, logger otelzap.LoggerWithCtx, request *pb.IdRequest) (*pb.Response, error) {
	return forumimpl.DeleteThread(logger.Context(), request)
}

func deleteMessage(forumimpl pb.ForumClient, logger otelzap.LoggerWithCtx, request *pb.IdRequest) (*pb.Response, error) {
	return forumimpl.DeleteMessage(logger.Context(), request)
}

func convertContents(list []*pb.Content, users map[uint64]service.UserProfile, dateFormat string) []forumservice.ForumContent {
	contents := make([]forumservice.ForumContent, 0, len(list))
	for _, content := range list {
		contents = append(contents, convertContent(content, users[content.UserId], dateFormat))
	}
	return contents
}

func convertContent(content *pb.Content, creator service.UserProfile, dateFormat string) forumservice.ForumContent {
	createdAt := time.Unix(content.CreatedAt, 0)
	return forumservice.ForumContent{
		Id: content.Id, Creator: creator, Date: createdAt.Format(dateFormat), Text: content.Text,
	}
}

func logCommentThreadNotFound(logger *slog.Logger, objectId uint64, elemTitle string) error {
	logger.Warn("comment thread not found", "objectId", objectId, "elemTitle", elemTitle)
	return common.ErrTechnical
}

// no duplicate check, there is one in GetProfiles
func extractUserIds(list []*pb.Content) []uint64 {
	userIds := make([]uint64, 0, len(list))
	for _, content := range list {
		userIds = append(userIds, content.UserId)
	}
	return userIds
}
