package main

import (
	"JGBot/channels/channelctl"
	"JGBot/config"
	"JGBot/database"
	"JGBot/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting system...")
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	log := logger.GetLogger(config.Conf.LogLevel)

	log.Info("Initializing database...")
	err = database.InitConnection()
	if err != nil {
		log.Error("Fail to initialize database", "error", err)
		os.Exit(1)
	}

	log.Info("Initializing channels...")
	channelCtl, err := channelctl.InitChannelCtl()
	if err != nil {
		log.Error("Fail to initialize channels", "error", err)
		os.Exit(1)
	}

	// whatsapp, err := whatsappchannel.NewWhatsAppCtl("db/whatsmeow.db")
	// if err != nil {
	// 	panic(err)
	// }

	// whatsapp.OnMsg = func(msg *events.Message) {
	// 	fmt.Println("\033[36m  ### Received a message >>>\033[0m", msg.Message.GetConversation())
	// 	fmt.Println(msg.Info.PushName)
	// 	// fmt.Println(msg.Message)
	// 	whatsapp.SendMessage(msg.Info.Chat, "Hi..."+msg.Message.GetConversation())
	// 	whatsapp.ReactMessage(msg.Info.Chat, msg.Info.Sender, msg.Info.ID, "👋")
	// }

	channelCtl.OnMessage(func(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
		fmt.Println("\033[36m  ### Received a message >>>\033[0m", channel, chatID, chatName, senderID, senderName, messageID, message)
		channelCtl.SendMessage(channel, chatID, "Hi..."+senderName)
		channelCtl.ReactMessage(channel, chatID, messageID, "👍")
	})

	log.Info("System ready and running...")

	log.Info("Press Ctrl+C to exit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Closing channels...")
	channelCtl.Close()

	log.Info("System stopped...")
}
