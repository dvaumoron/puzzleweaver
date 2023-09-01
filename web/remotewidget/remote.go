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

package remotewidget

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web"
	"github.com/dvaumoron/puzzleweaver/web/common"
	widgetservice "github.com/dvaumoron/puzzleweaver/web/remotewidget/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

const initMsg = "Failed to init remote widget"

type handlerDesc struct {
	httpMethod string
	path       string
	handler    gin.HandlerFunc
}

type remoteWidget struct {
	handlers []handlerDesc
}

func (w remoteWidget) LoadInto(router gin.IRouter) {
	for _, desc := range w.handlers {
		router.Handle(desc.httpMethod, desc.path, desc.handler)
	}
}

func MakeRemotePage(pageName string, ctx context.Context, logger *slog.Logger, widgetService widgetservice.WidgetService) (web.Page, bool) {
	actions, err := widgetService.GetDesc(ctx)
	if err != nil {
		logger.Error(initMsg, common.ErrorKey, err)
		return web.Page{}, false
	}

	handlers := make([]handlerDesc, 0, len(actions))
	for _, action := range actions {
		httpMethod := action.Kind
		actionName := action.Name
		actionPath := action.Path
		pathKeys := extractKeysFromPath(actionPath)
		queryKeys := extractQueryKeys(action.QueryNames)
		var handler gin.HandlerFunc
		switch httpMethod {
		case http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace:
			dataAdder := func(data gin.H, c *gin.Context) {
				retrieveContextData(pathKeys, queryKeys, data, c)
			}
			handler = createHandler(actionName, dataAdder, widgetService)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			dataAdder := func(data gin.H, c *gin.Context) {
				data[remoteservice.FormKey] = c.PostFormMap(remoteservice.FormKey)
				retrieveContextData(pathKeys, queryKeys, data, c)
			}
			handler = createHandler(actionName, dataAdder, widgetService)
		case widgetservice.RawResult:
			httpMethod = http.MethodGet
			handler = func(c *gin.Context) {
				data := gin.H{}
				retrieveContextData(pathKeys, queryKeys, data, c)
				_, _, resData, err := widgetService.Process(ctx, actionName, data, map[string][]byte{})
				if err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Data(http.StatusOK, http.DetectContentType(resData), resData)
			}
		default:
			logger.Error(initMsg, "unknownActionKind", httpMethod)
			return web.Page{}, false
		}
		handlers = append(handlers, handlerDesc{httpMethod: httpMethod, path: actionPath, handler: handler})
	}

	p := web.MakePage(pageName)
	p.Widget = remoteWidget{handlers: handlers}
	return p, true
}

func extractKeysFromPath(path string) [][2]string {
	splitted := strings.Split(path, "/")
	keys := make([][2]string, 0, len(splitted))
	for _, part := range splitted {
		if len(part) != 0 && part[0] == ':' {
			key := part[1:]
			keys = append(keys, [2]string{remoteservice.PathKeySlash + key, key})
		}
	}
	return keys
}

func extractQueryKeys(names []string) [][2]string {
	keys := make([][2]string, 0, len(names))
	for _, name := range names {
		key := strings.TrimSpace(name)
		if len(key) != 0 {
			keys = append(keys, [2]string{remoteservice.QueryKeySlash + key, key})
		}
	}
	return keys
}

func retrieveContextData(pathKeys [][2]string, queryKeys [][2]string, data gin.H, c *gin.Context) {
	for _, key := range pathKeys {
		data[key[0]] = c.Param(key[1])
	}
	for _, key := range queryKeys {
		data[key[0]] = c.Query(key[1])
	}
	data[common.SessionName] = web.GetSession(c).AsMap()
}

func readFiles(c *gin.Context) (map[string][]byte, error) {
	files := map[string][]byte{}
	fileList := c.PostForm("fileList")
	if len(fileList) == 0 {
		return files, nil
	}

	for _, name := range strings.Split(fileList, ",") {
		if trimmed := strings.TrimSpace(name); len(trimmed) != 0 {
			if err := readFile(trimmed, files, c); err != nil {
				return nil, err
			}
		}
	}
	return files, nil
}

func readFile(name string, files map[string][]byte, c *gin.Context) error {
	header, err := c.FormFile(name)
	if err != nil {
		return nil // ignore non existing file here (widget should handle)
	}

	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil || len(fileData) == 0 {
		return err
	}
	files[name] = fileData
	return nil
}

func createHandler(actionName string, dataAdder common.DataAdder, widgetService widgetservice.WidgetService) gin.HandlerFunc {
	return web.CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
		ctx := c.Request.Context()
		logger := web.GetLogger(c)
		dataAdder(data, c)
		files, err := readFiles(c)
		if err != nil {
			logger.Error("Failed to retrieve post file", common.ErrorKey, err)
			return "", common.DefaultErrorRedirect(logger, common.ErrorTechnicalKey)
		}
		redirect, templateName, resData, err := widgetService.Process(ctx, actionName, data, files)
		if err != nil {
			return "", common.DefaultErrorRedirect(logger, err.Error())
		}
		if redirect != "" {
			return "", redirect
		}

		if err = updateDataAndSession(data, resData, c); err != nil {
			logger.Error("Failed to unmarshal json from remote widget", common.ErrorKey, err)
			return "", common.DefaultErrorRedirect(logger, common.ErrorTechnicalKey)
		}
		return templateName, ""
	})
}

func updateDataAndSession(data gin.H, resData []byte, c *gin.Context) error {
	var newData gin.H
	if err := json.Unmarshal(resData, &newData); err != nil {
		return err
	}
	for key, value := range newData {
		data[key] = value
	}
	sessionMap, sessionUpdate := newData[common.SessionName]
	if sessionUpdate {
		casted, ok := sessionMap.(map[string]string)
		if ok {
			session := web.GetSession(c)
			for key, value := range casted {
				session.Store(key, value)
			}
		}
	}
	return nil
}
