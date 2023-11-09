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
package galleryimpl

import (
	"context"

	mongoclient "github.com/dvaumoron/puzzleweaver/client/mongo"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	galleryservice "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/gallery/service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "images"

	setOperator = "$set"

	galleryIdKey = "galleryId"
	imageIdKey   = "imageId"
	userIdKey    = "userId"
	titleKey     = "title"
	descKey      = "desc"
	imageKey     = "imageData"
)

var (
	optsCreateUnexisting     = options.Update().SetUpsert(true)
	optsMaxImageId           = options.FindOne().SetSort(bson.D{{Key: imageIdKey, Value: -1}}).SetProjection(bson.D{{Key: imageIdKey, Value: true}})
	optsOnlyImageField       = options.FindOne().SetProjection(bson.D{{Key: imageKey, Value: true}})
	optsOneExcludeImageField = options.FindOne().SetProjection(bson.D{{Key: imageKey, Value: false}})
)

type galleryImpl struct {
	clientOptions *options.ClientOptions
	databaseName  string
	loggerGetter  servicecommon.LoggerGetter
}

func New(serverAddress string, databaseName string, loggerGetter servicecommon.LoggerGetter) galleryservice.GalleryService {
	clientOptions := mongoclient.New(serverAddress)
	return galleryImpl{clientOptions: clientOptions, databaseName: databaseName, loggerGetter: loggerGetter}
}

func (impl galleryImpl) GetImages(ctx context.Context, galleryId uint64, start uint64, end uint64) (uint64, []galleryservice.GalleryImage, error) {
	logger := impl.loggerGetter.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.clientOptions)
	if err != nil {
		return 0, nil, err
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.databaseName).Collection(collectionName)
	filter := bson.D{{Key: galleryIdKey, Value: galleryId}}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	cursor, err := collection.Find(ctx, filter, initPaginationOpts(start, end))
	if err != nil {
		return 0, nil, err
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, nil, err
	}
	return uint64(total), mongoclient.ConvertSlice(results, convertToImage), nil
}

func (impl galleryImpl) GetImage(ctx context.Context, imageId uint64) (galleryservice.GalleryImage, error) {
	logger := impl.loggerGetter.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.clientOptions)
	if err != nil {
		return galleryservice.GalleryImage{}, err
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.databaseName).Collection(collectionName)

	var result bson.M
	err = collection.FindOne(
		ctx, bson.D{{Key: imageIdKey, Value: imageId}}, optsOneExcludeImageField,
	).Decode(&result)
	if err != nil {
		return galleryservice.GalleryImage{}, err
	}
	return convertToImage(result), nil
}

func (impl galleryImpl) GetImageData(ctx context.Context, imageId uint64) ([]byte, error) {
	logger := impl.loggerGetter.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.clientOptions)
	if err != nil {
		return nil, err
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.databaseName).Collection(collectionName)

	var result bson.D
	err = collection.FindOne(
		ctx, bson.D{{Key: imageIdKey, Value: imageId}}, optsOnlyImageField,
	).Decode(&result)
	if err != nil {
		return nil, err
	}

	// call [1] to get image because result has only the id and one field
	return mongoclient.ExtractBinary(result[1].Value), nil
}

func (impl galleryImpl) UpdateImage(ctx context.Context, galleryId uint64, info galleryservice.GalleryImage, data []byte) (uint64, error) {
	logger := impl.loggerGetter.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.clientOptions)
	if err != nil {
		return 0, err
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.databaseName).Collection(collectionName)

	imageId := info.ImageId
	image := bson.M{galleryIdKey: galleryId, imageIdKey: imageId, userIdKey: info.CreatorId, titleKey: info.Title, descKey: info.Desc}
	if len(data) != 0 {
		image[imageKey] = data
	}

	if imageId == 0 {
		return createImage(collection, ctx, image)
	}
	return imageId, updateImage(collection, ctx, image)
}

func (impl galleryImpl) DeleteImage(ctx context.Context, imageId uint64) error {
	logger := impl.loggerGetter.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.clientOptions)
	if err != nil {
		return err
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.databaseName).Collection(collectionName)

	_, err = collection.DeleteMany(
		ctx, bson.D{{Key: imageIdKey, Value: imageId}},
	)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func createImage(collection *mongo.Collection, ctx context.Context, image bson.M) (uint64, error) {
	// rely on the mongo server to ensure there will be no duplicate
	imageId := uint64(1)

	var err error
	var result bson.D
GenerateImageIdStep:
	err = collection.FindOne(ctx, bson.D{{Key: galleryIdKey, Value: image[galleryIdKey]}}, optsMaxImageId).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			goto CreateImageStep
		}

		return 0, err
	}

	// call [1] to get imageId because result has only the id and one field
	imageId = mongoclient.ExtractUint64(result[1].Value) + 1

CreateImageStep:
	image[imageIdKey] = imageId
	_, err = collection.InsertOne(ctx, image)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// retry
			goto GenerateImageIdStep
		}

		return 0, err
	}
	return imageId, nil
}

func updateImage(collection *mongo.Collection, ctx context.Context, image bson.M) error {
	request := bson.D{{Key: setOperator, Value: image}}
	_, err := collection.UpdateOne(
		ctx, bson.D{{Key: imageIdKey, Value: image[imageIdKey]}}, request, optsCreateUnexisting,
	)
	return err
}

func initPaginationOpts(start uint64, end uint64) *options.FindOptions {
	opts := options.Find().SetSort(bson.D{{Key: imageIdKey, Value: -1}}).SetProjection(bson.D{{Key: imageKey, Value: false}})
	castedStart := int64(start)
	return opts.SetSkip(castedStart).SetLimit(int64(end) - castedStart)
}

func convertToImage(image bson.M) galleryservice.GalleryImage {
	title, _ := image[titleKey].(string)
	desc, _ := image[descKey].(string)
	return galleryservice.GalleryImage{
		ImageId:   mongoclient.ExtractUint64(image[imageIdKey]),
		CreatorId: mongoclient.ExtractUint64(image[userIdKey]),
		Title:     title, Desc: desc,
	}
}
