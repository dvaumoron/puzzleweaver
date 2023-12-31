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

package fsclient

import (
	"errors"

	"github.com/spf13/afero"
)

var errUnknonwKind = errors.New("unknown file system kind")

type FsConf struct {
	Kind   string
	Params map[string]string
}

func New(conf FsConf) (afero.Fs, error) {
	// TODO add other kind
	switch conf.Kind {
	case "local":
		return afero.NewOsFs(), nil
	default:
		return nil, errUnknonwKind
	}
}
