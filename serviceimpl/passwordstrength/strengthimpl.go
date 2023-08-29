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

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type PasswordStrengthService service.PasswordStrengthService

type strengthImpl struct {
	weaver.Implements[PasswordStrengthService]
	weaver.WithConfig[strengthConf]
	initializedConf initializedStrengthConf
}

func (impl *strengthImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf, err = initStrengthConf(impl.Logger(ctx), impl.Config())
	return
}

func (impl *strengthImpl) Validate(ctx context.Context, password string) error {
	logger := impl.Logger(ctx)
	err := passwordvalidator.Validate(password, impl.initializedConf.minEntropy)
	if err != nil {
		logger.Error("Password not validated", common.ErrorKey, err)
		return common.ErrWeakPassword
	}
	return nil
}

func (impl *strengthImpl) GetRules(ctx context.Context, lang string) (string, error) {
	logger := impl.Logger(ctx)
	description, ok := impl.initializedConf.localizedRules[lang]
	if !ok {
		logger.Error("Locale not found")
		return "", servicecommon.ErrInternal
	}
	return description, nil
}
