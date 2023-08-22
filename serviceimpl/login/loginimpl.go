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

package loginimpl

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
)

// check matching with interface
var _ remoteservice.RemoteLoginService = &loginImpl{}

type loginImpl struct {
	weaver.Implements[remoteservice.RemoteLoginService]
}

func (impl *loginImpl) Verify(ctx context.Context, login string, salted string) (uint64, error) {
	id := uint64(0)
	// TODO
	return id, nil
}

func (impl *loginImpl) Register(ctx context.Context, login string, salted string) (uint64, error) {
	id := uint64(0)
	// TODO
	return id, nil
}

func (impl *loginImpl) GetUsers(ctx context.Context, userIds []uint64) (map[uint64]remoteservice.RawUser, error) {
	// TODO
	return nil, nil
}

func (impl *loginImpl) ChangeLogin(ctx context.Context, userId uint64, newLogin string, oldSalted string, newSalted string) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl *loginImpl) ChangePassword(ctx context.Context, userId uint64, oldSalted string, newSalted string) error {
	success := oldSalted == newSalted
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl *loginImpl) ListUsers(ctx context.Context, start uint64, end uint64, filter string) (uint64, []remoteservice.RawUser, error) {
	// TODO
	total := uint64(0)
	return total, nil, nil
}

func (impl *loginImpl) Delete(ctx context.Context, userId uint64) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}
