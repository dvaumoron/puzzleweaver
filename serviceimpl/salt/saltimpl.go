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

package saltimpl

import (
	"context"
	"crypto/rand"

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

const redisCallMsg = "Failed during Redis call"
const generateMsg = "Failed to generate"

// check matching with interface
type SaltService service.SaltService

type saltImpl struct {
	weaver.Implements[SaltService]
	weaver.WithConfig[saltConf]
	initializedConf initializedSaltConf
}

func (impl *saltImpl) Init(ctx context.Context) error {
	impl.initializedConf = initSaltConf(impl.Logger(ctx), impl.Config())
	return nil
}

func (impl *saltImpl) LoadOrGenerate(ctx context.Context, logins ...string) ([][]byte, error) {
	logger := impl.Logger(ctx)
	salts := make([][]byte, 0, len(logins))
	for _, login := range logins {
		salt, err := impl.innerLoadOrGenerate(ctx, logger, login)
		if err != nil {
			return nil, err
		}
		salts = append(salts, salt)
	}
	return salts, nil
}

func (impl *saltImpl) innerLoadOrGenerate(ctx context.Context, logger *slog.Logger, login string) ([]byte, error) {
	rdb := impl.initializedConf.rdb
	salt, err := rdb.Get(ctx, login).Result()
	if err == nil {
		return []byte(salt), nil
	}
	if err != redis.Nil {
		logger.Error(redisCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	saltBuffer := make([]byte, impl.Config().SaltLen)
	_, err = rand.Read(saltBuffer)
	if err != nil {
		logger.Error(generateMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	salt = string(saltBuffer)
	if err = rdb.Set(ctx, login, salt, 0).Err(); err != nil {
		logger.Error(redisCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return saltBuffer, nil
}
