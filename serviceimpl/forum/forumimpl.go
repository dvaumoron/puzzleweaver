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

	"github.com/ServiceWeaver/weaver"
	dbclient "github.com/dvaumoron/puzzleweaver/client/db"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/serviceimpl/forum/model"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"gorm.io/gorm"
)

type RemoteForumService remoteservice.RemoteForumService

type remoteForumImpl struct {
	weaver.Implements[RemoteForumService]
	weaver.WithConfig[forumConf]
	initializedConf initializedForumConf
}

func (impl *remoteForumImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf, err = initForumConf(ctx, impl.Config())
	return
}

func (impl *remoteForumImpl) CreateThread(ctx context.Context, objectId uint64, userId uint64, title string, message string) (uint64, error) {
	db := impl.initializedConf.db.WithContext(ctx)
	thread := model.Thread{
		ObjectId: objectId, UserId: userId, Title: title,
	}
	if message != "" {
		thread.Messages = []model.Message{{UserId: userId, Text: message}}
	}
	if err := db.Create(&thread).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, common.ErrUpdate
	}
	return thread.ID, nil
}

func (impl *remoteForumImpl) CreateMessage(ctx context.Context, objectId uint64, userId uint64, threadId uint64, message string) error {
	db := impl.initializedConf.db.WithContext(ctx)
	mMessage := model.Message{ThreadID: threadId, UserId: userId, Text: message}
	if err := db.Create(&mMessage).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *remoteForumImpl) GetThread(ctx context.Context, objectId uint64, threadId uint64, start uint64, end uint64, filter string) (uint64, remoteservice.RawForumContent, []remoteservice.RawForumContent, error) {
	db := impl.initializedConf.db.WithContext(ctx)

	paginate := func(tx *gorm.DB) *gorm.DB {
		return dbclient.Paginate(tx, start, end).Order("created_at asc")
	}

	preloadFilter := paginate
	messageRequest := db.Model(&model.Message{})
	if filter == "" {
		messageRequest.Where("thread_id = ?", threadId)
	} else {
		filter = dbclient.BuildLikeFilter(filter)
		messageRequest.Where("thread_id = ? AND text LIKE ?", threadId, filter)

		preloadFilter = func(tx *gorm.DB) *gorm.DB {
			return paginate(tx).Where("text LIKE ?", filter)
		}
	}

	var total int64
	err := messageRequest.Count(&total).Error
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, remoteservice.RawForumContent{}, nil, servicecommon.ErrInternal
	}

	var thread model.Thread
	if err := db.Preload("Messages", preloadFilter).First(&thread, threadId).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, remoteservice.RawForumContent{}, nil, servicecommon.ErrInternal
	}
	return uint64(total), convertThreadFromModel(thread), convertMessagesFromModel(thread.Messages), nil
}

func (impl *remoteForumImpl) GetThreads(ctx context.Context, objectId uint64, start uint64, end uint64, filter string) (uint64, []remoteservice.RawForumContent, error) {
	db := impl.initializedConf.db.WithContext(ctx)
	noFilter := filter == ""

	threadRequest := db.Model(&model.Thread{})
	if noFilter {
		threadRequest.Where("object_id = ?", objectId)
	} else {
		filter = dbclient.BuildLikeFilter(filter)
		threadRequest.Where("object_id = ? AND title LIKE ?", objectId, filter)
	}

	var total int64
	err := threadRequest.Count(&total).Error
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}
	if total == 0 {
		return 0, nil, nil
	}

	var threads []model.Thread
	page := dbclient.Paginate(db, start, end).Order("created_at desc")
	if noFilter {
		err = page.Find(&threads, "object_id = ?", objectId).Error
	} else {
		err = page.Find(&threads, "object_id = ? AND title LIKE ?", objectId, filter).Error
	}

	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}
	return uint64(total), convertThreadsFromModel(threads), nil
}

func (impl *remoteForumImpl) DeleteThread(ctx context.Context, containerId uint64, id uint64) error {
	db := impl.initializedConf.db.WithContext(ctx)
	if err := db.Delete(&model.Thread{}, id).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *remoteForumImpl) DeleteMessage(ctx context.Context, containerId uint64, id uint64) error {
	db := impl.initializedConf.db.WithContext(ctx)
	if err := db.Delete(&model.Message{}, id).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func convertThreadFromModel(thread model.Thread) remoteservice.RawForumContent {
	return remoteservice.RawForumContent{
		Id: thread.ID, CreatedAt: thread.CreatedAt.Unix(), CreatorId: thread.UserId, Text: thread.Title,
	}
}

func convertThreadsFromModel(threads []model.Thread) []remoteservice.RawForumContent {
	resThreads := make([]remoteservice.RawForumContent, 0, len(threads))
	for _, thread := range threads {
		resThreads = append(resThreads, convertThreadFromModel(thread))
	}
	return resThreads
}

func convertMessagesFromModel(messages []model.Message) []remoteservice.RawForumContent {
	resMessages := make([]remoteservice.RawForumContent, 0, len(messages))
	for _, message := range messages {
		resMessages = append(resMessages, remoteservice.RawForumContent{
			Id: message.ID, CreatedAt: message.CreatedAt.Unix(), CreatorId: message.UserId, Text: message.Text,
		})
	}
	return resMessages
}
