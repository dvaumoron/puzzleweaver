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
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/gin-gonic/gin"
)

type profileWidget struct {
	defaultHandler        gin.HandlerFunc
	viewHandler           gin.HandlerFunc
	editHandler           gin.HandlerFunc
	saveHandler           gin.HandlerFunc
	changeLoginHandler    gin.HandlerFunc
	changePasswordHandler gin.HandlerFunc
	pictureHandler        gin.HandlerFunc
}

func defaultRedirecter(c *gin.Context) string {
	userId := GetSessionUserId(c)
	if userId == 0 {
		return common.DefaultErrorRedirect(GetLogger(c), unknownUserKey)
	}
	return profileUrlBuilder(userId).String()
}

func (w profileWidget) LoadInto(router gin.IRouter) {
	router.GET("/", w.defaultHandler)
	router.GET("/view/:UserId", w.viewHandler)
	router.GET("/edit", w.editHandler)
	router.POST("/save", w.saveHandler)
	router.POST("/changeLogin", w.changeLoginHandler)
	router.POST("/changePassword", w.changePasswordHandler)
	router.GET("/picture/:UserId", w.pictureHandler)
}

func newProfilePage(globalConfig *config.GlobalServiceConfig) Page {
	profileService := globalConfig.ProfileService
	adminService := globalConfig.AdminService
	loginService := globalConfig.LoginService

	p := MakeHiddenPage("profile")
	p.Widget = profileWidget{
		defaultHandler: common.CreateRedirect(defaultRedirecter),
		viewHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			ctx := c.Request.Context()
			viewedUserId := GetRequestedUserId(c)
			if viewedUserId == 0 {
				return "", common.DefaultErrorRedirect(GetLogger(c), common.ErrorTechnicalKey)
			}

			currentUserId, _ := data[common.UserIdName].(uint64)
			updateRight := viewedUserId == currentUserId
			if !updateRight {
				if err := profileService.ViewRight(ctx, currentUserId); err != nil {
					return "", common.DefaultErrorRedirect(GetLogger(c), err.Error())
				}
			}

			profiles, err := profileService.GetProfiles(ctx, []uint64{viewedUserId})
			if err != nil {
				return "", common.DefaultErrorRedirect(GetLogger(c), err.Error())
			}

			groups, err := adminService.GetUserRoles(ctx, currentUserId, viewedUserId)
			if err == nil {
				data["UserRight"] = displayGroups(groups)
			} else if err != common.ErrNotAuthorized {
				// ignore ErrNotAuthorized
				return "", common.DefaultErrorRedirect(GetLogger(c), err.Error())
			}

			userProfile := profiles[viewedUserId]
			data[common.AllowedToUpdateName] = updateRight
			data[common.ViewedUserName] = userProfile
			return "profile/view", ""
		}),
		editHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			userId, _ := data[common.UserIdName].(uint64)
			if userId == 0 {
				return "", common.DefaultErrorRedirect(GetLogger(c), unknownUserKey)
			}

			profiles, err := profileService.GetProfiles(c.Request.Context(), []uint64{userId})
			if err != nil {
				return "", common.DefaultErrorRedirect(GetLogger(c), err.Error())
			}

			userProfile := profiles[userId]
			data[common.ViewedUserName] = userProfile
			return "profile/edit", ""
		}),
		saveHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			logger := GetLogger(c)
			userId := GetSessionUserId(c)
			if userId == 0 {
				return common.DefaultErrorRedirect(logger, unknownUserKey)
			}

			desc := c.PostForm("userDesc")
			info := c.PostFormMap("userInfo")

			picture, err := c.FormFile("picture")
			if err != nil {
				logger.Error("Failed to retrieve picture file", common.ErrorKey, err)
				return common.DefaultErrorRedirect(logger, common.ErrorTechnicalKey)
			}

			if picture != nil {
				var pictureFile multipart.File
				pictureFile, err = picture.Open()
				if err != nil {
					logger.Error("Failed to open retrieve picture file ", common.ErrorKey, err)
					return common.DefaultErrorRedirect(logger, common.ErrorTechnicalKey)
				}
				defer pictureFile.Close()

				var pictureData []byte
				pictureData, err = io.ReadAll(pictureFile)
				if err != nil {
					logger.Error("Failed to read picture file ", common.ErrorKey, err)
					return common.DefaultErrorRedirect(logger, common.ErrorTechnicalKey)
				}

				err = profileService.UpdatePicture(ctx, userId, pictureData)
			}

			if err == nil {
				err = profileService.UpdateProfile(ctx, userId, desc, info)
			}

			targetBuilder := profileUrlBuilder(userId)
			if err != nil {
				common.WriteError(targetBuilder, logger, err.Error())
			}
			return targetBuilder.String()
		}),
		changeLoginHandler: common.CreateRedirect(func(c *gin.Context) string {
			session := GetSession(c)
			userId := GetSessionUserId(c)
			if userId == 0 {
				return common.DefaultErrorRedirect(GetLogger(c), unknownUserKey)
			}

			oldLogin := session.Load(loginName)
			newLogin := c.PostForm(loginName)
			password := c.PostForm(passwordName)

			err := common.ErrEmptyLogin
			if newLogin != "" {
				err = loginService.ChangeLogin(c.Request.Context(), userId, oldLogin, newLogin, password)
			}

			targetBuilder := profileUrlBuilder(userId)
			if err == nil {
				session.Store(loginName, newLogin)
			} else {
				common.WriteError(targetBuilder, GetLogger(c), err.Error())
			}
			return targetBuilder.String()
		}),
		changePasswordHandler: common.CreateRedirect(func(c *gin.Context) string {
			session := GetSession(c)
			userId := GetSessionUserId(c)
			if userId == 0 {
				return common.DefaultErrorRedirect(GetLogger(c), unknownUserKey)
			}

			login := session.Load(loginName)
			oldPassword := c.PostForm("oldPassword")
			newPassword := c.PostForm("newPassword")
			confirmPassword := c.PostForm(confirmPasswordName)

			err := common.ErrEmptyPassword
			if newPassword != "" {
				err = common.ErrWrongConfirm
				if newPassword == confirmPassword {
					err = loginService.ChangePassword(c.Request.Context(), userId, login, oldPassword, newPassword)
				}
			}

			targetBuilder := profileUrlBuilder(userId)
			if err != nil {
				common.WriteError(targetBuilder, GetLogger(c), err.Error())
			}
			return targetBuilder.String()
		}),
		pictureHandler: func(c *gin.Context) {
			userId := GetRequestedUserId(c)
			if userId == 0 {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			data := profileService.GetPicture(c.Request.Context(), userId)
			c.Data(http.StatusOK, http.DetectContentType(data), data)
		},
	}
	return p
}

func GetRequestedUserId(c *gin.Context) uint64 {
	userId, err := strconv.ParseUint(c.Param(userIdName), 10, 64)
	if err != nil {
		GetLogger(c).Warn("Failed to parse userId from request", common.ErrorKey, err)
	}
	return userId
}

func profileUrlBuilder(userId uint64) *strings.Builder {
	targetBuilder := new(strings.Builder)
	targetBuilder.WriteString("/profile/view/")
	targetBuilder.WriteString(strconv.FormatUint(userId, 10))
	return targetBuilder
}