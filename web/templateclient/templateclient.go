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

package templateclient

import (
	"context"
	"encoding/json"

	templatesimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/templates"
	"github.com/dvaumoron/puzzleweb/common"
	"github.com/dvaumoron/puzzleweb/common/log"
	templateservice "github.com/dvaumoron/puzzleweb/templates/service"
	"go.uber.org/zap"
)

type templateServiceWrapper struct {
	templateService templatesimpl.TemplateService
	loggerGetter    log.LoggerGetter
}

func MakeTemplateServiceWrapper(templateService templatesimpl.TemplateService, loggerGetter log.LoggerGetter) templateservice.TemplateService {
	return templateServiceWrapper{templateService: templateService, loggerGetter: loggerGetter}
}

func (client templateServiceWrapper) Render(ctx context.Context, templateName string, data any) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		client.loggerGetter.Logger(ctx).Error("Failed to marshal data", zap.Error(err))
		return nil, common.ErrTechnical
	}
	return client.templateService.Render(ctx, templateName, dataBytes)
}
