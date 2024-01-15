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
	"log/slog"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweb/common"
	"github.com/redis/go-redis/v9"
)

// this key maintains the existence of the session when there is no other data,
// but it is never send to client nor updated by it
const creationTimeName = "sessionCreationTime"

var errGenerateRetry = errors.New("generate reached maximum number of retries")

type sessionImpl struct {
	weaver.Implements[SessionService]
	weaver.WithConfig[sessionConf]
	initializedConf initializedSessionConf
	generateMutex   sync.Mutex
}

func (impl *sessionImpl) Init(ctx context.Context) (err error) {
	impl.initializedConf, err = initSessionConf(impl.Logger(ctx), impl.Config())
	return
}

func (impl *sessionImpl) updateWithDefaultTTL(ctx context.Context, logger *slog.Logger, id string) {
	if err := impl.initializedConf.rdb.Expire(ctx, id, impl.Config().SessionTimeout).Err(); err != nil {
		logger.Info("Failed to set TTL", common.ErrorKey, err)
	}
}

func (impl *sessionImpl) Generate(ctx context.Context) (uint64, error) {
	logger := impl.Logger(ctx)

	// avoid id clash when generating, but possible bottleneck
	impl.generateMutex.Lock()
	defer impl.generateMutex.Unlock()
	for i := 0; i < impl.Config().RetryNumber; i++ {
		id := rand.Uint64()
		idStr := strconv.FormatUint(id, 10)
		nb, err := impl.initializedConf.rdb.Exists(ctx, idStr).Result()
		if err != nil {
			logger.Error(servicecommon.RedisCallMsg, common.ErrorKey, err)
			return 0, servicecommon.ErrInternal
		}
		if nb == 0 {
			err := impl.initializedConf.rdb.HSet(ctx, idStr, creationTimeName, time.Now().String()).Err()
			if err != nil {
				logger.Error(servicecommon.RedisCallMsg, common.ErrorKey, err)
				return 0, servicecommon.ErrInternal
			}
			impl.updateWithDefaultTTL(ctx, logger, idStr)
			return id, nil
		}
	}
	return 0, errGenerateRetry
}

func (impl *sessionImpl) Get(ctx context.Context, id uint64) (map[string]string, error) {
	logger := impl.Logger(ctx)

	idStr := strconv.FormatUint(id, 10)
	info, err := impl.initializedConf.rdb.HGetAll(ctx, idStr).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		logger.Error(servicecommon.RedisCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	impl.updateWithDefaultTTL(ctx, logger, idStr)
	delete(info, creationTimeName)
	return info, nil
}

func (impl *sessionImpl) Update(ctx context.Context, id uint64, info map[string]string) error {
	logger := impl.Logger(ctx)

	infoCopy := map[string]any{}
	var keyToDelete []string
	for k, v := range info {
		if k == creationTimeName {
			continue
		} else if v == "" {
			keyToDelete = append(keyToDelete, k)
		} else {
			infoCopy[k] = v
		}
	}
	idStr := strconv.FormatUint(id, 10)
	if err := impl.initializedConf.updater(impl.initializedConf.rdb, ctx, idStr, keyToDelete, infoCopy); err != nil {
		logger.Error(servicecommon.RedisCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	impl.updateWithDefaultTTL(ctx, logger, idStr)
	return nil
}

func updateSessionInfoTx(rdb *redis.Client, ctx context.Context, id string, keyToDelete []string, info map[string]any) error {
	haveActions := false
	pipe := rdb.TxPipeline()
	if len(keyToDelete) != 0 {
		haveActions = true
		pipe.HDel(ctx, id, keyToDelete...)
	}
	if len(info) != 0 {
		haveActions = true
		pipe.HSet(ctx, id, info)
	}
	if haveActions {
		_, err := pipe.Exec(ctx)
		return err
	}
	return nil
}

func updateSessionInfo(rdb *redis.Client, ctx context.Context, id string, keyToDelete []string, info map[string]any) error {
	if len(keyToDelete) != 0 {
		if err := rdb.HDel(ctx, id, keyToDelete...).Err(); err != nil {
			return err
		}
	}
	if len(info) != 0 {
		return rdb.HSet(ctx, id, info).Err()
	}
	return nil
}
