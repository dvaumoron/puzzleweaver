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

// check matching with interface
var _ service.AdminService = &AdminImpl{}

type AdminImpl struct {
	weaver.Implements[service.AdminService]
	weaver.WithConfig[adminConf]
	confMutex       sync.RWMutex
	initializedConf *initializedAdminConf
}

func (impl *AdminImpl) getInitializedConf(ctx context.Context) (*initializedAdminConf, error) {
	impl.confMutex.RLock()
	initializedConf := impl.initializedConf
	impl.confMutex.RUnlock()
	if initializedConf != nil {
		return initializedConf, nil
	}

	impl.confMutex.Lock()
	defer impl.confMutex.Unlock()
	if impl.initializedConf == nil {
		var err error
		impl.initializedConf, err = initAdminConf(impl.Config())
		if err != nil {
			return nil, err
		}
	}
	return impl.initializedConf, nil
}

func (impl *AdminImpl) getGroupId(ctx context.Context, groupName string) (uint64, error) {
	initializedConf, err := impl.getInitializedConf(ctx)
	if err != nil {
		return 0, err
	}
	return initializedConf.nameToGroupId[groupName], nil
}

func (impl *AdminImpl) getGroupName(ctx context.Context, groupId uint64) (string, error) {
	initializedConf, err := impl.getInitializedConf(ctx)
	if err != nil {
		return "", err
	}
	return initializedConf.groupIdToName[groupId], nil
}

func (impl *AdminImpl) GetAllGroups(ctx context.Context) ([]service.Group, error) {
	initializedConf, err := impl.getInitializedConf(ctx)
	if err != nil {
		return nil, err
	}

	groupIdToName := initializedConf.groupIdToName
	groups := make([]service.Group, 0, len(groupIdToName))
	for id, name := range groupIdToName {
		groups = append(groups, service.Group{Id: id, Name: name})
	}
	return groups, nil
}

func (impl *AdminImpl) AuthQuery(ctx context.Context, userId uint64, groupId uint64, action string) error {
	success := true
	// TODO
	if !success {
		return common.ErrNotAuthorized
	}
	return nil
}

func (impl *AdminImpl) GetAllRoles(ctx context.Context, adminId uint64) ([]service.Role, error) {
	initializedConf, err := impl.getInitializedConf(ctx)
	if err != nil {
		return nil, err
	}

	groupIdToName := initializedConf.groupIdToName
	groupIds := make([]uint64, 0, len(groupIdToName))
	for groupId := range groupIdToName {
		groupIds = append(groupIds, groupId)
	}
	return impl.getGroupRoles(ctx, adminId, groupIds)
}

func (impl *AdminImpl) GetActions(ctx context.Context, adminId uint64, roleName string, groupName string) ([]string, error) {
	success := true
	// TODO
	if !success {
		return nil, common.ErrNotAuthorized
	}
	// TODO
	return nil, nil
}

func (impl *AdminImpl) UpdateUser(ctx context.Context, adminId uint64, userId uint64, roles []service.Role) error {
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

func (impl *AdminImpl) UpdateRole(ctx context.Context, adminId uint64, role service.Role) error {
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

func (impl *AdminImpl) GetUserRoles(ctx context.Context, adminId uint64, userId uint64) ([]service.Role, error) {

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

func (impl *AdminImpl) getGroupRoles(ctx context.Context, adminId uint64, groupIds []uint64) ([]service.Role, error) {
	success := true
	// TODO
	if !success {
		return nil, common.ErrNotAuthorized
	}
	// TODO
	return nil, nil
}

func (impl *AdminImpl) getUserRoles(ctx context.Context, userId uint64) ([]service.Role, error) {
	// TODO
	return nil, nil
}
