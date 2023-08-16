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

package web

import (
	"errors"
	"sort"
	"strings"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/dvaumoron/puzzleweaver/web/locale"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

const roleNameName = "RoleName"
const groupName = "Group"
const groupsName = "Groups"
const viewAdminName = "ViewAdmin"

const (
	accessKey = "AccessLabel"
	createKey = "CreateLabel"
	updateKey = "UpdateLabel"
	deleteKey = "DeleteLabel"
)

var errBadName = errors.New("ErrorBadRoleName")

type GroupDisplay struct {
	Id           uint64
	Name         string
	DisplayName  string
	Roles        []RoleDisplay
	AddableRoles []RoleDisplay
}

func NewGroupDisplay(id uint64, name string) *GroupDisplay {
	return &GroupDisplay{Id: id, Name: name, DisplayName: getGroupDisplayNameKey(name)}
}

type RoleDisplay struct {
	Name    string
	Actions []string
}

func MakeRoleDisplay(role service.Role) RoleDisplay {
	return RoleDisplay{Name: role.Name, Actions: displayActions(role.Actions)}
}

type sortableGroups []*GroupDisplay

func (s sortableGroups) Len() int {
	return len(s)
}

func (s sortableGroups) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s sortableGroups) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortableRoles []RoleDisplay

func (s sortableRoles) Len() int {
	return len(s)
}

func (s sortableRoles) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s sortableRoles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type adminWidget struct {
	displayHandler    gin.HandlerFunc
	listUserHandler   gin.HandlerFunc
	viewUserHandler   gin.HandlerFunc
	editUserHandler   gin.HandlerFunc
	saveUserHandler   gin.HandlerFunc
	deleteUserHandler gin.HandlerFunc
	listRoleHandler   gin.HandlerFunc
	editRoleHandler   gin.HandlerFunc
	saveRoleHandler   gin.HandlerFunc
}

func (w adminWidget) LoadInto(router gin.IRouter) {
	router.GET("/", w.displayHandler)
	router.GET("/user/list", w.listUserHandler)
	router.GET("/user/view/:UserId", w.viewUserHandler)
	router.GET("/user/edit/:UserId", w.editUserHandler)
	router.POST("/user/save/:UserId", w.saveUserHandler)
	router.GET("/user/delete/:UserId", w.deleteUserHandler)
	router.GET("/role/list", w.listRoleHandler)
	router.GET("/role/edit/:RoleName/:Group", w.editRoleHandler)
	router.POST("/role/save", w.saveRoleHandler)
}

