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
	"sync"

	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
	"golang.org/x/exp/slog"
)

const wikiRefName = "wikiRef"

type wikiCache struct {
	mutex sync.RWMutex
	cache map[string]*wikiservice.WikiContent
}

func newCache() *wikiCache {
	return &wikiCache{cache: map[string]*wikiservice.WikiContent{}}
}

func (wiki *wikiCache) load(logger *slog.Logger, wikiRef string) *wikiservice.WikiContent {
	wiki.mutex.RLock()
	content, ok := wiki.cache[wikiRef]
	wiki.mutex.RUnlock()
	if !ok {
		logger.Debug("wikiCache miss", wikiRefName, wikiRef)
	}
	return content
}

func (wiki *wikiCache) store(logger *slog.Logger, wikiRef string, content *wikiservice.WikiContent) {
	wiki.mutex.Lock()
	wiki.cache[wikiRef] = content
	wiki.mutex.Unlock()
	logger.Debug("wikiCache store", wikiRefName, wikiRef)
}

func (wiki *wikiCache) delete(logger *slog.Logger, wikiRef string) {
	wiki.mutex.Lock()
	delete(wiki.cache, wikiRef)
	wiki.mutex.Unlock()
	logger.Debug("wikiCache delete", wikiRefName, wikiRef)
}
