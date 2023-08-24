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

package mongoclient

import (
	"context"
	"time"

	"github.com/dvaumoron/puzzleweaver/web/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"golang.org/x/exp/slog"
)

// Must be called after TracerProvider initialization.
func New(serverAddr string) *options.ClientOptions {
	clientOptions := options.Client()
	clientOptions.Monitor = otelmongo.NewMonitor()
	clientOptions.ApplyURI(serverAddr)
	return clientOptions
}

func Disconnect(client *mongo.Client, ctx context.Context, logger *slog.Logger) {
	if err := client.Disconnect(ctx); err != nil {
		logger.Error("Error during MongoDB disconnect", common.ErrorKey, err)
	}
}

func ExtractCreateDate(doc bson.M) time.Time {
	id, _ := doc["_id"].(primitive.ObjectID)
	return id.Timestamp()
}

func ExtractUint64(value any) uint64 {
	switch casted := value.(type) {
	case int32:
		return uint64(casted)
	case int64:
		return uint64(casted)
	}
	return 0
}

func ExtractBinary(value any) []byte {
	binary, _ := value.(primitive.Binary)
	return binary.Data
}

func ExtractStringMap(value any) map[string]string {
	resMap := map[string]string{}
	switch casted := value.(type) {
	case bson.D:
		for _, elem := range casted {
			resMap[elem.Key], _ = elem.Value.(string)
		}
	case bson.M:
		for key, innerValue := range casted {
			resMap[key], _ = innerValue.(string)
		}
	}
	return resMap
}

func ConvertSlice[T any](docs []bson.M, converter func(bson.M) T) []T {
	resSlice := make([]T, 0, len(docs))
	for _, doc := range docs {
		resSlice = append(resSlice, converter(doc))
	}
	return resSlice
}
