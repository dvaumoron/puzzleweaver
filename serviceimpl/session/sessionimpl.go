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

package sessionimpl

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"golang.org/x/exp/slog"
)

// this key maintains the existence of the session when there is no other data,
// but it is never send to client nor updated by it
const creationTimeName = "sessionCreationTime"

const redisCallMsg = "Failed during Redis call"

var errInternal = errors.New("internal service error")

type SessionService service.SessionService

type sessionImpl struct {
	weaver.Implements[SessionService]
	weaver.WithConfig[sessionConf]
	confMutex       sync.RWMutex
	initializedConf *initializedSessionConf
	generateMutex   sync.Mutex
}

func (impl *sessionImpl) getInitializedConf(ctx context.Context, logger *slog.Logger) (*initializedSessionConf, error) {
	impl.confMutex.RLock()
	initializedConf := impl.initializedConf
	impl.confMutex.RUnlock()
	if initializedConf != nil {
		return initializedConf, nil
	}

	impl.confMutex.Lock()
	defer impl.confMutex.Unlock()
	if impl.initializedConf == nil {
		var err error
		impl.initializedConf, err = initSessionConf(impl.Config())
		if err != nil {
			logger.Error("Failed to init config", common.ErrorKey, err)
			return nil, err
		}
	}
	return impl.initializedConf, nil
}

func (impl *sessionImpl) updateWithDefaultTTL(ctx context.Context, logger *slog.Logger, id string) {
	if initializedConf, err := impl.getInitializedConf(ctx, logger); err == nil {
		if err = initializedConf.rdb.Expire(ctx, id, impl.Config().sessionTimeout).Err(); err != nil {
			logger.Info("Failed to set TTL", common.ErrorKey, err)
		}
	}
}

func (impl *sessionImpl) Generate(ctx context.Context) (uint64, error) {
	logger := impl.Logger(ctx)
	initializedConf, err := impl.getInitializedConf(ctx, logger)
	if err != nil {
		return 0, errInternal
	}

	// avoid id clash when generating, but possible bottleneck
	impl.generateMutex.Lock()
	defer impl.generateMutex.Unlock()
	for i := 0; i < impl.Config().retryNumber; i++ {
		id := rand.Uint64()
		idStr := strconv.FormatUint(id, 10)
		nb, err := initializedConf.rdb.Exists(ctx, idStr).Result()
		if err != nil {
			logger.Error(redisCallMsg, common.ErrorKey, err)
			return 0, errInternal
		}
		if nb == 0 {
			err := initializedConf.rdb.HSet(ctx, idStr, creationTimeName, time.Now().String()).Err()
			if err != nil {
				logger.Error(redisCallMsg, common.ErrorKey, err)
				return 0, errInternal
			}
			impl.updateWithDefaultTTL(ctx, logger, idStr)
			return id, nil
		}
	}
	return 0, errors.New("generate reached maximum number of retries")
}

func (impl *sessionImpl) Get(ctx context.Context, id uint64) (map[string]string, error) {
	// TODO
	return nil, nil
}

func (impl *sessionImpl) Update(ctx context.Context, id uint64, info map[string]string) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}
