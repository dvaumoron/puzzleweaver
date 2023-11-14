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

package customwidgetclient

import (
	"context"
	"encoding/json"
	"net/http"

	customwidgetservice "github.com/dvaumoron/puzzleweaver/serviceimpl/customwidget/service"
	"github.com/dvaumoron/puzzleweb/common"
	"github.com/dvaumoron/puzzleweb/common/log"
	widgetservice "github.com/dvaumoron/puzzleweb/remotewidget/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type widgetServiceWrapper struct {
	widgetService customwidgetservice.CustomWidgetService
	loggerGetter  log.LoggerGetter
	widgetName    string
	objectId      uint64
	groupId       uint64
}

func MakeWidgetServiceWrapper(widgetService customwidgetservice.CustomWidgetService, loggerGetter log.LoggerGetter, widgetName string, objectId uint64, groupId uint64) widgetservice.WidgetService {
	return widgetServiceWrapper{
		widgetService: widgetService, loggerGetter: loggerGetter, widgetName: widgetName, objectId: objectId, groupId: groupId,
	}
}

func (client widgetServiceWrapper) GetDesc(ctx context.Context) ([]widgetservice.Action, error) {
	actions, err := client.widgetService.GetDesc(ctx, client.widgetName)
	if err != nil {
		return nil, err
	}
	return convertActions(actions), nil
}

func (client widgetServiceWrapper) Process(ctx context.Context, actionName string, data gin.H, files map[string][]byte) (string, string, []byte, error) {
	data[widgetservice.ObjectIdKey] = client.objectId
	data[widgetservice.GroupIdKey] = client.groupId
	dataBytes, err := json.Marshal(data)
	if err != nil {
		client.loggerGetter.Logger(ctx).Error("Failed to marshal data", zap.Error(err))
		return "", "", nil, common.ErrTechnical
	}

	files[widgetservice.DataKey] = dataBytes
	return client.widgetService.Process(ctx, client.widgetName, actionName, files)
}

func convertActions(actions []customwidgetservice.RawWidgetAction) []widgetservice.Action {
	res := make([]widgetservice.Action, 0, len(actions))
	for _, action := range actions {
		res = append(res, widgetservice.Action{
			Kind: converKind(action.Kind), Name: action.Name, Path: action.Path, QueryNames: action.QueryNames},
		)
	}
	return res
}

func converKind(kind uint8) string {
	switch kind {
	case customwidgetservice.KIND_HEAD:
		return http.MethodHead
	case customwidgetservice.KIND_POST:
		return http.MethodPost
	case customwidgetservice.KIND_PUT:
		return http.MethodPut
	case customwidgetservice.KIND_PATCH:
		return http.MethodPatch
	case customwidgetservice.KIND_DELETE:
		return http.MethodDelete
	case customwidgetservice.KIND_CONNECT:
		return http.MethodConnect
	case customwidgetservice.KIND_OPTIONS:
		return http.MethodOptions
	case customwidgetservice.KIND_TRACE:
		return http.MethodTrace
	case customwidgetservice.KIND_RAW:
		return widgetservice.RawResult
	}
	return http.MethodGet
}
