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

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleloginserver/model"
	dbclient "github.com/dvaumoron/puzzleweaver/client/db"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweb/common"
	"gorm.io/gorm"
)

type loginImpl struct {
	weaver.Implements[RemoteLoginService]
	weaver.WithConfig[loginConf]
	initializedConf initializedLoginConf
}

func (impl *loginImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf, err = initLoginConf(impl.Config())
	return
}

func (impl *loginImpl) Verify(ctx context.Context, login string, salted string) (uint64, error) {
	var user model.User
	if err := impl.initializedConf.db.First(&user, "login = ?", login).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, common.ErrWrongLogin
		}

		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, servicecommon.ErrInternal
	}

	if salted != user.Password {
		return 0, common.ErrWrongLogin
	}
	return user.ID, nil
}

func (impl *loginImpl) Register(ctx context.Context, login string, salted string) (uint64, error) {
	if login == "" {
		return 0, common.ErrEmptyLogin
	}

	var user model.User
	err := impl.initializedConf.db.First(&user, "login = ?", login).Error
	if err == nil {
		return 0, common.ErrExistingLogin
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// some technical error, send it
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, servicecommon.ErrInternal
	}

	// unknown user, create new
	user = model.User{Login: login, Password: salted}
	if err = impl.initializedConf.db.Create(&user).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, servicecommon.ErrInternal
	}
	return user.ID, nil
}

func (impl *loginImpl) GetUsers(ctx context.Context, userIds []uint64) (map[uint64]RawUser, error) {
	var users []model.User
	if err := impl.initializedConf.db.Find(&users, "id IN ?", userIds).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return convertUsersMapFromModel(users), nil
}

func (impl *loginImpl) ChangeLogin(ctx context.Context, userId uint64, newLogin string, oldSalted string, newSalted string) error {
	if newLogin == "" {
		return common.ErrEmptyLogin
	}

	var user model.User
	err := impl.initializedConf.db.First(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrWrongLogin
		}

		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}

	if oldSalted != user.Password {
		return common.ErrWrongLogin
	}

	err = impl.initializedConf.db.First(&user, "login = ?", newLogin).Error
	if err == nil {
		return common.ErrExistingLogin
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}

	err = impl.initializedConf.db.Model(&user).Updates(map[string]any{
		"login": newLogin, "password": newSalted,
	}).Error
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *loginImpl) ChangePassword(ctx context.Context, userId uint64, oldSalted string, newSalted string) error {
	var user model.User
	err := impl.initializedConf.db.First(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrWrongLogin
		}

		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}

	if oldSalted != user.Password {
		return common.ErrWrongLogin
	}
	if err = impl.initializedConf.db.Model(&user).Update("password", newSalted).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *loginImpl) ListUsers(ctx context.Context, start uint64, end uint64, filter string) (uint64, []RawUser, error) {
	noFilter := filter == ""

	userRequest := impl.initializedConf.db.Model(&model.User{})
	if !noFilter {
		filter = dbclient.BuildLikeFilter(filter)
		userRequest.Where("login LIKE ?", filter)
	}
	var total int64
	err := userRequest.Count(&total).Error
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}
	if total == 0 {
		return 0, nil, nil
	}

	var users []model.User
	page := dbclient.Paginate(impl.initializedConf.db, start, end).Order("login asc")
	if noFilter {
		err = page.Find(&users).Error
	} else {
		err = page.Find(&users, "login LIKE ?", filter).Error
	}

	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}
	return uint64(total), convertUsersFromModel(users), nil
}

func (impl *loginImpl) Delete(ctx context.Context, userId uint64) error {
	if err := impl.initializedConf.db.Delete(&model.User{}, userId).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func convertUsersFromModel(users []model.User) []RawUser {
	resUsers := make([]RawUser, 0, len(users))
	for _, user := range users {
		resUsers = append(resUsers, RawUser{Id: user.ID, Login: user.Login, RegistredAt: user.CreatedAt.Unix()})
	}
	return resUsers
}

func convertUsersMapFromModel(users []model.User) map[uint64]RawUser {
	resUsers := make(map[uint64]RawUser, len(users))
	for _, user := range users {
		resUsers[user.ID] = RawUser{Id: user.ID, Login: user.Login, RegistredAt: user.CreatedAt.Unix()}
	}
	return resUsers
}
