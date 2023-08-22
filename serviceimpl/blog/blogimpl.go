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

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/remoteservice"
)

// check matching with interface
var _ remoteservice.RemoteBlogService = &remoteBlogImpl{}

type remoteBlogImpl struct {
	weaver.Implements[remoteservice.RemoteBlogService]
}

func (*remoteBlogImpl) CreatePost(ctx context.Context, blogId uint64, userId uint64, title string, content string) (uint64, error) {
	return 0, nil
}

func (*remoteBlogImpl) GetPost(ctx context.Context, blogId uint64, postId uint64) (remoteservice.RawBlogPost, error) {
	return remoteservice.RawBlogPost{}, nil
}

func (*remoteBlogImpl) GetPosts(ctx context.Context, blogId uint64, start uint64, end uint64, filter string) (uint64, []remoteservice.RawBlogPost, error) {
	return 0, nil, nil
}

func (*remoteBlogImpl) Delete(ctx context.Context, blogId uint64, postId uint64) error {
	return nil
}
