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

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const QueryError = "?error="

const WrongLangKey = "WrongLang"

const ReportingPlaceName = "reporting_place"

// error displayed to user
const ErrorNotAuthorizedKey = "ErrorNotAuthorized"
const ErrorTechnicalKey = "ErrorTechnicalProblem"
const ErrorUpdateKey = "ErrorUpdate"

var ErrNotAuthorized = errors.New(ErrorNotAuthorizedKey)
var ErrTechnical = errors.New(ErrorTechnicalKey)
var ErrUpdate = errors.New(ErrorUpdateKey)

func LogOriginalError(logger otelzap.LoggerWithCtx, err error) error {
	logger.WithOptions(zap.AddCallerSkip(1)).Warn("Original error", zap.Error(err))
	return ErrTechnical
}

func WriteError(urlBuilder *strings.Builder, errorMsg string) {
	urlBuilder.WriteString(QueryError)
	urlBuilder.WriteString(errorMsg)
}

func DefaultErrorRedirect(errorMsg string) string {
	return "/?error=" + errorMsg
}
