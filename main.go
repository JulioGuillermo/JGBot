package main

import (
	"JGBot/agent"
	"JGBot/channels/channelctl"
	"JGBot/conf"
	"JGBot/database"
	"JGBot/log"
	"JGBot/session"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting system...")
	err := conf.InitConfig()
	if err != nil {
		panic(err)
	}
	log.InitLogger(conf.Conf.LogLevel)

	log.Info("Initializing database...")
	err = database.InitConnection()
	if err != nil {
		log.Error("Fail to initialize database", "error", err)
		os.Exit(1)
	}

	log.Info("Initializing agent...")
	agent, err := agent.NewAgentsCtl()
	if err != nil {
		log.Error("Fail to initialize agent", "error", err)
		os.Exit(1)
	}

	log.Info("Initializing channels...")
	channelCtl, err := channelctl.InitChannelCtl()
	if err != nil {
		log.Error("Fail to initialize channels", "error", err)
		os.Exit(1)
	}

	log.Info("Initializing session ctl...")
	session, err := session.NewSessionCtl(
		channelCtl,
		agent,
	)
	if err != nil || session == nil {
		log.Error("Fail to initialize session", "error", err)
		os.Exit(1)
	}

	// channelCtl.OnMessage(func(channel string, chatID uint, chatName string, senderID uint, senderName string, messageID uint, message string) {
	// 	fmt.Println("\033[36m  ### Received a message >>>\033[0m", channel, chatID, chatName, senderID, senderName, messageID, message)
	// 	channelCtl.SendMessage(channel, chatID, "Hi..."+senderName)
	// 	channelCtl.ReactMessage(channel, chatID, messageID, "👍")
	// })

	log.Info("System ready and running...")

	log.Info("Press Ctrl+C to exit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Closing channels...")
	channelCtl.Close()

	log.Info("System stopped...")
}
