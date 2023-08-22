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

package templatesimpl

import (
	"context"
	"encoding/json"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

// check matching with interface
var _ service.TemplateService = &templateClient{}

type templateClient struct {
	weaver.Implements[service.TemplateService]
}

func (impl *templateClient) Render(ctx context.Context, templateName string, data any) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		impl.Logger(ctx).Error("Failed to marshal data", common.ErrorKey, err)
		return nil, common.ErrTechnical
	}

	// TODO
	return dataBytes, nil
}
