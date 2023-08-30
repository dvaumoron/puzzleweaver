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

package adminimpl

import (
	"context"
	"errors"
	"sync"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/serviceimpl/admin/model"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/open-policy-agent/opa/rego"
	"gorm.io/gorm"
)

const (
	accessFlag = 1 << iota
	createFlag
	updateFlag
	deleteFlag
)

type AdminService service.AdminService

type adminImpl struct {
	weaver.Implements[AdminService]
	weaver.WithConfig[adminConf]
	initializedConf initializedAdminConf
	idToNameMutex   sync.RWMutex
	idToName        map[uint64]string
}

func (impl *adminImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf, err = initAdminConf(ctx, impl.Config())
	impl.idToName = map[uint64]string{}
	return
}

func (impl *adminImpl) GetAllGroups(ctx context.Context) ([]service.Group, error) {
	return impl.initializedConf.groups, nil
}

func (impl *adminImpl) AuthQuery(ctx context.Context, userId uint64, groupId uint64, action string) error {
	db := impl.initializedConf.db.WithContext(ctx)
	return impl.innerAuthQuery(ctx, db, userId, groupId, convertActionToFlag(action))
}

func (impl *adminImpl) GetActions(ctx context.Context, adminId uint64, roleName string, groupName string) ([]string, error) {
	db := impl.initializedConf.db.WithContext(ctx)
	err := impl.innerAuthQuery(ctx, db, adminId, service.AdminGroupId, accessFlag)
	if err != nil {
		return nil, err
	}

	var role model.Role
	groupId := impl.initializedConf.nameToGroupId[groupName]
	subQuery := db.Model(&model.RoleName{}).Select("id").Where("name = ?", roleName)
	if err = db.First(&role, "name_id IN (?) AND object_id = ?", subQuery, groupId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ignore unknown role
			return nil, nil
		}

		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return convertActionsFromFlags(role.ActionFlags), nil
}

func (impl *adminImpl) UpdateUser(ctx context.Context, adminId uint64, userId uint64, roles []service.Role) error {
	db := impl.initializedConf.db.WithContext(ctx)
	err := impl.innerAuthQuery(ctx, db, adminId, service.AdminGroupId, updateFlag)
	if err != nil {
		return err
	}

	mRoles, err := impl.loadRoles(ctx, db, roles)
	if err != nil {
		return err
	}

	rolesLen := len(mRoles)
	if rolesLen == 0 {
		// delete unused user
		return impl.handleUpdateError(ctx, db.Delete(&model.UserRoles{}, "user_id = ?", userId).Error)
	}

	tx := db.Begin()
	defer impl.commitOrRollBack(ctx, tx, &err)

	err = tx.Delete(&model.UserRoles{}, "user_id = ?", userId).Error
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}

	userRoles := make([]model.UserRoles, 0, rolesLen)
	for _, role := range mRoles {
		userRoles = append(userRoles, model.UserRoles{UserId: userId, RoleId: role.ID})
	}
	return impl.handleUpdateError(ctx, tx.Create(&userRoles).Error)
}

func (impl *adminImpl) UpdateRole(ctx context.Context, adminId uint64, role service.Role) error {
	db := impl.initializedConf.db.WithContext(ctx)
	err := impl.innerAuthQuery(ctx, db, adminId, service.AdminGroupId, updateFlag)
	if err != nil {
		return err
	}

	roleGroupId := impl.initializedConf.nameToGroupId[role.Group.Name]
	if roleGroupId == service.PublicGroupId {
		// right on public part are not updatable
		return common.ErrUpdate
	}

	actionFlags := convertActionsToFlags(role.Actions)
	if actionFlags == 0 {
		// delete unused role
		nameSubQuery := db.Model(&model.RoleName{}).Select("id").Where("name = ?", role.Name)
		var mRole model.Role
		if err = db.First(&mRole, "name_id IN (?) AND object_id = ?", nameSubQuery, roleGroupId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
			return servicecommon.ErrInternal
		}

		if err = db.Delete(&model.Role{}, mRole.ID).Error; err != nil {
			impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
			return servicecommon.ErrInternal
		}

		// we delete the names without roles
		roleSubQuery := db.Model(&model.Role{}).Distinct("name_id")
		if err = db.Delete(&model.RoleName{}, "id NOT IN (?)", roleSubQuery).Error; err != nil {
			impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
			return servicecommon.ErrInternal
		}

		// invalidate the cache of name
		impl.idToNameMutex.Lock()
		impl.idToName = map[uint64]string{}
		impl.idToNameMutex.Unlock()
		return nil
	}

	tx := db.Begin()
	defer impl.commitOrRollBack(ctx, tx, &err)

	var roleName model.RoleName
	if err = tx.FirstOrCreate(&roleName, model.RoleName{Name: role.Name}).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}

	var mRole model.Role
	err = db.First(&mRole, "name_id = ? AND object_id = ?", roleName.ID, roleGroupId).Error
	if err == nil {
		return impl.handleUpdateError(ctx, tx.Model(&role).Update("action_flags", actionFlags).Error)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}

	mRole = model.Role{NameId: roleName.ID, ObjectId: roleGroupId, ActionFlags: actionFlags}
	return impl.handleUpdateError(ctx, tx.Create(&mRole).Error)
}

