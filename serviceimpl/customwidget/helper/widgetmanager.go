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

package remotewidgethelper

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ActionHandler = func(context.Context, gin.H) (string, string, []byte, error)

type Action struct {
	Kind       uint8
	Path       string
	QueryNames []string
	Handler    ActionHandler
}

type Widget map[string]Action

// based on gin path convention, with the path "/view/:id/:name"
// the map passed to handler will contains "pathData/id" and "pathData/name" entries
// handler returned values are supposed to be redirect, templateName and data :
//
//  1. redirect is a redirect path (ignored if empty), to build an absolute one on the site the map contains the "CurrentUrl" entry
//
//  2. data could be :
//
//     - a json marshalled map which entries will be added to the data passed to the template engine with templateName
//
//     - or any raw data when the action kind is remoteservice.KIND_RAW
func (w Widget) AddAction(actionName string, kind uint8, path string, handler ActionHandler) {
	w[actionName] = Action{Kind: kind, Path: path, Handler: handler}
}

// Like AddAction but allow to indicate which query parameters should be transmitted.
func (w Widget) AddActionWithQuery(actionName string, kind uint8, path string, queryNames []string, handler ActionHandler) {
	w[actionName] = Action{Kind: kind, Path: path, QueryNames: queryNames, Handler: handler}
}

type WidgetManager map[string]Widget

func NewManager() WidgetManager {
	return WidgetManager{}
}

func (manager WidgetManager) CreateWidget(widgetName string) Widget {
	widget, ok := manager[widgetName]
	if !ok {
		widget = Widget{}
		manager[widgetName] = widget
	}
	return widget
}
