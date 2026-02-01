package database

import (
	"JGBot/conf"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitConnection() error {
	if conf.Conf.Database == "" {
		conf.Conf.Database = "db/database.db"
		err := conf.Conf.Save()
		if err != nil {
			return err
		}
	}

	dbDir := filepath.Dir(conf.Conf.Database)
	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		return err
	}

	db, err := gorm.Open(
		sqlite.Open(conf.Conf.Database),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}
	DB = db
	return nil
}
