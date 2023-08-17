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

package sessionimpl

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

// check matching with interface
var _ service.SessionService = &sessionImpl{}

type sessionImpl struct {
	weaver.Implements[service.SessionService]
}

func (impl sessionImpl) Generate(ctx context.Context) (uint64, error) {
	// TODO
	return 0, nil
}

func (impl sessionImpl) Get(ctx context.Context, id uint64) (map[string]string, error) {
	// TODO
	return nil, nil
}

func (impl sessionImpl) Update(ctx context.Context, id uint64, info map[string]string) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}
