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
	"strconv"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
	pb "github.com/dvaumoron/puzzlewikiservice"
	"go.uber.org/zap"
)

// check matching with interface
var _ wikiservice.WikiService = &wikiImpl{}

type wikiImpl struct {
	weaver.Implements[wikiservice.WikiService]
	authService    weaver.Ref[service.AuthService]
	profileService weaver.Ref[service.ProfileService]
	cache          *wikiCache
	wikiId         uint64
	groupId        uint64
	dateFormat     string
}

func (impl wikiImpl) LoadContent(ctx context.Context, userId uint64, lang string, title string, versionStr string) (*wikiservice.WikiContent, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionAccess)
	if err != nil {
		return nil, err
	}

	var version uint64
	if versionStr != "" {
		version, err = strconv.ParseUint(versionStr, 10, 64)
		if err != nil {
			impl.Logger(ctx).Info("Failed to parse wiki version, falling to last", common.ErrorKey, err)
		}
	}
	return impl.loadContent(ctx, buildRef(lang, title), version)
}

func (impl wikiImpl) StoreContent(ctx context.Context, userId uint64, lang string, title string, last string, markdown string) (bool, error) {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionCreate)
	if err != nil {
		return false, err
	}

	version, err := strconv.ParseUint(last, 10, 64)
	if err != nil {
		impl.Logger(ctx).Warn("Failed to parse wiki last version", common.ErrorKey, err)
		return false, common.ErrTechnical
	}
	return impl.storeContent(ctx, userId, buildRef(lang, title), version, markdown)
}

func (client wikiImpl) GetVersions(ctx context.Context, userId uint64, lang string, title string) ([]wikiservice.Version, error) {
	err := client.authService.Get().AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return nil, err
	}
	return client.getVersions(ctx, buildRef(lang, title))
}

func (impl wikiImpl) DeleteContent(ctx context.Context, userId uint64, lang string, title string, versionStr string) error {
	err := impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionDelete)
	if err != nil {
		return err
	}

	version, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		impl.Logger(ctx).Warn("Failed to parse wiki version to delete", zap.Error(err))
		return common.ErrTechnical
	}
	return impl.deleteContent(ctx, buildRef(lang, title), version)
}

func (impl wikiImpl) DeleteRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionDelete)
}

func (impl wikiImpl) loadContent(ctx context.Context, wikiRef string, version uint64) (*wikiservice.WikiContent, error) {

	if version != 0 {
		return impl.innerLoadContent(ctx, wikiRef, version)
	}

	list := []*pb.Version{}
	// TODO

	if lastVersion := maxVersion(list); lastVersion != nil {
		content := impl.cache.load(impl.Logger(ctx), wikiRef)
		if content != nil && lastVersion.Number == content.Version {
			return content, nil
		}
	}
	return impl.innerLoadContent(ctx, wikiRef, 0)
}

func (impl wikiImpl) innerLoadContent(ctx context.Context, wikiRef string, askedVersion uint64) (*wikiservice.WikiContent, error) {
	// TODO
	version := uint64(0)
	if version == 0 { // no stored wiki page
		return nil, nil
	}

	content := &wikiservice.WikiContent{Version: version, Markdown: "todo"}
	if askedVersion == 0 {
		impl.cache.store(impl.Logger(ctx), wikiRef, content)
	}
	return content, nil
}

func (impl wikiImpl) storeContent(ctx context.Context, userId uint64, wikiRef string, last uint64, markdown string) (bool, error) {
	version := last
	success := true
	// TODO
	if success {
		impl.cache.store(impl.Logger(ctx), wikiRef, &wikiservice.WikiContent{
			Version: version, Markdown: markdown,
		})
	}
	return success, nil
}

func (impl wikiImpl) getVersions(ctx context.Context, wikiRef string) ([]wikiservice.Version, error) {
	list := []*pb.Version{}
	// TODO
	return impl.sortConvertVersions(ctx, list)
}

func (impl wikiImpl) deleteContent(ctx context.Context, wikiRef string, version uint64) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}

	logger := impl.Logger(ctx)
	content := impl.cache.load(logger, wikiRef)
	if content != nil && version == content.Version {
		impl.cache.delete(logger, wikiRef)
	}
	return nil
}

func (impl wikiImpl) sortConvertVersions(ctx context.Context, list []*pb.Version) ([]wikiservice.Version, error) {
	size := len(list)
	if size == 0 {
		return nil, nil
	}

	valueSet := make([]*pb.Version, maxVersion(list).Number+1)
	// no duplicate check, there is one in GetProfiles
	userIds := make([]uint64, 0, size)
	for _, value := range list {
		valueSet[value.Number] = value
		userIds = append(userIds, value.UserId)
	}
	profiles, err := impl.profileService.Get().GetProfiles(ctx, userIds)
	if err != nil {
		return nil, err
	}

	newList := make([]wikiservice.Version, 0, size)
	for _, value := range valueSet {
		if value != nil {
			newList = append(newList, wikiservice.Version{Number: value.Number, Creator: profiles[value.UserId]})
		}
	}
	return newList, nil
}

func buildRef(lang string, title string) string {
	var refBuilder strings.Builder
	refBuilder.WriteString(lang)
	refBuilder.WriteString("/")
	refBuilder.WriteString(title)
	return refBuilder.String()
}

func maxVersion(list []*pb.Version) *pb.Version {
	var res *pb.Version
	if len(list) != 0 {
		res = list[0]
		for _, current := range list {
			if current.Number > res.Number {
				res = current
			}
		}
	}
	return res
}
