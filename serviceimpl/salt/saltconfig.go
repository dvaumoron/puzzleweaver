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
	redisclient "github.com/dvaumoron/puzzleweaver/client/redis"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type saltConf struct {
	RedisAddr     string
	RedisUser     string
	RedisPassword string
	RedisDBNum    int
	SaltLen       int
}

type initializedSaltConf struct {
	rdb *redis.Client
}

func initSaltConf(logger *slog.Logger, conf *saltConf) initializedSaltConf {
	rdb := redisclient.New(logger, &redis.Options{
		Addr:     conf.RedisAddr,
		Username: conf.RedisUser,
		Password: conf.RedisPassword,
		DB:       conf.RedisDBNum,
	})
	return initializedSaltConf{rdb: rdb}
}
