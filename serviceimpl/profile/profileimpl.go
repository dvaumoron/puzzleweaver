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
	"github.com/dvaumoron/puzzleweaver/remoteservice"
	"github.com/dvaumoron/puzzleweaver/web/common"
)

// check matching with interface
var _ remoteservice.RemoteProfileService = &profileImpl{}

type profileImpl struct {
	weaver.Implements[remoteservice.RemoteProfileService]
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
	return nil, nil
}

func (impl profileImpl) GetProfiles(ctx context.Context, userIds []uint64) (map[uint64]remoteservice.RawUserProfile, error) {
	// TODO
	list := []*pb.UserProfile{}
	profiles := map[uint64]remoteservice.RawUserProfile{}
	for _, profile := range list {
		userId := profile.UserId
		profiles[userId] = remoteservice.RawUserProfile{Desc: profile.Desc, Info: profile.Info}
	}
	return profiles, nil
}

func (impl profileImpl) Delete(ctx context.Context, userId uint64) error {
	success := true
	// TODO
	if !success {
		return common.ErrUpdate
	}
	return nil
}
