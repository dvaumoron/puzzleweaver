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

package common

import (
	"errors"
	"strings"

	"golang.org/x/exp/slog"
)

const ErrorKey = "error"

const QueryError = "?error="

// error displayed to user
const (
	ErrorBaseVersionKey   = "BaseVersionOutdated"
	ErrorExistingLoginKey = "ExistingLogin"
	ErrorNotAuthorizedKey = "ErrorNotAuthorized"
	ErrorTechnicalKey     = "ErrorTechnicalProblem"
	ErrorUpdateKey        = "ErrorUpdate"
	ErrorWeakPasswordKey  = "WeakPassword"
	ErrorWrongLangKey     = "WrongLang"
	ErrorWrongLoginKey    = "WrongLogin"
)

const originalErrorMsg = "Original error"

var (
	ErrBaseVersion   = errors.New(ErrorBaseVersionKey)
	ErrExistingLogin = errors.New(ErrorExistingLoginKey)
	ErrNotAuthorized = errors.New(ErrorNotAuthorizedKey)
	ErrTechnical     = errors.New(ErrorTechnicalKey)
	ErrUpdate        = errors.New(ErrorUpdateKey)
	ErrWeakPassword  = errors.New(ErrorWeakPasswordKey)
	ErrWrongLogin    = errors.New(ErrorWrongLangKey)
)

func LogOriginalError(logger *slog.Logger, err error) {
	logger.Warn(originalErrorMsg, ErrorKey, err.Error())
}

func WriteError(urlBuilder *strings.Builder, logger *slog.Logger, errorMsg string) {
	urlBuilder.WriteString(QueryError)
	urlBuilder.WriteString(filterErrorMsg(logger, errorMsg))
}

func DefaultErrorRedirect(logger *slog.Logger, errorMsg string) string {
	return "/?error=" + filterErrorMsg(logger, errorMsg)
}

func filterErrorMsg(logger *slog.Logger, errorMsg string) string {
	if errorMsg == ErrorBaseVersionKey || errorMsg == ErrorExistingLoginKey || errorMsg == ErrorNotAuthorizedKey || errorMsg == ErrorTechnicalKey || errorMsg == ErrorUpdateKey || errorMsg == ErrorWeakPasswordKey || errorMsg == ErrorWrongLangKey || errorMsg == ErrorWrongLoginKey {
		return errorMsg
	}
	logger.Warn(originalErrorMsg, ErrorKey, errorMsg)
	return ErrorTechnicalKey
}