func (impl *adminImpl) GetUserRoles(ctx context.Context, adminId uint64, userId uint64) ([]service.Role, error) {
	db := impl.initializedConf.db.WithContext(ctx)
	if adminId == userId {
		return impl.getUserRoles(ctx, db, userId)
	}

	err := impl.innerAuthQuery(ctx, db, adminId, service.AdminGroupId, accessFlag)
	if err != nil {
		return nil, err
	}
	return impl.getUserRoles(ctx, db, userId)
}

func (impl *adminImpl) ViewUserRoles(ctx context.Context, adminId uint64, userId uint64) (bool, []service.Role, error) {
	db := impl.initializedConf.db.WithContext(ctx)
	adminRoles, err := impl.retrieveUserRoles(ctx, db, adminId, service.AdminGroupId)
	if err != nil {
		return false, nil, err
	}

	if adminId == userId {
		updateRight := impl.evalOPA(ctx, userId, service.AdminGroupId, updateFlag, adminRoles) == nil
		userRoles, err := impl.getUserRoles(ctx, db, userId)
		return updateRight, userRoles, err
	}

	err = impl.evalOPA(ctx, userId, service.AdminGroupId, accessFlag, adminRoles)
	if err != nil {
		return false, nil, err
	}

	updateRight := impl.evalOPA(ctx, userId, service.AdminGroupId, updateFlag, adminRoles) == nil
	userRoles, err := impl.getUserRoles(ctx, db, userId)
	return updateRight, userRoles, err
}

func (impl *adminImpl) EditUserRoles(ctx context.Context, adminId uint64, userId uint64) ([]service.Role, []service.Role, error) {
	db := impl.initializedConf.db.WithContext(ctx)
	allRoles, err := impl.getAllRoles(ctx, db, adminId)
	if err != nil {
		return nil, nil, err
	}

	userRoles, err := impl.getUserRoles(ctx, db, userId)
	return userRoles, allRoles, err
}

func (impl *adminImpl) getUserRoles(ctx context.Context, db *gorm.DB, userId uint64) ([]service.Role, error) {
	subQuery := db.Model(&model.UserRoles{}).Select("role_id").Where("user_id = ?", userId)

	var roles []model.Role
	err := db.Find(&roles, "id IN (?)", subQuery).Error
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return impl.convertRolesFromModel(ctx, db, roles)
}

func (impl *adminImpl) getAllRoles(ctx context.Context, db *gorm.DB, adminId uint64) ([]service.Role, error) {
	err := impl.innerAuthQuery(ctx, db, adminId, service.AdminGroupId, accessFlag)
	if err != nil {
		return nil, err
	}

	var roles []model.Role
	if err := db.Find(&roles, "object_id IN ?", impl.initializedConf.groupIds).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return impl.convertRolesFromModel(ctx, db, roles)
}

func (impl *adminImpl) innerAuthQuery(ctx context.Context, db *gorm.DB, userId uint64, groupId uint64, actionFlag uint8) error {
	userRoles, err := impl.retrieveUserRoles(ctx, db, userId, groupId)
	if err != nil {
		return err
	}
	return impl.evalOPA(ctx, userId, groupId, actionFlag, userRoles)
}

func (impl *adminImpl) retrieveUserRoles(ctx context.Context, db *gorm.DB, userId uint64, groupId uint64) ([]any, error) {
	var roles []model.Role
	if userId != 0 {
		subQuery := db.Model(&model.UserRoles{}).Select("role_id").Where("user_id = ?", userId)
		if err := db.Find(&roles, "id in (?)", subQuery).Error; err != nil {
			impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
			return nil, servicecommon.ErrInternal
		}
	}
	return convertDataFromRolesModel(roles), nil
}

