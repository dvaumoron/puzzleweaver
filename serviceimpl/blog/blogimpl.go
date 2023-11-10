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

package blogimpl

import (
	"context"
	"strings"

	"github.com/ServiceWeaver/weaver"
	mongoclient "github.com/dvaumoron/puzzleweaver/client/mongo"
	servicecommon "github.com/dvaumoron/puzzleweaver/serviceimpl/common"
	"github.com/dvaumoron/puzzleweb/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "posts"

const blogIdKey = "blogId"
const postIdKey = "postId"
const userIdKey = "userId"
const titleKey = "title"
const textKey = "text"

var optsMaxPostId = options.FindOne().SetSort(bson.D{{Key: postIdKey, Value: -1}}).SetProjection(bson.D{{Key: postIdKey, Value: true}})

type remoteBlogImpl struct {
	weaver.Implements[RemoteBlogService]
	weaver.WithConfig[blogConf]
	initializedConf initializedBlogConf
}

func (impl *remoteBlogImpl) Init(ctx context.Context) error {
	impl.initializedConf = initBlogConf(impl.Config())
	return nil
}

func (impl *remoteBlogImpl) CreatePost(ctx context.Context, blogId uint64, userId uint64, title string, content string) (uint64, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	filter := bson.D{{Key: blogIdKey, Value: blogId}}
	post := bson.M{
		blogIdKey: blogId, userIdKey: userId, titleKey: title, textKey: content,
	}

	// rely on the mongo server to ensure there will be no duplicate
	newPostId := uint64(1)

GeneratePostIdStep:
	var result bson.D
	err = collection.FindOne(ctx, filter, optsMaxPostId).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			goto CreatePostStep
		}

		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, servicecommon.ErrInternal
	}

	// call [1] to get postId because result has only the id and one field
	newPostId = mongoclient.ExtractUint64(result[1].Value) + 1

CreatePostStep:
	post[postIdKey] = newPostId
	if _, err = collection.InsertOne(ctx, post); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// retry
			goto GeneratePostIdStep
		}

		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, servicecommon.ErrInternal
	}
	return newPostId, nil
}

func (impl *remoteBlogImpl) GetPost(ctx context.Context, blogId uint64, postId uint64) (RawBlogPost, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return RawBlogPost{}, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	var result bson.M
	err = collection.FindOne(
		ctx, bson.D{{Key: blogIdKey, Value: blogId}, {Key: postIdKey, Value: postId}},
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("No blog post found with requested ids", "blogId", blogId, "postId", postId)
		} else {
			logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		}
		return RawBlogPost{}, servicecommon.ErrInternal
	}
	return convertToPost(result), nil
}

func (impl *remoteBlogImpl) GetPosts(ctx context.Context, blogId uint64, start uint64, end uint64, filter string) (uint64, []RawBlogPost, error) {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	filters := bson.D{{Key: blogIdKey, Value: blogId}}
	if filter != "" {
		filters = append(filters, bson.E{Key: titleKey, Value: buildRegexFilter(filter)})
	}

	total, err := collection.CountDocuments(ctx, filters)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}

	paginate := options.Find().SetSort(bson.D{{Key: postIdKey, Value: -1}})
	paginate.SetSkip(int64(start)).SetLimit(int64(end - start))

	cursor, err := collection.Find(ctx, filters, paginate)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return 0, nil, servicecommon.ErrInternal
	}
	return uint64(total), servicecommon.ConvertSlice(results, convertToPost), nil
}

func (impl *remoteBlogImpl) Delete(ctx context.Context, blogId uint64, postId uint64) error {
	logger := impl.Logger(ctx)
	client, err := mongo.Connect(ctx, impl.initializedConf.clientOptions)
	if err != nil {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return servicecommon.ErrInternal
	}
	defer mongoclient.Disconnect(client, ctx, logger)

	collection := client.Database(impl.Config().MongoDatabaseName).Collection(collectionName)

	_, err = collection.DeleteMany(
		ctx, bson.D{{Key: blogIdKey, Value: blogId}, {Key: postIdKey, Value: postId}},
	)
	if err != nil && err != mongo.ErrNoDocuments {
		logger.Error(servicecommon.MongoCallMsg, common.ErrorKey, err)
		return common.ErrUpdate
	}
	return nil
}

func convertToPost(post bson.M) RawBlogPost {
	title, _ := post[titleKey].(string)
	text, _ := post[textKey].(string)
	return RawBlogPost{
		Id: mongoclient.ExtractUint64(post[postIdKey]), CreatorId: mongoclient.ExtractUint64(post[userIdKey]),
		CreatedAt: mongoclient.ExtractCreateDate(post).Unix(), Title: title, Content: text,
	}
}

func buildRegexFilter(filter string) bson.D {
	filter = strings.ReplaceAll(filter, "%", ".*")
	var regexBuilder strings.Builder
	if strings.Index(filter, ".*") != 0 {
		regexBuilder.WriteString(".*")
	}
	regexBuilder.WriteString(filter)
	if strings.LastIndex(filter, ".*") != len(filter)-2 {
		regexBuilder.WriteString(".*")
	}
	return bson.D{{Key: "$regex", Value: regexBuilder.String()}}
}
