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
	"strconv"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/dvaumoron/puzzleweaver/web/locale"
	"github.com/gin-gonic/gin"
)

const userIdName = "UserId"
const loginName = "Login" // current connected user login
const passwordName = "Password"
const confirmPasswordName = "ConfirmPassword"
const loginUrlName = "LoginUrl"
const prevUrlWithErrorName = "PrevUrlWithError"

type loginWidget struct {
	displayHandler gin.HandlerFunc
	submitHandler  gin.HandlerFunc
	logoutHandler  gin.HandlerFunc
}

func (w loginWidget) LoadInto(router gin.IRouter) {
	router.GET("/", w.displayHandler)
	router.POST("/submit", w.submitHandler)
	router.GET("/logout", w.logoutHandler)
}

func newLoginPage(loginService service.LoginService, settingsManager *SettingsManager) Page {
	p := MakeHiddenPage("login")
	p.Widget = loginWidget{
		displayHandler: CreateTemplate(func(data gin.H, c *gin.Context) (string, string) {
			data[common.RedirectName] = c.Query(common.RedirectName)

			currentUrl := c.Request.URL
			var errorKey string
			if len(currentUrl.Query()) == 0 {
				errorKey = common.QueryError
			} else {
				errorKey = "&error="
			}
			data[prevUrlWithErrorName] = currentUrl.String() + errorKey

			// To hide the connection link
			delete(data, loginUrlName)

			return "login", ""
		}),
		submitHandler: common.CreateRedirect(func(c *gin.Context) string {
			ctx := c.Request.Context()
			logger := GetLogger(c)
			login := c.PostForm(loginName)
			password := c.PostForm(passwordName)
			register := c.PostForm("Register") == "true"

			if login == "" {
				return c.PostForm(prevUrlWithErrorName) + common.ErrorEmptyLoginKey
			}
			if password == "" {
				return c.PostForm(prevUrlWithErrorName) + common.ErrorEmptyPasswordKey
			}

			var userId uint64
			var err error
			if register {
				if c.PostForm(confirmPasswordName) != password {
					return c.PostForm(prevUrlWithErrorName) + common.ErrorWrongConfirmPasswordKey
				}

				userId, err = loginService.Register(ctx, login, password)
			} else {
				userId, err = loginService.Verify(ctx, login, password)
			}

			if err != nil {
				return c.PostForm(prevUrlWithErrorName) + common.FilterErrorMsg(logger, err.Error())
			}

			s := GetSession(c)
			s.Store(loginName, login)
			s.Store(userIdName, strconv.FormatUint(userId, 10))

			GetLocalesManager(c).SetLangCookie(logger, settingsManager.Get(userId, c)[locale.LangName], c)

			return c.PostForm(common.RedirectName)
		}),
		logoutHandler: common.CreateRedirect(func(c *gin.Context) string {
			s := GetSession(c)
			s.Delete(loginName)
			s.Delete(userIdName)
			return c.Query(common.RedirectName)
		}),
	}
	return p
}
