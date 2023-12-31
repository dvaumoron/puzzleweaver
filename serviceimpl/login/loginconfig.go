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

package loginimpl

import (
	"github.com/dvaumoron/puzzleloginserver/model"
	dbclient "github.com/dvaumoron/puzzleweaver/client/db"
	"gorm.io/gorm"
)

type loginConf struct {
	DatabaseKind    string
	DatabaseAddress string
}

type initializedLoginConf struct {
	db *gorm.DB
}

func initLoginConf(conf *loginConf) (initializedLoginConf, error) {
	db, err := dbclient.New(conf.DatabaseKind, conf.DatabaseAddress)
	if err == nil {
		err = db.AutoMigrate(&model.User{})
	}
	return initializedLoginConf{db: db}, err
}
