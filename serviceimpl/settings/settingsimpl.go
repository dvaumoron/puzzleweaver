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

package settingsimpl

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	mongoclient "github.com/dvaumoron/puzzleweaver/client/mongo"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "settings"

const userIdKey = "userId"
const settingsKey = collectionName // currently the same

var optsOnlySettingsField = options.FindOne().SetProjection(bson.D{{Key: settingsKey, Value: true}})
var optsCreateUnexisting = options.Replace().SetUpsert(true)

type SettingsService service.SettingsService

type settingsImpl struct {
	weaver.Implements[SettingsService]
	weaver.WithConfig[settingsConf]
	initializedConf initializedSettingsConf
}

func (impl *settingsImpl) Init(ctx context.Context) error {
	impl.initializedConf = initSettingsConf(impl.Config())
	return nil
}

func (impl *settingsImpl) Get(ctx context.Context, id uint64) (map[string]string, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)
	var result bson.D
	err = collection.FindOne(
		ctx, bson.D{{Key: userIdKey, Value: id}}, optsOnlySettingsField,
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	// call [1] to get picture because result has only the id and one field
	return mongoclient.ExtractStringMap(result[1].Value), nil
}

func (impl *settingsImpl) Update(ctx context.Context, id uint64, info map[string]string) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	settings := bson.M{userIdKey: id, settingsKey: info}
	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)
	_, err = collection.ReplaceOne(
		ctx, bson.D{{Key: userIdKey, Value: id}}, settings, optsCreateUnexisting,
	)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	return nil
}
