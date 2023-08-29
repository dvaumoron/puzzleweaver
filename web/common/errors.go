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
const QueryError = "?" + ErrorKey + "="

// error displayed to user
const (
	ErrorBaseVersionKey          = "BaseVersionOutdated"
	ErrorEmptyLoginKey           = "EmptyLogin"
	ErrorEmptyPasswordKey        = "EmptyPassword"
	ErrorExistingLoginKey        = "ExistingLogin"
	ErrorNotAuthorizedKey        = "ErrorNotAuthorized"
	ErrorTechnicalKey            = "ErrorTechnicalProblem"
	ErrorUpdateKey               = "ErrorUpdate"
	ErrorWeakPasswordKey         = "WeakPassword"
	ErrorWrongConfirmPasswordKey = "WrongConfirmPassword"
	ErrorWrongLangKey            = "WrongLang"
	ErrorWrongLoginKey           = "WrongLogin"
)

const originalErrorMsg = "Original error"

var (
	ErrBaseVersion   = errors.New(ErrorBaseVersionKey)
	ErrEmptyLogin    = errors.New(ErrorEmptyLoginKey)
	ErrEmptyPassword = errors.New(ErrorEmptyPasswordKey)
	ErrExistingLogin = errors.New(ErrorExistingLoginKey)
	ErrNotAuthorized = errors.New(ErrorNotAuthorizedKey)
	ErrTechnical     = errors.New(ErrorTechnicalKey)
	ErrUpdate        = errors.New(ErrorUpdateKey)
	ErrWeakPassword  = errors.New(ErrorWeakPasswordKey)
	ErrWrongConfirm  = errors.New(ErrorWrongConfirmPasswordKey)
	ErrWrongLogin    = errors.New(ErrorWrongLangKey)
)

func LogOriginalError(logger *slog.Logger, err error) {
	logger.Warn(originalErrorMsg, ErrorKey, err.Error())
}

func WriteError(urlBuilder *strings.Builder, logger *slog.Logger, errorMsg string) {
	urlBuilder.WriteString(QueryError)
	urlBuilder.WriteString(FilterErrorMsg(logger, errorMsg))
}

func DefaultErrorRedirect(logger *slog.Logger, errorMsg string) string {
	return "/?error=" + FilterErrorMsg(logger, errorMsg)
}

func FilterErrorMsg(logger *slog.Logger, errorMsg string) string {
	if errorMsg == ErrorBaseVersionKey || errorMsg == ErrorEmptyLoginKey || errorMsg == ErrorEmptyPasswordKey ||
		errorMsg == ErrorExistingLoginKey || errorMsg == ErrorNotAuthorizedKey || errorMsg == ErrorTechnicalKey ||
		errorMsg == ErrorUpdateKey || errorMsg == ErrorWeakPasswordKey || errorMsg == ErrorWrongConfirmPasswordKey ||
		errorMsg == ErrorWrongLangKey || errorMsg == ErrorWrongLoginKey {
		return errorMsg
	}
	logger.Error(originalErrorMsg, ErrorKey, errorMsg)
	return ErrorTechnicalKey
}