func (impl *adminImpl) evalOPA(ctx context.Context, userId uint64, groupId uint64, actionFlag uint8, userRoles []any) error {
	input := map[string]any{
		"userId": userId, "objectId": groupId, "actionFlag": actionFlag, "userRoles": userRoles,
	}
	results, err := impl.initializedConf.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		impl.Logger(ctx).Error("OPA evaluation failed", common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	if !results.Allowed() {
		return common.ErrNotAuthorized
	}
	return nil
}

func (impl *adminImpl) loadRoles(ctx context.Context, db *gorm.DB, roles []service.Role) ([]model.Role, error) {
	resRoles := make([]model.Role, 0, len(roles)) // probably lot more space than necessary
	for name, objectIdSet := range impl.extractNamesToObjectIdSet(roles) {
		subQuery := db.Model(&model.RoleName{}).Select("id").Where("name = ?", name)

		var tempRoles []model.Role
		err := db.Find(&tempRoles, "name_id IN (?) AND object_id IN ?", subQuery, objectIdSet.Slice()).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
			return nil, servicecommon.ErrInternal
		}
		if len(tempRoles) != 0 {
			resRoles = append(resRoles, tempRoles...)
		}
	}
	return resRoles, nil
}

func (impl *adminImpl) convertRolesFromModel(ctx context.Context, db *gorm.DB, roles []model.Role) ([]service.Role, error) {
	allThere := true
	resRoles := make([]service.Role, 0, len(roles))
	impl.idToNameMutex.RLock()
	for _, role := range roles {
		var name string
		id := role.NameId
		name, allThere = impl.idToName[id]
		if !allThere {
			break
		}
		resRoles = append(resRoles, impl.convertRoleFromModel(name, role))
	}
	impl.idToNameMutex.RUnlock()
	if allThere {
		return resRoles, nil
	}

	impl.idToNameMutex.Lock()
	defer impl.idToNameMutex.Unlock()
	allThere = true
	resRoles = resRoles[:0]
	missingIdSet := common.MakeSet[uint64](nil)
	for _, role := range roles {
		id := role.NameId
		name, ok := impl.idToName[id]
		if ok {
			resRoles = append(resRoles, impl.convertRoleFromModel(name, role))
		} else {
			allThere = false
			missingIdSet.Add(id)
		}
	}
	if allThere {
		return resRoles, nil
	}

	var roleNames []model.RoleName
	if err := db.Find(&roleNames, "id IN ?", missingIdSet.Slice()).Error; err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	for _, roleName := range roleNames {
		impl.idToName[roleName.ID] = roleName.Name
	}

	resRoles = resRoles[:0]
	for _, role := range roles {
		resRoles = append(resRoles, impl.convertRoleFromModel(impl.idToName[role.NameId], role))
	}
	return resRoles, nil
}

func (impl *adminImpl) commitOrRollBack(ctx context.Context, tx *gorm.DB, err *error) {
	if r := recover(); r != nil {
		tx.Rollback()
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, "recover", r)
	} else if *err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}
}

func (impl *adminImpl) handleUpdateError(ctx context.Context, err error) error {
	if err != nil {
		impl.Logger(ctx).Error(servicecommon.DBAccessMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func convertDataFromRolesModel(roles []model.Role) []any {
	res := make([]any, 0, len(roles))
	for _, role := range roles {
		res = append(res, map[string]any{
			"objectId":    role.ObjectId,
			"actionFlags": role.ActionFlags,
		})
	}
	return res
}

func (impl *adminImpl) convertRoleFromModel(name string, role model.Role) service.Role {
	return service.Role{
		Name: name, Group: service.Group{
			Id: role.ObjectId, Name: impl.initializedConf.groupIdToName[role.ObjectId],
		},
		Actions: convertActionsFromFlags(role.ActionFlags),
	}
}

func convertActionsFromFlags(actionFlags uint8) []string {
	resActions := make([]string, 0, 4)
	if actionFlags&accessFlag != 0 {
		resActions = append(resActions, service.ActionAccess)
	}
	if actionFlags&createFlag != 0 {
		resActions = append(resActions, service.ActionCreate)
	}
	if actionFlags&updateFlag != 0 {
		resActions = append(resActions, service.ActionUpdate)
	}
	if actionFlags&deleteFlag != 0 {
		resActions = append(resActions, service.ActionDelete)
	}
	return resActions
}

func (impl *adminImpl) extractNamesToObjectIdSet(roles []service.Role) map[string]common.Set[uint64] {
	nameToObjectIdSet := map[string]common.Set[uint64]{}
	for _, role := range roles {
		objectIdSet := nameToObjectIdSet[role.Name]
		if objectIdSet == nil {
			objectIdSet = common.MakeSet[uint64](nil)
			nameToObjectIdSet[role.Name] = objectIdSet
		}
		objectIdSet.Add(impl.initializedConf.nameToGroupId[role.Group.Name])
	}
	return nameToObjectIdSet
}

func convertActionsToFlags(actions []string) uint8 {
	var flags uint8
	for _, action := range actions {
		flags |= convertActionToFlag(action)
	}
	return flags
}

func convertActionToFlag(action string) uint8 {
	switch action {
	case service.ActionAccess:
		return accessFlag
	case service.ActionCreate:
		return createFlag
	case service.ActionUpdate:
		return updateFlag
	case service.ActionDelete:
		return deleteFlag
	}
	return accessFlag
}
