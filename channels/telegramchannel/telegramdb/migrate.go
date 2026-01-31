package telegramdb

import (
	"JGBot/database"
)

func Migrate() error {
	err := database.DB.AutoMigrate(&TelegramChat{})
	if err != nil {
		return err
	}

	err = database.DB.AutoMigrate(&TelegramMessage{})
	if err != nil {
		return err
	}

	err = database.DB.AutoMigrate(&TelegramSender{})
	if err != nil {
		return err
	}

	return nil
}
