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
	"sync"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

type AdminService service.AdminService

type adminImpl struct {
	weaver.Implements[AdminService]
	weaver.WithConfig[adminConf]
	confMutex       sync.RWMutex
	initializedConf *initializedAdminConf
}

func (impl *adminImpl) getInitializedConf() *initializedAdminConf {
	impl.confMutex.RLock()
	initializedConf := impl.initializedConf
	impl.confMutex.RUnlock()
	if initializedConf != nil {
		return initializedConf
	}

	impl.confMutex.Lock()
	defer impl.confMutex.Unlock()
	if impl.initializedConf == nil {
		impl.initializedConf = initAdminConf(impl.Config())
	}
	return impl.initializedConf
}

func (impl *adminImpl) getGroupId(ctx context.Context, groupName string) (uint64, error) {
	return impl.getInitializedConf().nameToGroupId[groupName], nil
}

func (impl *adminImpl) getGroupName(ctx context.Context, groupId uint64) (string, error) {
	return impl.getInitializedConf().groupIdToName[groupId], nil
}

func (impl *adminImpl) GetAllGroups(ctx context.Context) ([]service.Group, error) {
	groupIdToName := impl.getInitializedConf().groupIdToName
	groups := make([]service.Group, 0, len(groupIdToName))
	for id, name := range groupIdToName {
		groups = append(groups, service.Group{Id: id, Name: name})
	}
	return groups, nil
}

func (impl *adminImpl) AuthQuery(ctx context.Context, userId uint64, groupId uint64, action string) error {
	success := true
	// TODO
	if !success {
		return common.ErrNotAuthorized
	}
	return nil
}

func (impl *adminImpl) GetAllRoles(ctx context.Context, adminId uint64) ([]service.Role, error) {
	groupIdToName := impl.getInitializedConf().groupIdToName
	groupIds := make([]uint64, 0, len(groupIdToName))
	for groupId := range groupIdToName {
		groupIds = append(groupIds, groupId)
	}
	return impl.getGroupRoles(ctx, adminId, groupIds)
}

func (impl *adminImpl) GetActions(ctx context.Context, adminId uint64, roleName string, groupName string) ([]string, error) {
	success := true
	// TODO
	if !success {
		return nil, common.ErrNotAuthorized
	}
	// TODO
	return nil, nil
}

func (impl *adminImpl) UpdateUser(ctx context.Context, adminId uint64, userId uint64, roles []service.Role) error {
	success := true
	// TODO
	if !success {
		return common.ErrNotAuthorized
	}

	success2 := true
	// TODO
	if !success2 {
		return common.ErrUpdate
	}
	return nil
}

func (impl *adminImpl) UpdateRole(ctx context.Context, adminId uint64, role service.Role) error {
	success := true
	// TODO
	if !success {
		return common.ErrNotAuthorized
	}

	success2 := true
	// TODO
	if !success2 {
		return common.ErrUpdate
	}
	return nil
}

func (impl *adminImpl) GetUserRoles(ctx context.Context, adminId uint64, userId uint64) ([]service.Role, error) {

	if adminId == userId {
		return impl.getUserRoles(ctx, userId)
	}

	success := true
	// TODO
	if !success {
		return nil, common.ErrNotAuthorized
	}
	return impl.getUserRoles(ctx, userId)
}

func (impl *adminImpl) getGroupRoles(ctx context.Context, adminId uint64, groupIds []uint64) ([]service.Role, error) {
	success := true
	// TODO
	if !success {
		return nil, common.ErrNotAuthorized
	}
	// TODO
	return nil, nil
}

func (impl *adminImpl) getUserRoles(ctx context.Context, userId uint64) ([]service.Role, error) {
	// TODO
	return nil, nil
}
