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

package redisclient

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

func New(logger *slog.Logger, options *redis.Options) (*redis.Client, error) {
	rdb := redis.NewClient(options)

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, err
	}

	// Enable metrics instrumentation.
	err := redisotel.InstrumentMetrics(rdb)
	return rdb, err
}
