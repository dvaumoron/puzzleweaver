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
	"log/slog"

	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	gallerywidget "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/gallery"
	galleryimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/gallery/service/impl"
	widgethelper "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/helper"
)

type widgetConf struct {
	GalleryMongoAddress      string
	GalleryMongoDatabaseName string
	DefaultPageSize          uint64
}

type initializedWidgetConf struct {
	widgets widgethelper.WidgetManager
}

func initWidgetConf(loggerGetter servicecommon.LoggerGetter, logger *slog.Logger, conf *widgetConf) initializedWidgetConf {
	galleryService := galleryimpl.New(conf.GalleryMongoAddress, conf.GalleryMongoDatabaseName, loggerGetter)

	widgets := widgethelper.NewManager()
	gallerywidget.InitWidget(widgets, logger, galleryService, conf.DefaultPageSize)
	return initializedWidgetConf{widgets: widgets}
}
