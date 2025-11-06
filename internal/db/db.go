package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func OpenSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite (%s): %w", dsn, err)
	}
	return db, nil
}
