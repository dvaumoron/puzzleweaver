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

package service

import (
	"context"
	"sync"

	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

type WikiContent struct {
	Version   uint64
	Markdown  string
	bodyMutex sync.RWMutex
	body      string
}

// Lazy loading for markdown application on body.
func (content *WikiContent) GetBody(ctx context.Context, markdownService service.MarkdownService) (string, error) {
	content.bodyMutex.RLock()
	body := content.body
	content.bodyMutex.RUnlock()
	if body != "" {
		return body, nil
	}
	markdown := content.Markdown
	if markdown == "" {
		return "", nil
	}

	content.bodyMutex.Lock()
	defer content.bodyMutex.Unlock()
	if body = content.body; body != "" {
		return body, nil
	}

	body, err := markdownService.Apply(ctx, markdown)
	if err != nil {
		return "", err
	}

	content.body = body
	return body, nil
}

type Version struct {
	Number  uint64
	Creator service.UserProfile
}

type WikiService interface {
	LoadContent(ctx context.Context, userId uint64, lang string, title string, versionStr string) (*WikiContent, error)
	StoreContent(ctx context.Context, userId uint64, lang string, title string, last string, markdown string) (bool, error)
	GetVersions(ctx context.Context, userId uint64, lang string, title string) ([]Version, error)
	DeleteContent(ctx context.Context, userId uint64, lang string, title string, versionStr string) error
	DeleteRight(ctx context.Context, userId uint64) bool
}
