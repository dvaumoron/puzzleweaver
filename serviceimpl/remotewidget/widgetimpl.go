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

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	widgethelper "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/helper"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/gin-gonic/gin"
)

const widgetNotFoundErrorMsg = "No widget found with requested name"
const widgetNameKey = "widgetName"

type RemoteWidgetService remoteservice.RemoteWidgetService

type remoteWidgetImpl struct {
	weaver.Implements[RemoteWidgetService]
	weaver.WithConfig[widgethelper.WidgetConf]
	initializedConf widgethelper.InitializedWidgetConf
}

func (impl *remoteWidgetImpl) Init(ctx context.Context) error {
	impl.initializedConf = widgethelper.InitWidgetConf(impl.Config())
	return nil
}

func (impl *remoteWidgetImpl) GetDesc(ctx context.Context, widgetName string) ([]remoteservice.RawWidgetAction, error) {
	widget, ok := impl.initializedConf.Widgets[widgetName]
	if !ok {
		impl.Logger(ctx).Error(widgetNotFoundErrorMsg, widgetNameKey, widgetName)
		return nil, servicecommon.ErrInternal
	}
	return convertActions(widget), nil
}

func (impl *remoteWidgetImpl) Process(ctx context.Context, widgetName string, actionName string, files map[string][]byte) (string, string, []byte, error) {
	widget, ok := impl.initializedConf.Widgets[widgetName]
	if !ok {
		impl.Logger(ctx).Error(widgetNotFoundErrorMsg, widgetNameKey, widgetName)
		return "", "", nil, servicecommon.ErrInternal
	}
	action, ok := widget[actionName]
	if !ok {
		impl.Logger(ctx).Error("No action found with requested names", widgetNameKey, widgetName, "actionName", actionName)
		return "", "", nil, servicecommon.ErrInternal
	}

	dataBytes := files[remoteservice.DataKey]

	var data gin.H
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		impl.Logger(ctx).Error("Failed to unmarshal data.json from call", common.ErrorKey, err)
		return "", "", nil, servicecommon.ErrInternal
	}
	// cleaning for GC
	dataBytes = nil
	delete(files, remoteservice.DataKey)

	if len(files) != 0 {
		data[widgethelper.FilesKey] = files
	}

	redirect, templateName, resData, err := action.Handler(ctx, data)
	if err != nil {
		impl.Logger(ctx).Error("Failed to handle action", common.ErrorKey, err)
		return "", "", nil, servicecommon.ErrInternal
	}
	return redirect, templateName, resData, nil

}

func convertActions(widget widgethelper.Widget) []remoteservice.RawWidgetAction {
	actions := make([]remoteservice.RawWidgetAction, 0, len(widget))
	for key, value := range widget {
		actions = append(actions, remoteservice.RawWidgetAction{
			Kind: value.Kind, Name: key, Path: value.Path, QueryNames: value.QueryNames,
		})
	}
	return actions
}
