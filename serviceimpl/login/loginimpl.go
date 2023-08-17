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
	"errors"
	"sort"
	"time"

	"github.com/ServiceWeaver/weaver"
	pb "github.com/dvaumoron/puzzleloginservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

// check matching with interface
var _ service.LoginService = &loginImpl{}

var errWeakPassword = errors.New("WeakPassword")

type loginImpl struct {
	weaver.Implements[service.FullLoginService]
	saltService     weaver.Ref[service.SaltService]
	strengthService weaver.Ref[service.PasswordStrengthService]
	dateFormat      string
}

type sortableContents []*pb.User

func (s sortableContents) Len() int {
	return len(s)
}

func (s sortableContents) Less(i, j int) bool {
	return s[i].Login < s[j].Login
}

func (s sortableContents) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (impl loginImpl) Verify(ctx context.Context, login string, password string) (bool, uint64, error) {
	_, err := impl.saltService.Get().Salt(ctx, login, password)
	if err != nil {
		return false, 0, common.LogOriginalError(impl.Logger(ctx), err)
	}

	success := true
	id := uint64(0)
	// TODO
	return success, id, nil
}

func (impl loginImpl) Register(ctx context.Context, login string, password string) (bool, uint64, error) {
	strong, err := impl.strengthService.Get().Validate(ctx, password)
	if err != nil {
		return false, 0, err
	}
	if !strong {
		return false, 0, errWeakPassword
	}

	_, err = impl.saltService.Get().Salt(ctx, login, password)
	if err != nil {
		return false, 0, common.LogOriginalError(impl.Logger(ctx), err)
	}

	success := true
	id := uint64(0)
	// TODO
	return success, id, nil
}

// You should remove duplicate id in list
func (impl loginImpl) GetUsers(ctx context.Context, userIds []uint64) (map[uint64]service.User, error) {
	list := []*pb.User{}
	// TODO
	logins := map[uint64]service.User{}
	for _, value := range list {
		logins[value.Id] = convertUser(value, impl.dateFormat)
	}
	return logins, nil
}

func (impl loginImpl) ChangeLogin(ctx context.Context, userId uint64, oldLogin string, newLogin string, password string) error {
	oldSalted, err := impl.saltService.Get().Salt(ctx, oldLogin, password)
	if err != nil {
		return common.LogOriginalError(impl.Logger(ctx), err)
	}

	newSalted, err := impl.saltService.Get().Salt(ctx, newLogin, password)
	if err != nil {
		return common.LogOriginalError(impl.Logger(ctx), err)
	}

	success := oldSalted == newSalted
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl loginImpl) ChangePassword(ctx context.Context, userId uint64, login string, oldPassword string, newPassword string) error {
	strong, err := impl.strengthService.Get().Validate(ctx, newPassword)
	if err != nil {
		return err
	}
	if !strong {
		return errWeakPassword
	}

	oldSalted, err := impl.saltService.Get().Salt(ctx, login, oldPassword)
	if err != nil {
		return common.LogOriginalError(impl.Logger(ctx), err)
	}

	newSalted, err := impl.saltService.Get().Salt(ctx, login, newPassword)
	if err != nil {
		return common.LogOriginalError(impl.Logger(ctx), err)
	}

	success := oldSalted == newSalted
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl loginImpl) ListUsers(ctx context.Context, start uint64, end uint64, filter string) (uint64, []service.User, error) {
	// TODO
	total := uint64(0)
	list := []*pb.User{}
	sort.Sort(sortableContents(list))
	users := make([]service.User, 0, len(list))
	for _, user := range list {
		users = append(users, convertUser(user, impl.dateFormat))
	}
	return total, users, nil
}

// no right check
func (impl loginImpl) Delete(ctx context.Context, userId uint64) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func convertUser(user *pb.User, dateFormat string) service.User {
	registredAt := time.Unix(user.RegistredAt, 0)
	return service.User{Id: user.Id, Login: user.Login, RegistredAt: registredAt.Format(dateFormat)}
}
