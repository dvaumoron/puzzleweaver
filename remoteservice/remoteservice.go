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

package remoteservice

import (
	"context"
)

type RawBlogPost struct {
	Id        uint64
	CreatorId uint64
	CreatedAt int64
	Title     string
	Content   string
}

type RemoteBlogService interface {
	CreatePost(ctx context.Context, blogId uint64, userId uint64, title string, content string) (uint64, error)
	GetPost(ctx context.Context, blogId uint64, postId uint64) (RawBlogPost, error)
	GetPosts(ctx context.Context, blogId uint64, start uint64, end uint64, filter string) (uint64, []RawBlogPost, error)
	Delete(ctx context.Context, blogId uint64, postId uint64) error
}

type RawForumContent struct {
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

type RawUserProfile struct {
	Desc string
	Info map[string]string
}

type RemoteProfileService interface {
	GetProfiles(ctx context.Context, userIds []uint64) (map[uint64]RawUserProfile, error)
	GetPicture(ctx context.Context, userId uint64) ([]byte, error)
	UpdateProfile(ctx context.Context, userId uint64, desc string, info map[string]string) error
	UpdatePicture(ctx context.Context, userId uint64, data []byte) error
	Delete(ctx context.Context, userId uint64) error
}

type RawWikiContent struct {
	Version   uint64
	CreatorId uint64
	CreatedAt int64
	Markdown  string
}

type RemoteWikiService interface {
	Load(ctx context.Context, wikiId uint64, wikiRef string, version uint64) (RawWikiContent, error)
	Store(ctx context.Context, wikiId uint64, userId uint64, wikiRef string, last uint64, markdown string) error
	GetVersions(ctx context.Context, wikiId uint64, wikiRef string) ([]RawWikiContent, error)
	Delete(ctx context.Context, wikiId uint64, wikiRef string, version uint64) error
}
