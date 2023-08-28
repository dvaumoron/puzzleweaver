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
	"bufio"
	"errors"
	"html/template"
	"io/fs"
	"strings"

	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/spf13/afero"
	"golang.org/x/exp/slog"
)

var errNoLocale = errors.New("no locales declared")

type templateConf struct {
	AllLang      []string
	TemplatePath string
	LocalesPath  string
}

type initializedTemplateConf struct {
	templates *template.Template
	messages  map[string]map[string]string
}

func initTemplateConf(logger *slog.Logger, conf *templateConf) *initializedTemplateConf {
	// TODO manage switch to network FS
	fileSystem := afero.NewOsFs()

	templates := loadTemplates(logger, fileSystem, conf)
	messages := loadLocales(logger, fileSystem, conf)
	return &initializedTemplateConf{templates: templates, messages: messages}
}

func loadTemplates(logger *slog.Logger, fileSystem afero.Fs, conf *templateConf) *template.Template {
	templatesPath := cleanPath(conf.TemplatePath)

	tmpl := template.New("")
	inSize := len(templatesPath)
	err := afero.Walk(fileSystem, templatesPath, func(path string, fi fs.FileInfo, err error) error {
		if err == nil && !fi.IsDir() {
			name := path[inSize:]
			if end := len(name) - 5; name[end:] == ".html" {
				var data []byte
				if data, err = afero.ReadFile(fileSystem, path); err == nil {
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

func loadLocales(logger *slog.Logger, fileSystem afero.Fs, conf *templateConf) map[string]map[string]string {
	if len(conf.AllLang) == 0 {
		logger.Error(servicecommon.NolocalesErrorMsg)
		return nil
	}

	localesPath := cleanPath(conf.LocalesPath)
	messages := make(map[string]map[string]string, len(conf.AllLang))
	for _, lang := range conf.AllLang {
		messagesLang := map[string]string{}
		messages[lang] = messagesLang

		var pathBuilder strings.Builder
		pathBuilder.WriteString(localesPath)
		pathBuilder.WriteString("/messages_")
		pathBuilder.WriteString(lang)
		pathBuilder.WriteString(".properties")

		if err := parseFile(fileSystem, pathBuilder.String(), messagesLang); err != nil {
			logger.Error(servicecommon.LoadFileErrorMsg, common.ErrorKey, err)
		}
	}

	defaultLang := conf.AllLang[0]
	messagesDefaultLang := messages[defaultLang]
	for _, lang := range conf.AllLang {
		if lang == defaultLang {
			continue
		}
		messagesLang := messages[lang]
		for key, value := range messagesDefaultLang {
			if messagesLang[key] == "" {
				messagesLang[key] = value
			}
		}
	}
	return messages
}

// separated function to close file sooner
func parseFile(fileSystem afero.Fs, path string, messagesLang map[string]string) error {
	file, err := fileSystem.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && line[0] != '#' {
			if equal := strings.Index(line, "="); equal > 0 {
				if key := strings.TrimSpace(line[:equal]); key != "" {
					if value := strings.TrimSpace(line[equal+1:]); value != "" {
						messagesLang[key] = value
					}
				}
			}
		}
	}
	return scanner.Err()
}

func cleanPath(path string) string {
	if last := len(path) - 1; last == -1 || path[last] != '/' {
		path += "/"
	}
	return path
}
