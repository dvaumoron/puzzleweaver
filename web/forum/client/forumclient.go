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

package forumclient

import (
	"context"
	"time"

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

type forumServiceWrapper struct {
	forumService   remoteservice.RemoteForumService
	authService    service.AuthService
	profileService service.ProfileService
	loggerGetter   common.LoggerGetter
	forumId        uint64
	groupId        uint64
	dateFormat     string
}

func MakeForumServiceWrapper(forumService remoteservice.RemoteForumService, authService service.AuthService, profileService service.ProfileService, loggerGetter common.LoggerGetter, forumId uint64, groupId uint64, dateFormat string) forumservice.FullForumService {
	return forumServiceWrapper{
		forumService: forumService, authService: authService, profileService: profileService,
		loggerGetter: loggerGetter, forumId: forumId, groupId: groupId, dateFormat: dateFormat,
	}
}

type deleteRequestKind func(remoteservice.RemoteForumService, context.Context, uint64, uint64) error

func cmpAsc(a remoteservice.RawForumContent, b remoteservice.RawForumContent) bool {
	return a.CreatedAt < b.CreatedAt
}

func cmpDesc(a remoteservice.RawForumContent, b remoteservice.RawForumContent) bool {
	return a.CreatedAt > b.CreatedAt
}

func (client forumServiceWrapper) CreateThread(ctx context.Context, userId uint64, title string, message string) (uint64, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionCreate)
	if err != nil {
		return 0, err
	}
	return client.forumService.CreateThread(ctx, client.forumId, userId, title, message)
}

func (client forumServiceWrapper) CreateCommentThread(ctx context.Context, userId uint64, elemTitle string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionCreate)
	if err != nil {
		return err
	}
	_, err = client.forumService.CreateThread(ctx, client.forumId, userId, elemTitle, "")
	return err
}

func (client forumServiceWrapper) CreateMessage(ctx context.Context, userId uint64, threadId uint64, message string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionUpdate)
	if err != nil {
		return err
	}
	return client.forumService.CreateMessage(ctx, client.forumId, userId, threadId, message)
}

func (client forumServiceWrapper) CreateComment(ctx context.Context, userId uint64, elemTitle string, comment string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return err
	}

	total, threads, err := client.forumService.GetThreads(ctx, client.forumId, 0, 1, elemTitle)
	if err != nil {
		return err
	}

	if total == 0 {
		logCommentThreadNotFound(client.loggerGetter.Logger(ctx), client.forumId, elemTitle)

		_, err = client.forumService.CreateThread(ctx, client.forumId, userId, elemTitle, comment)
		return err
	}

	threadId := threads[0].Id
	return client.forumService.CreateMessage(ctx, client.forumId, userId, threadId, comment)
}

func (client forumServiceWrapper) GetThread(ctx context.Context, userId uint64, threadId uint64, start uint64, end uint64, filter string) (uint64, forumservice.ForumContent, []forumservice.ForumContent, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return 0, forumservice.ForumContent{}, nil, err
	}

	total, rawThread, list, err := client.forumService.GetThread(ctx, client.forumId, threadId, start, end, filter)
	if err != nil {
		return 0, forumservice.ForumContent{}, nil, err
	}

	userIds := extractUserIds(list)
	threadCreatorId := uint64(0)
	userIds = append(userIds, threadCreatorId)

	users, err := client.profileService.GetProfiles(ctx, userIds)
	if err != nil {
		return 0, forumservice.ForumContent{}, nil, err
	}

	thread := convertContent(rawThread, users[threadCreatorId], client.dateFormat)
	slices.SortFunc(list, cmpAsc)
	messages := convertContents(list, users, client.dateFormat)
	return total, thread, messages, nil
}

func (client forumServiceWrapper) GetThreads(ctx context.Context, userId uint64, start uint64, end uint64, filter string) (uint64, []forumservice.ForumContent, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return 0, nil, err
	}

	total, message, err := client.forumService.GetThreads(ctx, client.forumId, start, end, filter)
	if err != nil {
		return 0, nil, err
	}
	if len(message) == 0 {
		return total, nil, nil
	}

	users, err := client.profileService.GetProfiles(ctx, extractUserIds(message))
	if err != nil {
		return 0, nil, err
	}
	slices.SortFunc(message, cmpDesc)
	return total, convertContents(message, users, client.dateFormat), nil
}

