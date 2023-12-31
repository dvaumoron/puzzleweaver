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

package wikiimpl

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type RawWikiContent struct {
	weaver.AutoMarshal
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
