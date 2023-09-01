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

package gallerywidget

import (
	"context"
	"encoding/json"

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	galleryservice "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/gallery/service"
	widgethelper "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/helper"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

const (
	imageKey       = "Image"
	imageIdKey     = "ImageId"
	pathImageIdKey = remoteservice.PathKeySlash + imageIdKey
)

func InitWidget(manager widgethelper.WidgetManager, logger *slog.Logger, service galleryservice.GalleryService, defaultPageSize uint64, args ...string) {
	viewTmpl := "gallery/view"
	editTmpl := "gallery/edit"
	switch len(args) {
	default:
		logger.Info("InitWidget should be called with 0 to 2 optional arguments.")
		fallthrough
	case 2:
		if args[1] != "" {
			editTmpl = args[1]
		}
		fallthrough
	case 1:
		if args[0] != "" {
			viewTmpl = args[0]
		}
	case 0:
	}

	w := manager.CreateWidget("gallery")
	w.AddActionWithQuery("list", remoteservice.KIND_GET, "/", widgethelper.GetPaginationNames(), func(ctx context.Context, data gin.H) (string, string, []byte, error) {
		pageNumber, start, end, _ := widgethelper.GetPagination(defaultPageSize, data)

		galleryId, err := widgethelper.AsUint64(data[remoteservice.ObjectIdKey])
		if err != nil {
			return "", "", nil, err
		}

		total, images, err := service.GetImages(ctx, galleryId, start, end)
		if err != nil {
			return "", "", nil, err
		}

		newData := gin.H{}
		common.InitPagination(newData, "", pageNumber, end, total)
		newData["Images"] = images
		resData, err := json.Marshal(newData)
		if err != nil {
			return "", "", nil, err
		}
		return "", viewTmpl, resData, nil
	})
	w.AddAction("retrieve", remoteservice.KIND_RAW, "/retrieve/:ImageId", func(ctx context.Context, data gin.H) (string, string, []byte, error) {
		imageId, err := widgethelper.AsUint64(data[pathImageIdKey])
		if err != nil {
			return "", "", nil, err
		}

		image, err := service.GetImageData(ctx, imageId)
		if err != nil {
			return "", "", nil, err
		}
		return "", "", image, nil
	})
	w.AddAction("create", remoteservice.KIND_GET, "/create", func(ctx context.Context, data gin.H) (string, string, []byte, error) {
		baseUrl, err := widgethelper.GetBaseUrl(1, data)
		if err != nil {
			return "", "", nil, err
		}

		newData := gin.H{}
		newData[imageKey] = galleryservice.GalleryImage{Title: "new"}
		newData[common.BaseUrlName] = baseUrl
		resData, err := json.Marshal(newData)
		if err != nil {
			return "", "", nil, err
		}
		return "", editTmpl, resData, nil
	})
	w.AddAction("edit", remoteservice.KIND_GET, "/edit/:ImageId", func(ctx context.Context, data gin.H) (string, string, []byte, error) {
		imageId, err := widgethelper.AsUint64(data[pathImageIdKey])
		if err != nil {
			return "", "", nil, err
		}

		image, err := service.GetImage(ctx, imageId)
		if err != nil {
			return "", "", nil, err
		}

		baseUrl, err := widgethelper.GetBaseUrl(2, data)
		if err != nil {
			return "", "", nil, err
		}

		newData := gin.H{}
		newData[imageKey] = image
		newData[common.BaseUrlName] = baseUrl
		resData, err := json.Marshal(newData)
		if err != nil {
			return "", "", nil, err
		}
		return "", editTmpl, resData, nil
	})
	w.AddAction("save", remoteservice.KIND_POST, "/save", func(ctx context.Context, data gin.H) (string, string, []byte, error) {
		galleryId, err := widgethelper.AsUint64(data[remoteservice.ObjectIdKey])
		if err != nil {
			return "", "", nil, err
		}

		listUrl, err := widgethelper.GetBaseUrl(1, data)
		if err != nil {
			return "", "", nil, err
		}

		userId, err := widgethelper.AsUint64(data[common.UserIdName])
		if err != nil {
			return "", "", nil, err
		}
		if userId == 0 {
			return listUrl + "?error=ErrorNotAuthorized", "", nil, nil
		}

		formData, err := widgethelper.GetFormData(data)
		if err != nil {
			return "", "", nil, err
		}

		imageId, err := widgethelper.AsUint64(formData[imageIdKey])
		if err != nil {
			return "", "", nil, err
		}

		title, err := widgethelper.AsString(formData["Title"])
		if err != nil {
			return "", "", nil, err
		}

		if title == "new" || title == "" {
			return listUrl + "?error=ErrorBadImageTitle", "", nil, nil
		}

		desc, err := widgethelper.AsString(formData["Desc"])
		if err != nil {
			return "", "", nil, err
		}

		imageInfo := galleryservice.GalleryImage{ImageId: imageId, CreatorId: userId, Title: title, Desc: desc}

		files, err := widgethelper.GetFiles(data)
		if err != nil {
			return "", "", nil, err
		}

		if _, err = service.UpdateImage(ctx, galleryId, imageInfo, files["image"]); err != nil {
			return "", "", nil, err
		}
		return listUrl, "", nil, nil
	})
	w.AddAction("delete", remoteservice.KIND_POST, "/delete/:ImageId", func(ctx context.Context, data gin.H) (string, string, []byte, error) {
		imageId, err := widgethelper.AsUint64(data[pathImageIdKey])
		if err != nil {
			return "", "", nil, err
		}

		if err = service.DeleteImage(ctx, imageId); err != nil {
			return "", "", nil, err
		}

		listUrl, err := widgethelper.GetBaseUrl(2, data)
		if err != nil {
			return "", "", nil, err
		}
		return listUrl, "", nil, nil
	})
}
