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

package remotewidgetservice

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

const (
	KIND_GET uint8 = iota
	KIND_HEAD
	KIND_POST
	KIND_PUT
	KIND_PATCH
	KIND_DELETE
	KIND_CONNECT
	KIND_OPTIONS
	KIND_TRACE
	KIND_RAW // added special category
)

type RawWidgetAction struct {
	weaver.AutoMarshal
	Kind       uint8
	Name       string
	Path       string
	QueryNames []string
}

type RemoteWidgetService interface {
	GetDesc(ctx context.Context, widgetName string) ([]RawWidgetAction, error)
	Process(ctx context.Context, widgetName string, actionName string, files map[string][]byte) (string, string, []byte, error)
}
