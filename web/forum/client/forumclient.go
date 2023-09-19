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
	"cmp"
	"context"
	"log/slog"
	"slices"
	"time"

	adminimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/admin"
	forumimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/forum"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	forumservice "github.com/dvaumoron/puzzleweaver/web/forum/service"
)

type forumServiceWrapper struct {
	forumService   forumimpl.RemoteForumService
	authService    adminimpl.AuthService
	profileService service.ProfileService
	loggerGetter   common.LoggerGetter
	forumId        uint64
	groupId        uint64
	dateFormat     string
}

func MakeForumServiceWrapper(forumService forumimpl.RemoteForumService, authService adminimpl.AuthService, profileService service.ProfileService, loggerGetter common.LoggerGetter, forumId uint64, groupId uint64, dateFormat string) forumservice.FullForumService {
	return forumServiceWrapper{
		forumService: forumService, authService: authService, profileService: profileService,
		loggerGetter: loggerGetter, forumId: forumId, groupId: groupId, dateFormat: dateFormat,
	}
}

type deleteRequestKind func(forumimpl.RemoteForumService, context.Context, uint64, uint64) error

func cmpAsc(a forumimpl.RawForumContent, b forumimpl.RawForumContent) int {
	return cmp.Compare(a.CreatedAt, b.CreatedAt)
}

func cmpDesc(a forumimpl.RawForumContent, b forumimpl.RawForumContent) int {
	return -cmp.Compare(a.CreatedAt, b.CreatedAt)
}

func (client forumServiceWrapper) CreateThread(ctx context.Context, userId uint64, title string, message string) (uint64, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionCreate)
	if err != nil {
		return 0, err
	}
	return client.forumService.CreateThread(ctx, client.forumId, userId, title, message)
}

func (client forumServiceWrapper) CreateCommentThread(ctx context.Context, userId uint64, elemTitle string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionCreate)
	if err != nil {
		return err
	}
	_, err = client.forumService.CreateThread(ctx, client.forumId, userId, elemTitle, "")
	return err
}

func (client forumServiceWrapper) CreateMessage(ctx context.Context, userId uint64, threadId uint64, message string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionUpdate)
	if err != nil {
		return err
	}
	return client.forumService.CreateMessage(ctx, client.forumId, userId, threadId, message)
}

func (client forumServiceWrapper) CreateComment(ctx context.Context, userId uint64, elemTitle string, comment string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionAccess)
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
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionAccess)
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
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionAccess)
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
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionAccess)
	if err != nil {
		return 0, nil, err
	}

	_, threads, err := client.forumService.GetThreads(ctx, client.forumId, 0, 1, elemTitle)
	if err != nil {
		logCommentThreadNotFound(client.loggerGetter.Logger(ctx), client.forumId, elemTitle)
		return 0, nil, err
	}
	if len(threads) == 0 {
		logCommentThreadNotFound(client.loggerGetter.Logger(ctx), client.forumId, elemTitle)
		return 0, nil, common.ErrTechnical
	}

	threadId := threads[0].Id
	total, _, messages, err := client.forumService.GetThread(ctx, client.forumId, threadId, start, end, "")
	if err != nil {
		return 0, nil, err
	}

	users, err := client.profileService.GetProfiles(ctx, extractUserIds(messages))
	if err != nil {
		return 0, nil, err
	}
	slices.SortFunc(messages, cmpAsc)
	return total, convertContents(messages, users, client.dateFormat), nil
}

func (client forumServiceWrapper) DeleteThread(ctx context.Context, userId uint64, threadId uint64) error {
	return client.deleteContent(ctx, userId, forumimpl.RemoteForumService.DeleteThread, client.forumId, threadId)
}

func (client forumServiceWrapper) DeleteCommentThread(ctx context.Context, userId uint64, elemTitle string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionDelete)
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
	return client.deleteContent(ctx, userId, forumimpl.RemoteForumService.DeleteMessage, threadId, messageId)
}

func (client forumServiceWrapper) DeleteComment(ctx context.Context, userId uint64, elemTitle string, commentId uint64) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionDelete)
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
	return client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionCreate) == nil
}

func (client forumServiceWrapper) CreateMessageRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionUpdate) == nil
}

func (client forumServiceWrapper) DeleteRight(ctx context.Context, userId uint64) bool {
	return client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionDelete) == nil
}

func (client forumServiceWrapper) deleteContent(ctx context.Context, userId uint64, kind deleteRequestKind, containerId uint64, id uint64) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, adminimpl.ActionDelete)
	if err != nil {
		return err
	}
	return kind(client.forumService, ctx, containerId, id)
}

func deleteThread(forumService forumimpl.RemoteForumService, ctx context.Context, containerId uint64, id uint64) error {
	return forumService.DeleteThread(ctx, containerId, id)
}

func deleteMessage(forumService forumimpl.RemoteForumService, ctx context.Context, containerId uint64, id uint64) error {
	return forumService.DeleteMessage(ctx, containerId, id)
}

func convertContents(list []forumimpl.RawForumContent, users map[uint64]service.UserProfile, dateFormat string) []forumservice.ForumContent {
	contents := make([]forumservice.ForumContent, 0, len(list))
	for _, content := range list {
		contents = append(contents, convertContent(content, users[content.CreatorId], dateFormat))
	}
	return contents
}

func convertContent(content forumimpl.RawForumContent, creator service.UserProfile, dateFormat string) forumservice.ForumContent {
	createdAt := time.Unix(content.CreatedAt, 0)
	return forumservice.ForumContent{
		Id: content.Id, Creator: creator, Date: createdAt.Format(dateFormat), Text: content.Text,
	}
}

func logCommentThreadNotFound(logger *slog.Logger, objectId uint64, elemTitle string) {
	logger.Warn("comment thread not found", "objectId", objectId, "elemTitle", elemTitle)
}

// no duplicate check, there is one in GetProfiles
func extractUserIds(list []forumimpl.RawForumContent) []uint64 {
	userIds := make([]uint64, 0, len(list))
	for _, content := range list {
		userIds = append(userIds, content.CreatorId)
	}
	return userIds
}
