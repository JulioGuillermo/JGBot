package domain

type MessageHandler func(channel string, origin string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string)
