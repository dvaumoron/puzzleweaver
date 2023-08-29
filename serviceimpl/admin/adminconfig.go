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

	dbclient "github.com/dvaumoron/puzzleweaver/client/db"
	"github.com/dvaumoron/puzzleweaver/serviceimpl/admin/model"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/afero"
	"gorm.io/gorm"
)

type permissionGroup struct {
	Id   uint64
	Name string
}

type adminConf struct {
	PermissionGroups []permissionGroup
	DatabaseKind     string
	DatabaseAddress  string
	OPAModulePath    string
}

type initializedAdminConf struct {
	db            *gorm.DB
	query         rego.PreparedEvalQuery
	groupIdToName map[uint64]string
	nameToGroupId map[string]uint64
	groups        []service.Group
	groupIds      []uint64
}

func initAdminConf(ctx context.Context, conf *adminConf) (initializedAdminConf, error) {
	// TODO manage switch to network FS
	fileSystem := afero.NewOsFs()

	db, err := dbclient.New(conf.DatabaseKind, conf.DatabaseAddress)
	if err == nil {
		err = db.AutoMigrate(&model.UserRoles{}, &model.Role{}, &model.RoleName{})
	}
	if err != nil {
		return initializedAdminConf{}, err
	}

	query, err := readRule(ctx, fileSystem, conf.OPAModulePath)
	if err != nil {
		return initializedAdminConf{}, err
	}

	groupIdToName, nameToGroupId, groups, groupIds := initMapping(conf.PermissionGroups)
	return initializedAdminConf{
		db: db, query: query, groupIdToName: groupIdToName, nameToGroupId: nameToGroupId, groups: groups, groupIds: groupIds,
	}, nil
}

func readRule(ctx context.Context, fileSystem afero.Fs, modulePath string) (rego.PreparedEvalQuery, error) {
	data, err := afero.ReadFile(fileSystem, modulePath)
	if err != nil {
		return rego.PreparedEvalQuery{}, err
	}

	rule := rego.New(
		rego.Query("x = data.auth.allow"),
		rego.Module("auth.rego", string(data)),
	)
	return rule.PrepareForEval(ctx)
}

func initMapping(permissionGroups []permissionGroup) (map[uint64]string, map[string]uint64, []service.Group, []uint64) {
	groupIdToName := map[uint64]string{
		service.PublicGroupId: service.PublicName, service.AdminGroupId: service.AdminName,
	}
	nameToGroupId := map[string]uint64{
		service.PublicName: service.PublicGroupId, service.AdminName: service.AdminGroupId,
	}
	for _, idName := range permissionGroups {
		groupIdToName[idName.Id] = idName.Name
		nameToGroupId[idName.Name] = idName.Id
	}
	size := len(groupIdToName)
	groups := make([]service.Group, 0, size)
	groupIds := make([]uint64, 0, size)
	for id, name := range groupIdToName {
		groups = append(groups, service.Group{Id: id, Name: name})
		groupIds = append(groupIds, id)
	}
	return groupIdToName, nameToGroupId, groups, groupIds
}
