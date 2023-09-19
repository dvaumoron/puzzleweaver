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
	"bytes"
	"context"
	"encoding/json"

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
)

type templateImpl struct {
	weaver.Implements[TemplateService]
	weaver.WithConfig[templateConf]
	initializedConf initializedTemplateConf
}

func (impl *templateImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf, err = initTemplateConf(impl.Config())
	return
}

func (impl *templateImpl) Render(ctx context.Context, templateName string, data []byte) ([]byte, error) {
	logger := impl.Logger(ctx)

	var err error
	var parsedData map[string]any
	if err = json.Unmarshal(data, &parsedData); err != nil {
		logger.Error("Failed to parse JSON", common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	parsedData["Messages"] = impl.initializedConf.messages[asString(parsedData["lang"])]
	var content bytes.Buffer
	if err = impl.initializedConf.templates.ExecuteTemplate(&content, templateName, parsedData); err != nil {
		logger.Error("Failed to call go template", common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return content.Bytes(), nil
}

func asString(value any) string {
	s, _ := value.(string)
	return s
}
