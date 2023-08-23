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

package templatesimpl

import (
	"html/template"
	"io/fs"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/spf13/afero"
	"golang.org/x/exp/slog"
)

type templateConf struct {
	TemplatePath string
}

type initializedTemplateConf struct {
	templates *template.Template
	messages  map[string]map[string]string
}

func initTemplateConf(logger *slog.Logger, conf *templateConf) *initializedTemplateConf {
	var fileSystem afero.Fs
	// TODO
	return &initializedTemplateConf{templates: load(logger, fileSystem, conf)}
}

func load(logger *slog.Logger, fileSystem afero.Fs, conf *templateConf) *template.Template {
	templatesPath := conf.TemplatePath
	if last := len(templatesPath) - 1; templatesPath[last] != '/' {
		templatesPath += "/"
	}

	tmpl := template.New("")
	inSize := len(templatesPath)
	err := afero.Walk(fileSystem, templatesPath, func(path string, d fs.FileInfo, err error) error {
		if err == nil && !d.IsDir() {
			name := path[inSize:]
			if end := len(name) - 5; name[end:] == ".html" {
				var data []byte
				data, err = afero.ReadFile(fileSystem, path)
				if err == nil {
					_, err = tmpl.New(name[:end]).Parse(string(data))
				}
			}
		}
		return err
	})

	if err != nil {
		logger.Error("Failed to load templates", common.ErrorKey, err)
	}
	return tmpl
}
