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
	"strings"

	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/spf13/afero"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/exp/slog"
)

type strengthConf struct {
	DefaultPassword string
	AllLang         []string
}

type initializedStrengthConf struct {
	minEntropy     float64
	localizedRules map[string]string
}

func initStrengthConf(logger *slog.Logger, conf *strengthConf) (initializedStrengthConf, error) {
	// TODO manage switch to network FS
	fileSystem := afero.NewOsFs()

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
		lang = strings.TrimSpace(lang)

		var pathBuilder strings.Builder
		pathBuilder.WriteString("rules/rules_")
		pathBuilder.WriteString(lang)
		pathBuilder.WriteString(".txt")
		content, err := afero.ReadFile(fileSystem, pathBuilder.String())
		if err == nil {
			localizedRules[lang] = strings.TrimSpace(string(content))
		} else {
			return nil, err
		}
	}
	return localizedRules, nil
}
