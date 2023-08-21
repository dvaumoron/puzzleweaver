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
const ErrorBaseVersionKey = "BaseVersionOutdated"
const ErrorNotAuthorizedKey = "ErrorNotAuthorized"
const ErrorTechnicalKey = "ErrorTechnicalProblem"
const ErrorUpdateKey = "ErrorUpdate"
const ErrorWrongLangKey = "WrongLang"

const originalErrorMsg = "Original error"

var ErrNotAuthorized = errors.New(ErrorNotAuthorizedKey)
var ErrTechnical = errors.New(ErrorTechnicalKey)
var ErrUpdate = errors.New(ErrorUpdateKey)

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
	if errorMsg == ErrorBaseVersionKey || errorMsg == ErrorNotAuthorizedKey || errorMsg == ErrorTechnicalKey || errorMsg == ErrorUpdateKey || errorMsg == ErrorWrongLangKey {
		return errorMsg
	}
	logger.Warn(originalErrorMsg, ErrorKey, errorMsg)
	return ErrorTechnicalKey
}
