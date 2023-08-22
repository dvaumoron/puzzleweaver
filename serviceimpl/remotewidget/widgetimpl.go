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

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	pb "github.com/dvaumoron/puzzlewidgetservice"
)

type RemoteWidgetService remoteservice.RemoteWidgetService

type remoteWidgetImpl struct {
	weaver.Implements[RemoteWidgetService]
}

func (impl *remoteWidgetImpl) GetDesc(ctx context.Context, name string) ([]remoteservice.RawWidgetAction, error) {
	var _ pb.MethodKind
	// TODO
	return nil, nil
}

func (impl *remoteWidgetImpl) Process(ctx context.Context, widgetName string, actionName string, files map[string][]byte) (string, string, []byte, error) {
	// TODO
	return "Redirect", "TemplateName", nil, nil
}
