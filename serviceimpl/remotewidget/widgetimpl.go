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

package remotewidgetimpl

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web/common"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	pb "github.com/dvaumoron/puzzlewidgetservice"
	"github.com/gin-gonic/gin"
)

// check matching with interface
var _ widgetservice.WidgetService = &widgetImpl{}

type widgetImpl struct {
	weaver.Implements[widgetservice.WidgetService]
	objectId uint64
	groupId  uint64
}

func (impl widgetImpl) GetDesc(ctx context.Context, name string) ([]widgetservice.Action, error) {
	// TODO
	return convertActions(nil), nil
}

func (impl widgetImpl) Process(ctx context.Context, widgetName string, actionName string, data gin.H, files map[string][]byte) (string, string, []byte, error) {
	data["objectId"] = impl.objectId
	data["groupId"] = impl.groupId
	dataBytes, err := json.Marshal(data)
	if err != nil {
		impl.Logger(ctx).Error("Failed to marshal data", common.ErrorKey, err)
		return "", "", nil, common.ErrTechnical
	}

	files["puzzledata.json"] = dataBytes
	// TODO
	return "Redirect", "TemplateName", nil, nil
}

func convertActions(actions []*pb.Action) []widgetservice.Action {
	res := make([]widgetservice.Action, 0, len(actions))
	for _, action := range actions {
		res = append(res, widgetservice.Action{
			Kind: converKind(action.Kind), Name: action.Name, Path: action.Path, QueryNames: action.QueryNames},
		)
	}
	return res
}

func converKind(kind pb.MethodKind) string {
	switch kind {
	case pb.MethodKind_HEAD:
		return http.MethodHead
	case pb.MethodKind_POST:
		return http.MethodPost
	case pb.MethodKind_PUT:
		return http.MethodPut
	case pb.MethodKind_PATCH:
		return http.MethodPatch
	case pb.MethodKind_DELETE:
		return http.MethodDelete
	case pb.MethodKind_CONNECT:
		return http.MethodConnect
	case pb.MethodKind_OPTIONS:
		return http.MethodOptions
	case pb.MethodKind_TRACE:
		return http.MethodTrace
	case pb.MethodKind_RAW:
		return widgetservice.RawResult
	}
	return http.MethodGet
}
