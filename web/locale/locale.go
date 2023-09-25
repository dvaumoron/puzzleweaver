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

package locale

import (
	"log/slog"
	"unicode"

	"github.com/dvaumoron/puzzleweaver/web/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const LangName = "lang"
const pathName = "Path"

type Manager interface {
	GetDefaultLang() string
	GetAllLang() []string
	GetMultipleLang() bool
	GetLang(*slog.Logger, *gin.Context) string
	CheckLang(*slog.Logger, string) string
	SetLangCookie(*slog.Logger, string, *gin.Context) string
}

type localesManager struct {
	SessionTimeOut int
	Domain         string
	AllLang        []string
	DefaultLang    string
	MultipleLang   bool
	matcher        language.Matcher
}

func NewManager(logger *slog.Logger, globalConfig *config.GlobalServiceConfig) Manager {
	sessionTimeOut := globalConfig.SessionTimeOut
	domain := globalConfig.Domain
	allLang := globalConfig.AllLang

	size := len(allLang)
	if size == 0 {
		logger.Error("No locales declared")
	}

	tags := make([]language.Tag, 0, size)
	for _, lang := range allLang {
		tags = append(tags, language.MustParse(lang))
	}

	return &localesManager{
		SessionTimeOut: sessionTimeOut, Domain: domain, AllLang: allLang,
		DefaultLang: globalConfig.DefaultLang, MultipleLang: size > 1, matcher: language.NewMatcher(tags),
	}
}

func (m *localesManager) GetDefaultLang() string {
	return m.DefaultLang
}

func (m *localesManager) GetAllLang() []string {
	return m.AllLang
}

func (m *localesManager) GetMultipleLang() bool {
	return m.MultipleLang
}

func (m *localesManager) GetLang(logger *slog.Logger, c *gin.Context) string {
	lang, err := c.Cookie(LangName)
	if err != nil {
		tag, _ := language.MatchStrings(m.matcher, c.GetHeader("Accept-Language"))
		return m.setLangCookie(tag.String(), c)
	}
	// check & refresh cookie
	return m.SetLangCookie(logger, lang, c)
}

func (m *localesManager) CheckLang(logger *slog.Logger, lang string) string {
	for _, l := range m.AllLang {
		if lang == l {
			return lang
		}
	}
	logger.Info("Asked not declared locale", "askedLocale", lang)
	return m.DefaultLang
}

func (m *localesManager) setLangCookie(lang string, c *gin.Context) string {
	c.SetCookie(LangName, lang, m.SessionTimeOut, "/", m.Domain, false, false)
	return lang
}

func (m *localesManager) SetLangCookie(logger *slog.Logger, lang string, c *gin.Context) string {
	return m.setLangCookie(m.CheckLang(logger, lang), c)
}

func CamelCase(word string) string {
	if word == "" {
		return ""
	}

	first := true
	chars := make([]rune, 0, len(word))
	for _, char := range word {
		if first {
			first = false
			char = unicode.ToTitle(char)
		}
		chars = append(chars, char)
	}
	return string(chars)
}
