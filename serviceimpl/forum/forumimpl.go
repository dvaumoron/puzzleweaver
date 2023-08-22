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
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
)

type RemoteForumService remoteservice.RemoteForumService

type remoteForumImpl struct {
	weaver.Implements[RemoteForumService]
}

func (*remoteForumImpl) CreateThread(ctx context.Context, objectId uint64, userId uint64, title string, message string) (uint64, error) {
	id := uint64(0)
	success := true
	// TODO
	if !success {
		return 0, common.ErrUpdate
	}
	return id, nil
}

func (*remoteForumImpl) CreateMessage(ctx context.Context, objectId uint64, userId uint64, threadId uint64, message string) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (*remoteForumImpl) GetThread(ctx context.Context, objectId uint64, threadId uint64, start uint64, end uint64, filter string) (uint64, remoteservice.RawForumContent, []remoteservice.RawForumContent, error) {
	// TODO
	return 0, remoteservice.RawForumContent{}, nil, nil
}

func (*remoteForumImpl) GetThreads(ctx context.Context, objectId uint64, start uint64, end uint64, filter string) (uint64, []remoteservice.RawForumContent, error) {
	// TODO
	return 0, nil, nil
}

func (*remoteForumImpl) DeleteThread(ctx context.Context, containerId uint64, id uint64) error {
	// TODO
	return nil
}

func (*remoteForumImpl) DeleteMessage(ctx context.Context, containerId uint64, id uint64) error {
	// TODO
	return nil
}
