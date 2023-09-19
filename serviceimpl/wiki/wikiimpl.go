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

package wikiimpl

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	mongoclient "github.com/dvaumoron/puzzleweaver/client/mongo"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "pages"

const wikiIdKey = "wikiId"
const wikiRefKey = "ref"
const versionKey = "version"
const textKey = "text"
const userIdKey = "userId"

var descVersion = bson.D{{Key: versionKey, Value: -1}}
var contentFields = bson.D{
	// exclude unused fields
	{Key: wikiIdKey, Value: false}, {Key: wikiRefKey, Value: false}, {Key: userIdKey, Value: false},
}
var optsContentMaxVersion = options.FindOne().SetSort(descVersion).SetProjection(contentFields)
var optsContentFields = options.FindOne().SetProjection(contentFields)
var optsVersion = options.Find().SetProjection(
	bson.D{{Key: versionKey, Value: true}, {Key: userIdKey, Value: true}},
)

type remoteWikiImpl struct {
	weaver.Implements[RemoteWikiService]
	weaver.WithConfig[wikiConf]
	initializedConf initializedWikiConf
}

func (impl *remoteWikiImpl) Init(ctx context.Context) error {
	impl.initializedConf = initWikiConf(impl.Config())
	return nil
}

func (impl *remoteWikiImpl) Load(ctx context.Context, wikiId uint64, wikiRef string, version uint64) (RawWikiContent, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return RawWikiContent{}, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	filters := bson.D{
		{Key: wikiIdKey, Value: wikiId}, {Key: wikiRefKey, Value: wikiRef},
	}

	opts := optsContentMaxVersion
	if version != 0 {
		filters = append(filters, bson.E{Key: versionKey, Value: version})
		opts = optsContentFields
	}

	var result bson.M
	err = collection.FindOne(ctx, filters, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// an empty Content has Version 0, which is recognized by client
			return RawWikiContent{}, nil
		}

		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return RawWikiContent{}, servicecommon.ErrInternal
	}
	return convertToContent(result), nil
}

func (impl *remoteWikiImpl) Store(ctx context.Context, wikiId uint64, userId uint64, wikiRef string, last uint64, markdown string) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	// rely on the mongo server to ensure there will be no duplicate
	page := bson.M{
		wikiIdKey: wikiId, wikiRefKey: wikiRef, versionKey: last + 1,
		userIdKey: userId, textKey: markdown,
	}

	if _, err = collection.InsertOne(ctx, page); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return common.ErrBaseVersion
		}

		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func (impl *remoteWikiImpl) GetVersions(ctx context.Context, wikiId uint64, wikiRef string) ([]RawWikiContent, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.D{
		{Key: wikiIdKey, Value: wikiId}, {Key: wikiRefKey, Value: wikiRef},
	}, optsVersion)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return nil, servicecommon.ErrInternal
	}
	return mongoclient.ConvertSlice(results, convertToVersion), nil
}

func (impl *remoteWikiImpl) Delete(ctx context.Context, wikiId uint64, wikiRef string, version uint64) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	_, err = collection.DeleteMany(ctx, bson.D{
		{Key: wikiIdKey, Value: wikiId}, {Key: wikiRefKey, Value: wikiRef},
		{Key: versionKey, Value: version},
	})
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func convertToContent(page bson.M) RawWikiContent {
	text, _ := page[textKey].(string)
	return RawWikiContent{
		Version:   mongoclient.ExtractUint64(page[versionKey]),
		CreatedAt: mongoclient.ExtractCreateDate(page).Unix(),
		Markdown:  text,
	}
}

func convertToVersion(page bson.M) RawWikiContent {
	return RawWikiContent{
		Version:   mongoclient.ExtractUint64(page[versionKey]),
		CreatorId: mongoclient.ExtractUint64(page[userIdKey]),
		CreatedAt: mongoclient.ExtractCreateDate(page).Unix(),
	}
}
