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

package wikiclient

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
	wikiservice "github.com/dvaumoron/puzzleweaver/web/wiki/service"
	"golang.org/x/exp/slog"
)

// check matching with interface
var _ wikiservice.WikiService = wikiServiceWrapper{}

type wikiServiceWrapper struct {
	wikiService    remoteservice.RemoteWikiService
	authService    service.AuthService
	profileService service.ProfileService
	loggerGetter   common.LoggerGetter
	wikiId         uint64
	groupId        uint64
	dateFormat     string
	cache          *wikiCache
}

func MakeWikiServiceWrapper(wikiService remoteservice.RemoteWikiService, authService service.AuthService, profileService service.ProfileService, loggerGetter common.LoggerGetter, wikiId uint64, groupId uint64, dateFormat string) wikiservice.WikiService {
	return wikiServiceWrapper{
		wikiService: wikiService, authService: authService, profileService: profileService, loggerGetter: loggerGetter,
		wikiId: wikiId, groupId: groupId, dateFormat: dateFormat, cache: newCache(),
	}
}

func (client wikiServiceWrapper) LoadContent(ctx context.Context, userId uint64, lang string, title string, versionStr string) (*wikiservice.WikiContent, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return nil, err
	}

	logger := client.loggerGetter.Logger(ctx)

	var version uint64
	if versionStr != "" {
		version, err = strconv.ParseUint(versionStr, 10, 64)
		if err != nil {
			logger.Info("Failed to parse wiki version, falling to last", common.ErrorKey, err)
		}
	}
	wikiRef := buildRef(lang, title)
	if version != 0 {
		return client.innerLoadContent(ctx, logger, wikiRef, version)
	}

	list, err := client.wikiService.GetVersions(ctx, client.wikiId, wikiRef)
	if err != nil {
		return nil, err
	}
	if len(list) != 0 {
		content := client.cache.load(client.loggerGetter.Logger(ctx), wikiRef)
		if content != nil && maxVersion(list) == content.Version {
			return content, nil
		}
	}
	return client.innerLoadContent(ctx, logger, wikiRef, 0)
}

func (client wikiServiceWrapper) StoreContent(ctx context.Context, userId uint64, lang string, title string, last string, markdown string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionCreate)
	if err != nil {
		return err
	}

	version, err := strconv.ParseUint(last, 10, 64)
	if err != nil {
		client.loggerGetter.Logger(ctx).Warn("Failed to parse wiki last version", common.ErrorKey, err)
		return common.ErrTechnical
	}

	wikiRef := buildRef(lang, title)
	err = client.wikiService.Store(ctx, client.wikiId, userId, wikiRef, version, markdown)
	if err != nil {
		return err
	}
	client.cache.store(client.loggerGetter.Logger(ctx), wikiRef, &wikiservice.WikiContent{
		Version: version, Markdown: markdown,
	})
	return nil
}

func (client wikiServiceWrapper) GetVersions(ctx context.Context, userId uint64, lang string, title string) ([]wikiservice.Version, error) {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
	if err != nil {
		return nil, err
	}

	list, err := client.wikiService.GetVersions(ctx, client.wikiId, buildRef(lang, title))
	if err != nil {
		return nil, err
	}

	size := len(list)
	if size == 0 {
		return nil, nil
	}

	valueSet := make([]*remoteservice.RawWikiContent, maxVersion(list)+1)
	// no duplicate check, there is one in GetProfiles
	userIds := make([]uint64, 0, size)
	for _, value := range list {
		valueSet[value.Version] = &value
		userIds = append(userIds, value.CreatorId)
	}

	profiles, err := client.profileService.GetProfiles(ctx, userIds)
	if err != nil {
		return nil, err
	}

	newList := make([]wikiservice.Version, 0, size)
	for _, value := range valueSet {
		if value != nil {
			createdAt := time.Unix(value.CreatedAt, 0)
			newList = append(newList, wikiservice.Version{
				Number: value.Version, Creator: profiles[value.CreatorId], Date: createdAt.Format(client.dateFormat),
			})
		}
	}
	return newList, nil
}

func (client wikiServiceWrapper) DeleteContent(ctx context.Context, userId uint64, lang string, title string, versionStr string) error {
	err := client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionDelete)
	if err != nil {
		return err
	}

	logger := client.loggerGetter.Logger(ctx)
	version, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		logger.Warn("Failed to parse wiki version to delete", common.ErrorKey, err)
		return common.ErrTechnical
	}

	wikiRef := buildRef(lang, title)
	content := client.cache.load(logger, wikiRef)
	if content != nil && version == content.Version {
		client.cache.delete(logger, wikiRef)
	}
	return nil
}

func (impl wikiServiceWrapper) DeleteRight(ctx context.Context, userId uint64) bool {
	return impl.authService.AuthQuery(ctx, userId, impl.groupId, service.ActionDelete) == nil
}

func (client wikiServiceWrapper) innerLoadContent(ctx context.Context, logger *slog.Logger, wikiRef string, askedVersion uint64) (*wikiservice.WikiContent, error) {
	res, err := client.wikiService.Load(ctx, client.wikiId, wikiRef, askedVersion)
	if err != nil || res.Version == 0 { // no stored wiki page
		return nil, err
	}

	content := &wikiservice.WikiContent{Version: res.Version, Markdown: res.Markdown}
	if askedVersion == 0 {
		client.cache.store(logger, wikiRef, content)
	}
	return content, nil
}

func buildRef(lang string, title string) string {
	var refBuilder strings.Builder
	refBuilder.WriteString(lang)
	refBuilder.WriteByte('/')
	refBuilder.WriteString(title)
	return refBuilder.String()
}

func maxVersion(list []remoteservice.RawWikiContent) uint64 {
	res := list[0].Version
	for _, current := range list[1:] {
		if current.Version > res {
			res = current.Version
		}
	}
	return res
}
