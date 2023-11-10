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

package adminclient

import (
	"context"

	adminimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/admin"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	adminservice "github.com/dvaumoron/puzzleweb/admin/service"
)

type adminServiceWrapper struct {
	adminService adminimpl.AdminService
}

func MakeAdminServiceWrapper(adminService adminimpl.AdminService) adminservice.AdminService {
	return adminServiceWrapper{adminService: adminService}
}

func (client adminServiceWrapper) AuthQuery(ctx context.Context, userId uint64, groupId uint64, action string) error {
	return client.adminService.AuthQuery(ctx, userId, groupId, action)
}

func (client adminServiceWrapper) GetAllGroups(ctx context.Context, adminId uint64) ([]adminservice.Group, error) {
	groups, err := client.adminService.GetAllGroups(ctx, adminId)
	return servicecommon.ConvertSlice(groups, convertGroupFrom), err
}

func (client adminServiceWrapper) GetActions(ctx context.Context, adminId uint64, roleName string, groupName string) ([]string, error) {
	return client.adminService.GetActions(ctx, adminId, roleName, groupName)
}

func (client adminServiceWrapper) UpdateUser(ctx context.Context, adminId uint64, userId uint64, roles []adminservice.Group) error {
	return client.adminService.UpdateUser(ctx, adminId, userId, servicecommon.ConvertSlice(roles, convertGroupTo))
}

func (client adminServiceWrapper) UpdateRole(ctx context.Context, adminId uint64, roleName string, groupName string, actions []string) error {
	return client.adminService.UpdateRole(ctx, adminId, roleName, groupName, actions)
}

func (client adminServiceWrapper) GetUserRoles(ctx context.Context, adminId uint64, userId uint64) ([]adminservice.Group, error) {
	groups, err := client.adminService.GetUserRoles(ctx, adminId, userId)
	return servicecommon.ConvertSlice(groups, convertGroupFrom), err
}

func (client adminServiceWrapper) ViewUserRoles(ctx context.Context, adminId uint64, userId uint64) (bool, []adminservice.Group, error) {
	b, groups, err := client.adminService.ViewUserRoles(ctx, adminId, userId)
	return b, servicecommon.ConvertSlice(groups, convertGroupFrom), err
}

func (client adminServiceWrapper) EditUserRoles(ctx context.Context, adminId uint64, userId uint64) ([]adminservice.Group, []adminservice.Group, error) {
	groups1, groups2, err := client.adminService.EditUserRoles(ctx, adminId, userId)
	return servicecommon.ConvertSlice(groups1, convertGroupFrom), servicecommon.ConvertSlice(groups2, convertGroupFrom), err
}

func convertGroupFrom(group adminimpl.Group) adminservice.Group {
	return adminservice.Group{Id: group.Id, Name: group.Name, Roles: servicecommon.ConvertSlice(group.Roles, convertRoleFrom)}
}

func convertGroupTo(group adminservice.Group) adminimpl.Group {
	return adminimpl.Group{Id: group.Id, Name: group.Name, Roles: servicecommon.ConvertSlice(group.Roles, convertRoleTo)}
}

func convertRoleFrom(role adminimpl.Role) adminservice.Role {
	return adminservice.Role{Name: role.Name, Actions: role.Actions}
}

func convertRoleTo(role adminservice.Role) adminimpl.Role {
	return adminimpl.Role{Name: role.Name, Actions: role.Actions}
}
