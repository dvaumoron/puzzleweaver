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

package templates

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

const contentTypeName = "Content-Type"

var htmlContentType = []string{"text/html; charset=utf-8"}

type ContextAndData struct {
	Ctx  context.Context
	Data gin.H
}

// match Render interface from gin.
type remoteHTML struct {
	remoteHTMLRender
	dataWithCtx  ContextAndData
	templateName string
}

func (r remoteHTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	ctx := r.dataWithCtx.Ctx
	data := r.dataWithCtx.Data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		r.loggerGetter.Logger(ctx).Error("Failed to marshal data", common.ErrorKey, err)
		return common.ErrTechnical
	}
	content, err := r.templateService.Render(ctx, r.templateName, dataBytes)
	if err != nil {
		return err
	}
	_, err = w.Write(content)
	return err
}

// Writes HTML ContentType.
func (r remoteHTML) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header[contentTypeName]; len(val) == 0 {
		header[contentTypeName] = htmlContentType
	}
}

// match HTMLRender interface from gin.
type remoteHTMLRender struct {
	templateService service.TemplateService
	loggerGetter    common.LoggerGetter
}

func (r remoteHTMLRender) Instance(name string, dataWithCtx any) render.Render {
	casted := dataWithCtx.(ContextAndData)
	return remoteHTML{remoteHTMLRender: r, dataWithCtx: casted, templateName: name}
}

func NewServiceRender(templateService service.TemplateService, loggerGetter common.LoggerGetter) render.HTMLRender {
	return remoteHTMLRender{templateService: templateService, loggerGetter: loggerGetter}
}
