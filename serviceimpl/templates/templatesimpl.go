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
	"sync"

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"golang.org/x/exp/slog"
)

type TemplateService service.TemplateService

type templateImpl struct {
	weaver.Implements[TemplateService]
	weaver.WithConfig[templateConf]
	confMutex       sync.RWMutex
	initializedConf *initializedTemplateConf
}

func (impl *templateImpl) getInitializedConf(logger *slog.Logger) (*initializedTemplateConf, error) {
	impl.confMutex.RLock()
	initializedConf := impl.initializedConf
	impl.confMutex.RUnlock()
	if initializedConf != nil {
		return initializedConf, nil
	}

	impl.confMutex.Lock()
	defer impl.confMutex.Unlock()
	if impl.initializedConf == nil {
		var err error
		impl.initializedConf, err = initTemplateConf(logger, impl.Config())
		if err != nil {
			return nil, err
		}
	}
	return impl.initializedConf, nil
}

func (impl *templateImpl) Render(ctx context.Context, templateName string, data []byte) ([]byte, error) {
	logger := impl.Logger(ctx)
	initializedConf, err := impl.getInitializedConf(logger)
	if err != nil {
		logger.Error("Failed to load templates configuration", common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	var parsedData map[string]any
	if err = json.Unmarshal(data, &parsedData); err != nil {
		logger.Error("Failed to parse JSON", common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	parsedData["Messages"] = initializedConf.messages[asString(parsedData["lang"])]
	var content bytes.Buffer
	if err = initializedConf.templates.ExecuteTemplate(&content, templateName, parsedData); err != nil {
		logger.Error("Failed to call go template", common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return content.Bytes(), nil
}

func asString(value any) string {
	s, _ := value.(string)
	return s
}
