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
)

type RawUser struct {
	weaver.AutoMarshal
	Id          uint64
	Login       string
	RegistredAt int64
}

type RemoteLoginService interface {
	GetUsers(ctx context.Context, userIds []uint64) (map[uint64]RawUser, error)
	ListUsers(ctx context.Context, start uint64, end uint64, filter string) (uint64, []RawUser, error)
	Delete(ctx context.Context, userId uint64) error
	Verify(ctx context.Context, login string, salted string) (uint64, error)
	Register(ctx context.Context, login string, salted string) (uint64, error)
	ChangeLogin(ctx context.Context, userId uint64, newLogin string, oldSalted string, newSalted string) error
	ChangePassword(ctx context.Context, userId uint64, oldSalted string, newSalted string) error
}
