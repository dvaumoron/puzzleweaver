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

package passwordstrengthimpl

import (
	"log/slog"
	"strings"

	fsclient "github.com/dvaumoron/puzzleweaver/client/fs"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/spf13/afero"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type strengthConf struct {
	DefaultPassword string
	AllLang         []string
	FsConf          fsclient.FsConf
	RuleFilePath    string
}

type initializedStrengthConf struct {
	minEntropy     float64
	localizedRules map[string]string
}

func initStrengthConf(logger *slog.Logger, conf *strengthConf) (initializedStrengthConf, error) {
	fileSystem, err := fsclient.New(conf.FsConf)
	if err != nil {
		return initializedStrengthConf{}, err
	}

	localizedRules, err := readRulesConfig(logger, fileSystem, conf)
	if err != nil {
		return initializedStrengthConf{}, err
	}

	return initializedStrengthConf{
		minEntropy: passwordvalidator.GetEntropy(conf.DefaultPassword), localizedRules: localizedRules,
	}, nil
}

func readRulesConfig(logger *slog.Logger, fileSystem afero.Fs, conf *strengthConf) (map[string]string, error) {
	if len(conf.AllLang) == 0 {
		return nil, servicecommon.ErrNolocales
	}

	localizedRules := make(map[string]string, len(conf.AllLang))
	for _, lang := range conf.AllLang {
		path := strings.ReplaceAll(conf.RuleFilePath, servicecommon.LangPlaceHolder, lang)
		content, err := afero.ReadFile(fileSystem, path)
		if err == nil {
			localizedRules[lang] = strings.TrimSpace(string(content))
		} else {
			return nil, err
		}
	}
	return localizedRules, nil
}
