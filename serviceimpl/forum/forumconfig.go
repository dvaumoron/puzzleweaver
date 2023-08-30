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

package forumimpl

import (
	"context"

	dbclient "github.com/dvaumoron/puzzleweaver/client/db"
	"github.com/dvaumoron/puzzleweaver/serviceimpl/forum/model"
	"gorm.io/gorm"
)

type forumConf struct {
	DatabaseKind    string
	DatabaseAddress string
}

type initializedForumConf struct {
	db *gorm.DB
}

func initForumConf(ctx context.Context, conf *forumConf) (initializedForumConf, error) {
	db, err := dbclient.New(conf.DatabaseKind, conf.DatabaseAddress)
	if err == nil {
		err = db.AutoMigrate(&model.Thread{}, &model.Message{})
	}
	return initializedForumConf{db: db}, err
}
