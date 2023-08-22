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

package loginclient

import (
	"context"
	"errors"
	"time"

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

var errNotEnoughValues = errors.New("not enough return values from saltService call")

type loginServiceWrapper struct {
	loginService    remoteservice.RemoteLoginService
	saltService     service.SaltService
	strengthService service.PasswordStrengthService
	dateFormat      string
}

func MakeLoginServiceWrapper(loginService remoteservice.RemoteLoginService, saltService service.SaltService, strengthService service.PasswordStrengthService, dateFormat string) service.LoginService {
	return loginServiceWrapper{
		loginService: loginService, saltService: saltService, strengthService: strengthService, dateFormat: dateFormat,
	}
}

func (client loginServiceWrapper) Verify(ctx context.Context, login string, password string) (uint64, error) {
	salteds, err := client.saltService.Salt(ctx, [2]string{login, password})
	if err != nil {
		return 0, err
	}
	if len(salteds) == 0 {
		return 0, errNotEnoughValues
	}
	return client.loginService.Verify(ctx, login, salteds[0])
}

func (client loginServiceWrapper) Register(ctx context.Context, login string, password string) (uint64, error) {
	err := client.strengthService.Validate(ctx, password)
	if err != nil {
		return 0, err
	}

	salteds, err := client.saltService.Salt(ctx, [2]string{login, password})
	if err != nil {
		return 0, err
	}
	if len(salteds) == 0 {
		return 0, errNotEnoughValues
	}
	return client.loginService.Register(ctx, login, salteds[0])
}

// You should remove duplicate id in list
func (client loginServiceWrapper) GetUsers(ctx context.Context, userIds []uint64) (map[uint64]service.User, error) {
	rawUsers, err := client.loginService.GetUsers(ctx, userIds)
	if err != nil {
		return nil, err
	}

	users := map[uint64]service.User{}
	for _, value := range rawUsers {
		users[value.Id] = convertUser(value, client.dateFormat)
	}
	return users, nil
}

func (client loginServiceWrapper) ChangeLogin(ctx context.Context, userId uint64, oldLogin string, newLogin string, password string) error {
	// avoid useless call
	if oldLogin == newLogin {
		return nil
	}

	salteds, err := client.saltService.Salt(ctx, [2]string{oldLogin, password}, [2]string{newLogin, password})
	if err != nil {
		return err
	}
	if len(salteds) < 2 {
		return errNotEnoughValues
	}
	return client.loginService.ChangeLogin(ctx, userId, newLogin, salteds[0], salteds[1])
}

func (client loginServiceWrapper) ChangePassword(ctx context.Context, userId uint64, login string, oldPassword string, newPassword string) error {
	// avoid useless call
	if oldPassword == newPassword {
		return nil
	}

	err := client.strengthService.Validate(ctx, newPassword)
	if err != nil {
		return err
	}

	salteds, err := client.saltService.Salt(ctx, [2]string{login, oldPassword}, [2]string{login, newPassword})
	if err != nil {
		return err
	}
	if len(salteds) < 2 {
		return errNotEnoughValues
	}

	// avoid useless call (unlikely since oldPassword != newPassword)
	if salteds[0] == salteds[1] {
		return nil
	}
	return client.loginService.ChangePassword(ctx, userId, salteds[0], salteds[1])
}

func (client loginServiceWrapper) ListUsers(ctx context.Context, start uint64, end uint64, filter string) (uint64, []service.User, error) {
	total, list, err := client.loginService.ListUsers(ctx, start, end, filter)
	if err != nil {
		return 0, nil, err
	}

	users := make([]service.User, 0, len(list))
	for _, user := range list {
		users = append(users, convertUser(user, client.dateFormat))
	}
	return total, users, nil
}

// no right check
func (client loginServiceWrapper) Delete(ctx context.Context, userId uint64) error {
	return client.loginService.Delete(ctx, userId)
}

func convertUser(user remoteservice.RawUser, dateFormat string) service.User {
	registredAt := time.Unix(user.RegistredAt, 0)
	return service.User{Id: user.Id, Login: user.Login, RegistredAt: registredAt.Format(dateFormat)}
}
