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
	"time"

	redisclient "github.com/dvaumoron/puzzleweaver/client/redis"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type sessionConf struct {
	SessionTimeout time.Duration
	RetryNumber    int
	RedisAddress   string
	RedisUser      string
	RedisPassword  string
	RedisDBNum     int
	Debug          bool
}

type initializedSessionConf struct {
	rdb     *redis.Client
	updater func(*redis.Client, context.Context, string, []string, map[string]any) error
}

func initSessionConf(logger *slog.Logger, conf *sessionConf) (initializedSessionConf, error) {
	rdb, err := redisclient.New(logger, &redis.Options{
		Addr:     conf.RedisAddress,
		Username: conf.RedisUser,
		Password: conf.RedisPassword,
		DB:       conf.RedisDBNum,
	})
	if err != nil {
		return initializedSessionConf{}, err
	}

	updater := updateSessionInfoTx
	if conf.Debug {
		logger.Info("Mode debug on")
		updater = updateSessionInfo
	}
	return initializedSessionConf{rdb: rdb, updater: updater}, nil
}
