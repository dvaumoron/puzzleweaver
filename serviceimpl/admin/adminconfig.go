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
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

type permissionGroup struct {
	Id   uint64
	Name string
}

type adminConf struct {
	PermissionGroups []permissionGroup
}

type initializedAdminConf struct {
	groupIdToName map[uint64]string
	nameToGroupId map[string]uint64
}

func initAdminConf(conf *adminConf) initializedAdminConf {
	groupIdToName := map[uint64]string{
		service.PublicGroupId: service.PublicName, service.AdminGroupId: service.AdminName,
	}
	nameToGroupId := map[string]uint64{
		service.PublicName: service.PublicGroupId, service.AdminName: service.AdminGroupId,
	}
	for _, idName := range conf.PermissionGroups {
		groupIdToName[idName.Id] = idName.Name
		nameToGroupId[idName.Name] = idName.Id
	}
	return initializedAdminConf{groupIdToName: groupIdToName, nameToGroupId: nameToGroupId}
}
