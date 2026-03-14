package whatsappdb

import "JGBot/database"

func Migrate() error {
	err := database.DB.AutoMigrate(&WhatsAppChat{})
	if err != nil {
		return err
	}

	err = database.DB.AutoMigrate(&WhatsAppSender{})
	if err != nil {
		return err
	}

	err = database.DB.AutoMigrate(&WhatsAppMessage{})
	if err != nil {
		return err
	}
	return nil
}