func newAdminPage(adminService service.AdminService, userService service.AdvancedUserService, profileService service.AdvancedProfileService, adminConfig config.AdminConfig) Page {
	defaultPageSize := adminConfig.PageSize

	p := MakeHiddenPage("admin")
	p.Widget = adminWidget{
		displayHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			viewAdmin, _ := data[viewAdminName].(bool)
			if !viewAdmin {
				return "", common.DefaultErrorRedirect(common.ErrorNotAuthorizedKey)
			}
			return "admin/index", ""
		}),
		listUserHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			viewAdmin, _ := data[viewAdminName].(bool)
			if !viewAdmin {
				return "", common.DefaultErrorRedirect(common.ErrorNotAuthorizedKey)
			}

			pageNumber, start, end, filter := common.GetPagination(defaultPageSize, c)

			total, users, err := userService.ListUsers(c.Request.Context(), start, end, filter)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			common.InitPagination(data, filter, pageNumber, end, total)
			data["Users"] = users
			InitNoELementMsg(data, len(users), c)
			return "admin/user/list", ""
		}),
		viewUserHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			ctx := c.Request.Context()
			adminId, _ := data[common.IdName].(uint64)
			userId := GetRequestedUserId(c)
			if userId == 0 {
				return "", common.DefaultErrorRedirect(common.ErrorTechnicalKey)
			}

			roles, err := adminService.GetUserRoles(ctx, adminId, userId)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			users, err := userService.GetUsers(ctx, []uint64{userId})
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			updateRight := adminService.AuthQuery(ctx, adminId, service.AdminGroupId, service.ActionUpdate) == nil

			user := users[userId]
			data[common.ViewedUserName] = user
			data[common.AllowedToUpdateName] = updateRight
			data[groupsName] = DisplayGroups(roles)
			return "admin/user/view", ""
		}),
		editUserHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			ctx := c.Request.Context()
			adminId, _ := data[common.IdName].(uint64)
			userId := GetRequestedUserId(c)
			if userId == 0 {
				return "", common.DefaultErrorRedirect(common.ErrorTechnicalKey)
			}

			allRoles, err := adminService.GetAllRoles(ctx, adminId)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			userRoles, err := adminService.GetUserRoles(ctx, adminId, userId)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			userIdToLogin, err := userService.GetUsers(ctx, []uint64{userId})
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			data[common.ViewedUserName] = userIdToLogin[userId]
			data[groupsName] = displayEditGroups(userRoles, allRoles)
			return "admin/user/edit", ""
		}),
		saveUserHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			userId := GetRequestedUserId(c)
			err := common.ErrTechnical
			if userId != 0 {
				rolesStr := c.PostFormArray("roles")
				roles := make([]service.Role, 0, len(rolesStr))
				for _, roleStr := range rolesStr {
					splitted := strings.Split(roleStr, "/")
					if len(splitted) > 1 {
						roles = append(roles, service.Role{Name: splitted[0], GroupName: splitted[1]})
					}
				}
				err = adminService.UpdateUser(ctx, GetSessionUserId(c), userId, roles)
			}

			targetBuilder := userListUrlBuilder()
			if err != nil {
				common.WriteError(targetBuilder, err.Error())
			}
			return targetBuilder.String()
		}),
		deleteUserHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			userId := GetRequestedUserId(c)
			err := common.ErrTechnical
			if userId != 0 {
				// an empty slice delete the user right
				// only the first service call do a right check
				err = adminService.UpdateUser(ctx, GetSessionUserId(c), userId, []service.Role{})
				if err == nil {
					err = profileService.Delete(ctx, userId)
					if err == nil {
						err = userService.Delete(ctx, userId)
					}
				}
			}

			targetBuilder := userListUrlBuilder()
			if err != nil {
				common.WriteError(targetBuilder, err.Error())
			}
			return targetBuilder.String()
		}),
		listRoleHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			ctx := c.Request.Context()
			adminId, _ := data[common.IdName].(uint64)
			allRoles, err := adminService.GetAllRoles(ctx, adminId)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}

			allGroups, err := adminService.GetAllGroups(ctx)
			if err != nil {
				return "", common.DefaultErrorRedirect(err.Error())
			}
			data[groupsName] = displayAllGroups(allGroups, allRoles)
			return "admin/role/list", ""
		}),
		editRoleHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			roleName := c.Param(roleNameName)
			group := c.Param(groupName)

			data[roleNameName] = roleName
			data[groupName] = group
			data["GroupDisplayName"] = getGroupDisplayNameKey(group)

			if roleName != "new" {
				adminId, _ := data[common.IdName].(uint64)
				actions, err := adminService.GetActions(c.Request.Context(), adminId, roleName, group)
				if err != nil {
					return "", common.DefaultErrorRedirect(err.Error())
				}

				actionSet := common.MakeSet(actions)
				setActionChecked(data, actionSet, service.ActionAccess, "Access")
				setActionChecked(data, actionSet, service.ActionCreate, "Create")
				setActionChecked(data, actionSet, service.ActionUpdate, "Update")
				setActionChecked(data, actionSet, service.ActionDelete, "Delete")
			}

			return "admin/role/edit", ""
		}),
		saveRoleHandler: common.CreateRedirect(func(c *gin.Context) string {
			roleName := c.PostForm(roleNameName)
			err := errBadName
			if roleName != "new" {
				group := c.PostForm(groupName)
				actions := c.PostFormArray("actions")
				err = adminService.UpdateRole(c.Request.Context(), GetSessionUserId(c), service.Role{
					Name: roleName, GroupName: group, Actions: actions,
				})
			}

			var targetBuilder strings.Builder
			targetBuilder.WriteString("/admin/role/list")
			if err != nil {
				common.WriteError(&targetBuilder, err.Error())
			}
			return targetBuilder.String()
		}),
	}
	return p
}

