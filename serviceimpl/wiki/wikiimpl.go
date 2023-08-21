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
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
)

// check matching with interface
var _ remoteservice.RemoteWikiService = &remoteWikiImpl{}

type remoteWikiImpl struct {
	weaver.Implements[remoteservice.RemoteWikiService]
}

func (remoteWikiImpl) Load(ctx context.Context, wikiId uint64, wikiRef string, version uint64) (remoteservice.RawWikiContent, error) {
	// TODO
	return remoteservice.RawWikiContent{}, nil
}

func (remoteWikiImpl) Store(ctx context.Context, wikiId uint64, userId uint64, wikiRef string, last uint64, markdown string) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (remoteWikiImpl) GetVersions(ctx context.Context, wikiId uint64, wikiRef string) ([]remoteservice.RawWikiContent, error) {
	// TODO
	return nil, nil
}

func (remoteWikiImpl) Delete(ctx context.Context, wikiId uint64, wikiRef string, version uint64) error {
	// TODO
	return nil
}
