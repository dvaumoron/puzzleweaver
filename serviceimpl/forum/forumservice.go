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
)

type RawForumContent struct {
	weaver.AutoMarshal
	Id        uint64
	CreatorId uint64
	CreatedAt int64
	Text      string
}

type RemoteForumService interface {
	CreateThread(ctx context.Context, objectId uint64, userId uint64, title string, message string) (uint64, error)
	CreateMessage(ctx context.Context, objectId uint64, userId uint64, threadId uint64, message string) error
	GetThread(ctx context.Context, objectId uint64, threadId uint64, start uint64, end uint64, filter string) (uint64, RawForumContent, []RawForumContent, error)
	GetThreads(ctx context.Context, objectId uint64, start uint64, end uint64, filter string) (uint64, []RawForumContent, error)
	DeleteThread(ctx context.Context, containerId uint64, id uint64) error
	DeleteMessage(ctx context.Context, containerId uint64, id uint64) error
}
