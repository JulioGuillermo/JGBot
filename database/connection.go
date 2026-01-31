package database

import (
	"JGBot/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitConnection() error {
	if config.Conf.Database == "" {
		config.Conf.Database = "db/database.db"
		err := config.Conf.Save()
		if err != nil {
			return err
		}
	}

	db, err := gorm.Open(
		sqlite.Open(config.Conf.Database),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}
	DB = db
	return nil
}
