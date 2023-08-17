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
	pb "github.com/dvaumoron/puzzleprofileservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
	"github.com/dvaumoron/puzzleweaver/web/common/service"
)

// check matching with interface
var _ service.AdvancedProfileService = &profileImpl{}

type profileImpl struct {
	weaver.Implements[service.AdvancedProfileService]
	userService    weaver.Ref[service.UserService]
	authService    weaver.Ref[service.AuthService]
	groupId        uint64
	defaultPicture []byte
}

func (impl profileImpl) UpdateProfile(ctx context.Context, userId uint64, desc string, info map[string]string) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl profileImpl) UpdatePicture(ctx context.Context, userId uint64, data []byte) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl profileImpl) GetPicture(ctx context.Context, userId uint64) ([]byte, error) {
	// TODO
	return impl.defaultPicture, nil
}

func (impl profileImpl) GetProfiles(ctx context.Context, userIds []uint64) (map[uint64]service.UserProfile, error) {

	// duplicate removal
	userIds = common.MakeSet(userIds).Slice()

	users, err := impl.userService.Get().GetUsers(ctx, userIds)
	if err != nil {
		return nil, err
	}

	list := []*pb.UserProfile{}
	tempProfiles := map[uint64]service.UserProfile{}
	for _, profile := range list {
		userId := profile.UserId
		tempProfiles[userId] = service.UserProfile{User: users[userId], Desc: profile.Desc, Info: profile.Info}
	}

	profiles := map[uint64]service.UserProfile{}
	for userId, user := range users {
		profile, ok := tempProfiles[userId]
		if ok {
			profiles[userId] = profile
		} else {
			// user who doesn't have profile data yet
			profiles[userId] = service.UserProfile{User: user}
		}
	}
	return profiles, err
}

// no right check
func (impl profileImpl) Delete(ctx context.Context, userId uint64) error {

	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}

func (impl profileImpl) ViewRight(ctx context.Context, userId uint64) error {
	return impl.authService.Get().AuthQuery(ctx, userId, impl.groupId, service.ActionAccess)
}
