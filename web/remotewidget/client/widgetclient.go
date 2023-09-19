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

package remotewidgetclient

import (
	"context"
	"encoding/json"
	"net/http"

	remotewidgetservice "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/service"
	"github.com/dvaumoron/puzzleweaver/web/common"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	"github.com/gin-gonic/gin"
)

type widgetServiceWrapper struct {
	widgetService remotewidgetservice.RemoteWidgetService
	loggerGetter  common.LoggerGetter
	widgetName    string
	objectId      uint64
	groupId       uint64
}

func MakeWidgetServiceWrapper(widgetService remotewidgetservice.RemoteWidgetService, loggerGetter common.LoggerGetter, widgetName string, objectId uint64, groupId uint64) widgetservice.WidgetService {
	return widgetServiceWrapper{
		widgetService: widgetService, loggerGetter: loggerGetter, widgetName: widgetName, objectId: objectId, groupId: groupId,
	}
}

func (client widgetServiceWrapper) GetDesc(ctx context.Context) ([]widgetservice.WidgetAction, error) {
	actions, err := client.widgetService.GetDesc(ctx, client.widgetName)
	if err != nil {
		return nil, err
	}
	return convertActions(actions), nil
}

func (client widgetServiceWrapper) Process(ctx context.Context, actionName string, data gin.H, files map[string][]byte) (string, string, []byte, error) {
	data[remotewidgetservice.ObjectIdKey] = client.objectId
	data[remotewidgetservice.GroupIdKey] = client.groupId
	dataBytes, err := json.Marshal(data)
	if err != nil {
		client.loggerGetter.Logger(ctx).Error("Failed to marshal data", common.ErrorKey, err)
		return "", "", nil, common.ErrTechnical
	}

	files[remotewidgetservice.DataKey] = dataBytes
	return client.widgetService.Process(ctx, client.widgetName, actionName, files)
}

func convertActions(actions []remotewidgetservice.RawWidgetAction) []widgetservice.WidgetAction {
	res := make([]widgetservice.WidgetAction, 0, len(actions))
	for _, action := range actions {
		res = append(res, widgetservice.WidgetAction{
			Kind: converKind(action.Kind), Name: action.Name, Path: action.Path, QueryNames: action.QueryNames},
		)
	}
	return res
}

func converKind(kind uint8) string {
	switch kind {
	case remotewidgetservice.KIND_HEAD:
		return http.MethodHead
	case remotewidgetservice.KIND_POST:
		return http.MethodPost
	case remotewidgetservice.KIND_PUT:
		return http.MethodPut
	case remotewidgetservice.KIND_PATCH:
		return http.MethodPatch
	case remotewidgetservice.KIND_DELETE:
		return http.MethodDelete
	case remotewidgetservice.KIND_CONNECT:
		return http.MethodConnect
	case remotewidgetservice.KIND_OPTIONS:
		return http.MethodOptions
	case remotewidgetservice.KIND_TRACE:
		return http.MethodTrace
	case remotewidgetservice.KIND_RAW:
		return widgetservice.RawResult
	}
	return http.MethodGet
}
