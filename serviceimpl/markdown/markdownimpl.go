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

package markdownimpl

import (
	"context"
	"strings"

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

type MarkdownService service.MarkdownService

type markdownImpl struct {
	weaver.Implements[MarkdownService]
	weaver.WithConfig[markdownConf]
	initializedConf initializedMarkdownConf
}

func (impl *markdownImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf = initMarkdownConf(impl.Config())
	return
}

func (impl *markdownImpl) Apply(ctx context.Context, text string) (string, error) {
	var resBuilder strings.Builder
	if err := impl.initializedConf.engine.Convert([]byte(text), &resBuilder); err != nil {
		impl.Logger(ctx).Error("Failed to transform markdown", common.ErrorKey, err)
		return "", servicecommon.ErrInternal
	}
	return resBuilder.String(), nil
}
