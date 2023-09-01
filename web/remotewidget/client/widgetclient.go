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

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	"github.com/gin-gonic/gin"
)

type widgetServiceWrapper struct {
	widgetService remoteservice.RemoteWidgetService
	loggerGetter  common.LoggerGetter
	widgetName    string
	objectId      uint64
	groupId       uint64
}

func MakeWidgetServiceWrapper(widgetService remoteservice.RemoteWidgetService, loggerGetter common.LoggerGetter, widgetName string, objectId uint64, groupId uint64) widgetservice.WidgetService {
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
	data[remoteservice.ObjectIdKey] = client.objectId
	data[remoteservice.GroupIdKey] = client.groupId
	dataBytes, err := json.Marshal(data)
	if err != nil {
		client.loggerGetter.Logger(ctx).Error("Failed to marshal data", common.ErrorKey, err)
		return "", "", nil, common.ErrTechnical
	}

	files[remoteservice.DataKey] = dataBytes
	return client.widgetService.Process(ctx, client.widgetName, actionName, files)
}

func convertActions(actions []remoteservice.RawWidgetAction) []widgetservice.WidgetAction {
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
	case remoteservice.KIND_HEAD:
		return http.MethodHead
	case remoteservice.KIND_POST:
		return http.MethodPost
	case remoteservice.KIND_PUT:
		return http.MethodPut
	case remoteservice.KIND_PATCH:
		return http.MethodPatch
	case remoteservice.KIND_DELETE:
		return http.MethodDelete
	case remoteservice.KIND_CONNECT:
		return http.MethodConnect
	case remoteservice.KIND_OPTIONS:
		return http.MethodOptions
	case remoteservice.KIND_TRACE:
		return http.MethodTrace
	case remoteservice.KIND_RAW:
		return widgetservice.RawResult
	}
	return http.MethodGet
}
