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

package dbclient

import (
	"errors"
	"strings"

	"github.com/glebarez/sqlite" // driver without cgo
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var errUnknownKind = errors.New("Unknown database type")

func New(kind string, addr string) (*gorm.DB, error) {
	kind = strings.ToLower(kind)
	var dialector gorm.Dialector
	switch kind {
	case "sqlite":
		dialector = sqlite.Open(addr)
	case "postgres":
		dialector = postgres.Open(addr)
	case "mysql":
		dialector = mysql.Open(addr)
	case "sqlserver":
		dialector = sqlserver.Open(addr)
	case "clickhouse":
		dialector = clickhouse.Open(addr)
	default:
		return nil, errUnknownKind
	}

	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	err = db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(kind)))
	return db, err
}

func Paginate(db *gorm.DB, start uint64, end uint64) *gorm.DB {
	return db.Offset(int(start)).Limit(int(end - start))
}

func BuildLikeFilter(filter string) string {
	filter = strings.ReplaceAll(filter, ".*", "%")
	var likeBuilder strings.Builder
	if strings.IndexByte(filter, '%') != 0 {
		likeBuilder.WriteByte('%')
	}
	likeBuilder.WriteString(filter)
	if strings.LastIndexByte(filter, '%') != len(filter)-1 {
		likeBuilder.WriteByte('%')
	}
	return likeBuilder.String()
}
