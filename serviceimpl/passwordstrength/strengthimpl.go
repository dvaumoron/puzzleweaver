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

package passwordstrengthimpl

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/exp/slog"
)

type PasswordStrengthService service.PasswordStrengthService

type strengthImpl struct {
	weaver.Implements[PasswordStrengthService]
	weaver.WithConfig[strengthConf]
	confMutex       sync.RWMutex
	initializedConf *initializedStrengthConf
}

func (impl *strengthImpl) getInitializedConf(logger *slog.Logger) *initializedStrengthConf {
	impl.confMutex.RLock()
	initializedConf := impl.initializedConf
	impl.confMutex.RUnlock()
	if initializedConf != nil {
		return initializedConf
	}

	impl.confMutex.Lock()
	defer impl.confMutex.Unlock()
	if impl.initializedConf == nil {
		impl.initializedConf = initStrengthConf(logger, impl.Config())
	}
	return impl.initializedConf
}

func (impl *strengthImpl) Validate(ctx context.Context, password string) error {
	logger := impl.Logger(ctx)
	err := passwordvalidator.Validate(password, impl.getInitializedConf(logger).minEntropy)
	if err != nil {
		logger.Error("Password not validated", common.ErrorKey, err)
		return common.ErrWeakPassword
	}
	return nil
}

func (impl *strengthImpl) GetRules(ctx context.Context, lang string) (string, error) {
	logger := impl.Logger(ctx)
	description, ok := impl.getInitializedConf(logger).localizedRules[lang]
	if !ok {
		logger.Error("Locale not found")
		return "", servicecommon.ErrInternal
	}
	return description, nil
}