func getGroupDisplayNameKey(name string) string {
	return "GroupLabel" + locale.CamelCase(name)
}

func DisplayGroups(roles []service.Role) []*GroupDisplay {
	nameToGroup := map[string]*GroupDisplay{}
	populateGroup(nameToGroup, roles, rolesAppender)
	return sortGroups(nameToGroup)
}

func populateGroup(nameToGroup map[string]*GroupDisplay, roles []service.Role, appender func(*GroupDisplay, service.Role)) {
	for _, role := range roles {
		groupName := role.GroupName
		group := nameToGroup[groupName]
		if group == nil {
			group = NewGroupDisplay(role.GroupId, groupName)
			nameToGroup[groupName] = group
		}
		appender(group, role)
	}
}

func rolesAppender(group *GroupDisplay, role service.Role) {
	group.Roles = append(group.Roles, MakeRoleDisplay(role))
}

// convert a string slice of codes in a displayable key slice,
// always in the same order : access, create, update, delete
func displayActions(actions []string) []string {
	actionSet := common.MakeSet(actions)
	res := make([]string, len(actions))
	if actionSet.Contains(service.ActionAccess) {
		res = append(res, accessKey)
	}
	if actionSet.Contains(service.ActionCreate) {
		res = append(res, createKey)
	}
	if actionSet.Contains(service.ActionUpdate) {
		res = append(res, updateKey)
	}
	if actionSet.Contains(service.ActionDelete) {
		res = append(res, deleteKey)
	}
	return res
}

func sortGroups(nameToGroup map[string]*GroupDisplay) []*GroupDisplay {
	groupRoles := common.MapToValueSlice(nameToGroup)
	sort.Sort(sortableGroups(groupRoles))
	for _, group := range groupRoles {
		sort.Sort(sortableRoles(group.Roles))
		sort.Sort(sortableRoles(group.AddableRoles))
	}
	return groupRoles
}

func displayEditGroups(userRoles []service.Role, allRoles []service.Role) []*GroupDisplay {
	nameToGroup := map[string]*GroupDisplay{}
	populateGroup(nameToGroup, userRoles, rolesAppender)
	populateGroup(nameToGroup, allRoles, addableRolesAppender)
	return sortGroups(nameToGroup)
}

func addableRolesAppender(group *GroupDisplay, role service.Role) {
	// check if the user already have this role
	contains := slices.ContainsFunc(group.Roles, func(roleDisplay RoleDisplay) bool {
		return roleDisplay.Name == role.Name
	})
	// no duplicate
	if !contains {
		group.AddableRoles = append(group.AddableRoles, MakeRoleDisplay(role))
	}
}

func displayAllGroups(groups []service.Group, roles []service.Role) []*GroupDisplay {
	nameToGroup := map[string]*GroupDisplay{}
	for _, group := range groups {
		nameToGroup[group.Name] = NewGroupDisplay(group.Id, group.Name)
	}
	populateGroup(nameToGroup, roles, rolesAppender)
	return sortGroups(nameToGroup)
}

func setActionChecked(data gin.H, actionSet common.Set[string], toTest string, name string) {
	if actionSet.Contains(toTest) {
		data[name] = true
	}
}

func userListUrlBuilder() *strings.Builder {
	targetBuilder := new(strings.Builder)
	targetBuilder.WriteString("/admin/user/list")
	return targetBuilder
}
