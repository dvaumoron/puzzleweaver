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

package profileclient

import (
	"context"

	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

type profileServiceWrapper struct {
	profileService remoteservice.RemoteProfileService
	userService    service.UserService
	authService    service.AuthService
	groupId        uint64
	defaultPicture []byte
}

func MakeProfileServiceWrapper(profileService remoteservice.RemoteProfileService, userService service.UserService, authService service.AuthService, groupId uint64) service.ProfileService {
	return profileServiceWrapper{
		profileService: profileService, userService: userService, authService: authService, groupId: groupId,
	}
}

func (client profileServiceWrapper) UpdateProfile(ctx context.Context, userId uint64, desc string, info map[string]string) error {
	return client.UpdateProfile(ctx, userId, desc, info)
}

func (client profileServiceWrapper) UpdatePicture(ctx context.Context, userId uint64, data []byte) error {
	return client.UpdatePicture(ctx, userId, data)
}

func (client profileServiceWrapper) GetPicture(ctx context.Context, userId uint64) []byte {
	// TODO
	return client.defaultPicture
}

func (client profileServiceWrapper) GetProfiles(ctx context.Context, userIds []uint64) (map[uint64]service.UserProfile, error) {
	// duplicate removal
	userIds = common.MakeSet(userIds).Slice()

	users, err := client.userService.GetUsers(ctx, userIds)
	if err != nil {
		return nil, err
	}

	idToRaw, err := client.profileService.GetProfiles(ctx, userIds)
	profiles := map[uint64]service.UserProfile{}
	for id, raw := range idToRaw {
		profiles[id] = service.UserProfile{User: users[id], Desc: raw.Desc, Info: raw.Info}
	}

	// add users who doesn't have profile data yet
	for userId, user := range users {
		if _, ok := profiles[userId]; !ok {
			profiles[userId] = service.UserProfile{User: user}
		}
	}
	return profiles, nil
}

// no right check
func (client profileServiceWrapper) Delete(ctx context.Context, userId uint64) error {
	return client.profileService.Delete(ctx, userId)
}

func (client profileServiceWrapper) ViewRight(ctx context.Context, userId uint64) error {
	return client.authService.AuthQuery(ctx, userId, client.groupId, service.ActionAccess)
}
