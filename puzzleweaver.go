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

package main

import (
	"context"
	"log"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/frame"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	widgethelper "github.com/dvaumoron/puzzleweaver/serviceimpl/customwidget/helper"
)

// can be overridden with ldflags
var version = "dev"

func main() {
	ctx := context.Background()
	widgethelper.Registerers = append(widgethelper.Registerers, func(wm widgethelper.WidgetManager, conf map[string]string, lg servicecommon.LoggerGetter) error {
		lg.Logger(ctx).Info("Custom widget initialization", "conf", conf)
		return nil
	})

	if err := weaver.Run(ctx, frame.NewFrameServe(version)); err != nil {
		log.Fatal(err)
	}
}
