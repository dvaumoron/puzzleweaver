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
package galleryservice

import "context"

type GalleryImage struct {
	ImageId   uint64
	CreatorId uint64
	Title     string
	Desc      string
}

type GalleryService interface {
	GetImages(ctx context.Context, galleryId uint64, start uint64, end uint64) (uint64, []GalleryImage, error)
	GetImage(ctx context.Context, imageId uint64) (GalleryImage, error)
	GetImageData(ctx context.Context, imageId uint64) ([]byte, error)
	UpdateImage(ctx context.Context, galleryId uint64, info GalleryImage, data []byte) (uint64, error)
	DeleteImage(ctx context.Context, imageId uint64) error
}
