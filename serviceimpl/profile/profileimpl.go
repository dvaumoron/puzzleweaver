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

package profileimpl

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	mongoclient "github.com/dvaumoron/puzzleweaver/client/mongo"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "profiles"

const setOperator = "$set"

const userIdKey = "userId"
const descKey = "desc"
const infoKey = "info"
const pictureKey = "pictureData"

var optsCreateUnexisting = options.Update().SetUpsert(true)
var optsExcludePictureField = options.Find().SetProjection(bson.D{{Key: pictureKey, Value: false}})
var optsOnlyPictureField = options.FindOne().SetProjection(bson.D{{Key: pictureKey, Value: true}})

type RemoteProfileService remoteservice.RemoteProfileService

type remoteProfileImpl struct {
	weaver.Implements[RemoteProfileService]
	weaver.WithConfig[profileConf]
	initializedConf initializedProfileConf
}

func (impl *remoteProfileImpl) Init(ctx context.Context) error {
	impl.initializedConf = initProfileConf(impl.Config())
	return nil
}

func (impl *remoteProfileImpl) UpdateProfile(ctx context.Context, userId uint64, desc string, info map[string]string) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	infoB := bson.M{}
	for k, v := range info {
		infoB[k] = v
	}
	profile := bson.D{{Key: setOperator, Value: bson.M{userIdKey: userId, descKey: desc, infoKey: infoB}}}
	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)
	_, err = collection.UpdateOne(
		ctx, bson.D{{Key: userIdKey, Value: userId}}, profile, optsCreateUnexisting,
	)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *remoteProfileImpl) UpdatePicture(ctx context.Context, userId uint64, data []byte) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	profile := bson.D{{Key: setOperator, Value: bson.M{userIdKey: userId, pictureKey: data}}}
	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)
	_, err = collection.UpdateOne(
		ctx, bson.D{{Key: userIdKey, Value: userId}}, profile, optsCreateUnexisting,
	)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *remoteProfileImpl) GetPicture(ctx context.Context, userId uint64) ([]byte, error) {
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
		ctx, bson.D{{Key: userIdKey, Value: userId}}, optsOnlyPictureField,
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, servicecommon.ErrPictureNotFound
		}

		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	// call [1] to get picture because result has only the id and one field
	picture := mongoclient.ExtractBinary(result[1].Value)
	if len(picture) == 0 {
		return nil, servicecommon.ErrPictureNotFound
	}
	return picture, nil
}

func (impl *remoteProfileImpl) GetProfiles(ctx context.Context, userIds []uint64) (map[uint64]remoteservice.RawUserProfile, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)
	filter := bson.D{{Key: userIdKey, Value: bson.D{{Key: "$in", Value: userIds}}}}
	cursor, err := collection.Find(ctx, filter, optsExcludePictureField)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	profiles := map[uint64]remoteservice.RawUserProfile{}
	for _, profile := range results {
		userId := mongoclient.ExtractUint64(profile[userIdKey])
		desc, _ := profile[descKey].(string)
		profiles[userId] = remoteservice.RawUserProfile{Desc: desc, Info: mongoclient.ExtractStringMap(profile[infoKey])}
	}
	return profiles, nil
}

func (impl *remoteProfileImpl) Delete(ctx context.Context, userId uint64) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)
	_, err = collection.DeleteMany(ctx, bson.D{{Key: userIdKey, Value: userId}})
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}