func (client forumServiceWrapper) GetCommentThread(ctx context.Context, userId uint64, elemTitle string, start uint64, end uint64) (uint64, []forumservice.ForumContent, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return 0, nil, err
	}

	total, threads, err := client.forumService.GetThreads(ctx, client.forumId, 0, 1, elemTitle)
	if err != nil {
		logCommentThreadNotFound(client.loggerGetter.Logger(ctx), client.forumId, elemTitle)
		return 0, nil, err
	}
	if len(threads) == 0 {
		return total, nil, nil
	}

	users, err := client.profileService.GetProfiles(ctx, extractUserIds(threads))
	if err != nil {
		return 0, nil, err
	}
	slices.SortFunc(threads, cmpAsc)
	return total, convertContents(threads, users, client.dateFormat), nil
}

func (client forumServiceWrapper) DeleteThread(ctx context.Context, userId uint64, threadId uint64) error {
	return client.deleteContent(ctx, userId, remoteservice.RemoteForumService.DeleteThread, client.forumId, threadId)
}

func (client forumServiceWrapper) DeleteCommentThread(ctx context.Context, userId uint64, elemTitle string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionDelete)
	if err != nil {
		return err
	}

	_, threads, err := client.forumService.GetThreads(ctx, client.forumId, 0, 1, elemTitle)
	if err != nil {
		logCommentThreadNotFound(client.loggerGetter.Logger(ctx), client.forumId, elemTitle)
		return err
	}
	if len(threads) == 0 {
		return nil
	}
	return client.forumService.DeleteThread(ctx, client.forumId, threads[0].Id)
}

func (client forumServiceWrapper) DeleteMessage(ctx context.Context, userId uint64, threadId uint64, messageId uint64) error {
	return client.deleteContent(ctx, userId, remoteservice.RemoteForumService.DeleteMessage, threadId, messageId)
}

func (client forumServiceWrapper) DeleteComment(ctx context.Context, userId uint64, elemTitle string, commentId uint64) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionDelete)
	if err != nil {
		return err
	}

	_, threads, err := client.forumService.GetThreads(ctx, client.forumId, 0, 1, elemTitle)
	if err != nil {
		logCommentThreadNotFound(client.loggerGetter.Logger(ctx), client.forumId, elemTitle)
		return err
	}
	if len(threads) == 0 {
		return nil
	}
	return client.forumService.DeleteMessage(ctx, threads[0].Id, commentId)
}

func (client forumServiceWrapper) CreateThreadRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionCreate) == nil
}

func (client forumServiceWrapper) CreateMessageRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionUpdate) == nil
}

func (client forumServiceWrapper) DeleteRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionDelete) == nil
}

func (client forumServiceWrapper) deleteContent(ctx context.Context, userId uint64, kind deleteRequestKind, containerId uint64, id uint64) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionDelete)
	if err != nil {
		return err
	}
	return kind(client.forumService, ctx, containerId, id)
}

func deleteThread(forumService remoteservice.RemoteForumService, ctx context.Context, containerId uint64, id uint64) error {
	return forumService.DeleteThread(ctx, containerId, id)
}

func deleteMessage(forumService remoteservice.RemoteForumService, ctx context.Context, containerId uint64, id uint64) error {
	return forumService.DeleteMessage(ctx, containerId, id)
}

func convertContents(list []remoteservice.RawForumContent, users map[uint64]service.UserProfile, dateFormat string) []forumservice.ForumContent {
	contents := make([]forumservice.ForumContent, 0, len(list))
	for _, content := range list {
		contents = append(contents, convertContent(content, users[content.CreatorId], dateFormat))
	}
	return contents
}

func convertContent(content remoteservice.RawForumContent, creator service.UserProfile, dateFormat string) forumservice.ForumContent {
	createdAt := time.Unix(content.CreatedAt, 0)
	return forumservice.ForumContent{
		Id: content.Id, Creator: creator, Date: createdAt.Format(dateFormat), Text: content.Text,
	}
}

func logCommentThreadNotFound(logger *slog.Logger, objectId uint64, elemTitle string) {
	logger.Warn("comment thread not found", "objectId", objectId, "elemTitle", elemTitle)
}

// no duplicate check, there is one in GetProfiles
func extractUserIds(list []remoteservice.RawForumContent) []uint64 {
	userIds := make([]uint64, 0, len(list))
	for _, content := range list {
		userIds = append(userIds, content.CreatorId)
	}
	return userIds
}
