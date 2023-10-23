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
	"context"
	"errors"
	"strings"

	settingsimpl "github.com/dvaumoron/puzzleweaver/serviceimpl/settings"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/locale"
	"github.com/gin-gonic/gin"
)

const settingsName = "Settings"

var errWrongLang = errors.New(common.ErrorWrongLangKey)

type SettingsManager struct {
	settingsService settingsimpl.SettingsService
	InitSettings    func(*gin.Context) map[string]string
	CheckSettings   func(map[string]string, *gin.Context) error
}

func NewSettingsManager(settingsService settingsimpl.SettingsService) *SettingsManager {
	return &SettingsManager{settingsService: settingsService, InitSettings: initSettings, CheckSettings: checkSettings}
}

func initSettings(c *gin.Context) map[string]string {
	return map[string]string{locale.LangName: GetLocalesManager(c).GetLang(GetLogger(c), c)}
}

func checkSettings(settings map[string]string, c *gin.Context) error {
	askedLang := settings[locale.LangName]
	lang := GetLocalesManager(c).SetLangCookie(GetLogger(c), askedLang, c)
	settings[locale.LangName] = lang
	if lang != askedLang {
		return errWrongLang
	}
	return nil
}

func (m *SettingsManager) Get(userId uint64, c *gin.Context) map[string]string {
	ctx := c.Request.Context()
	userSettings := c.GetStringMapString(settingsName)
	if len(userSettings) != 0 {
		return userSettings
	}

	userSettings, err := m.settingsService.Get(ctx, userId)
	if err != nil {
		GetLogger(c).Warn("Failed to retrieve user settings", common.ErrorKey, err)
	}

	if len(userSettings) == 0 {
		userSettings = m.InitSettings(c)
		err = m.settingsService.Update(ctx, userId, userSettings)
		if err != nil {
			GetLogger(c).Warn("Failed to create user settings", common.ErrorKey, err)
		}
	}
	c.Set(settingsName, userSettings)
	return userSettings
}

func (m *SettingsManager) Update(ctx context.Context, userId uint64, settings map[string]string) error {
	return m.settingsService.Update(ctx, userId, settings)
}

type settingsWidget struct {
	editHandler gin.HandlerFunc
	saveHandler gin.HandlerFunc
}

func (w settingsWidget) LoadInto(router gin.IRouter) {
	router.GET("/", w.editHandler)
	router.POST("/save", w.saveHandler)
}

func newSettingsPage(settingsManager *SettingsManager) Page {
	p := MakeHiddenPage("settings")
	p.Widget = settingsWidget{
		editHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			userId, _ := data[common.UserIdName].(uint64)
			if userId == 0 {
				return "", common.DefaultErrorRedirect(GetLogger(c), unknownUserKey)
			}

			data["Settings"] = settingsManager.Get(userId, c)
			return "settings/edit", ""
		}),
		saveHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			userId := GetSessionUserId(c)
			if userId == 0 {
				return common.DefaultErrorRedirect(GetLogger(c), unknownUserKey)
			}

			settings := c.PostFormMap("settings")
			err := settingsManager.CheckSettings(settings, c)
			if err == nil {
				err = settingsManager.Update(ctx, userId, settings)
			}

			var targetBuilder strings.Builder
			targetBuilder.WriteString("/settings")
			if err != nil {
				common.WriteError(&targetBuilder, GetLogger(c), err.Error())
			}
			return targetBuilder.String()
		}),
	}
	return p
}
