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
	"strings"
	"text/template"
	"time"

	"github.com/dvaumoron/partrenderer"
	fsclient "github.com/dvaumoron/puzzleweaver/client/fs"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/spf13/afero"
)

type templateConf struct {
	AllLang        []string
	FsConf         fsclient.FsConf
	ComponentsPath string
	ViewsPath      string
	LocaleFilePath string
	DateFormat     string
}

type initializedTemplateConf struct {
	templates partrenderer.PartRenderer
	messages  map[string]map[string]string
}

func initTemplateConf(conf *templateConf) (initializedTemplateConf, error) {
	fileSystem, err := fsclient.New(conf.FsConf)
	if err != nil {
		return initializedTemplateConf{}, err
	}

	templates, err := loadTemplates(fileSystem, conf)
	if err != nil {
		return initializedTemplateConf{}, err
	}

	messages, err := loadLocales(fileSystem, conf)
	if err != nil {
		return initializedTemplateConf{}, err
	}
	return initializedTemplateConf{templates: templates, messages: messages}, nil
}

func loadTemplates(fileSystem afero.Fs, conf *templateConf) (partrenderer.PartRenderer, error) {
	sourceFormat := conf.DateFormat
	customFuncs := template.FuncMap{"date": func(value string, targetFormat string) string {
		if sourceFormat == targetFormat {
			return value
		}
		date, err := time.Parse(sourceFormat, value)
		if err != nil {
			return value
		}
		return date.Format(targetFormat)
	}}
	return partrenderer.MakePartRenderer(conf.ComponentsPath, conf.ViewsPath, partrenderer.WithFs(fileSystem), partrenderer.WithFuncs(customFuncs))
}

// TODO merge with puzzlelocaleloader to call a shared library
func loadLocales(fileSystem afero.Fs, conf *templateConf) (map[string]map[string]string, error) {
	if len(conf.AllLang) == 0 {
		return nil, servicecommon.ErrNolocales
	}

	messages := make(map[string]map[string]string, len(conf.AllLang))
	for _, lang := range conf.AllLang {
		messagesLang := map[string]string{}
		messages[lang] = messagesLang

		path := strings.ReplaceAll(conf.LocaleFilePath, servicecommon.LangPlaceHolder, lang)
		if err := parseFile(fileSystem, path, messagesLang); err != nil {
			return nil, err
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
	return messages, nil
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
