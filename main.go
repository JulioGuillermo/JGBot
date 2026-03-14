package main

import (
	"JGBot/agent"
	"JGBot/channels/infrastructure/channelctl"
	"JGBot/conf"
	"JGBot/database"
	"JGBot/log"
	"JGBot/session"
	"JGBot/skill"
	taskdomain "JGBot/task/domain"
	taskports "JGBot/task/ports"
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

	log.Info("Initializing cron and timer...")
	taskports.InitCronService()
	taskports.InitTimerService()

	log.Info("Loading skills...")
	err = skill.InitSkills()
	if err != nil {
		log.Error("Fail to load skills", "error", err)
		os.Exit(1)
	}

	log.Info("Initializing agent...")
	agent, err := agent.NewAgentsCtl()
	if err != nil {
		log.Error("Fail to initialize agent", "error", err)
		os.Exit(1)
	}

	log.Info("Initializing channels...")
	channelCtl, err := channelctl.NewChannelCtl()
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
	taskports.CronService.OnActivation = func(t *taskdomain.Task, s string) {
		session.OnAutoActivation(
			t.TaskOriginInfo.Origin,
			t.TaskOriginInfo.Channel,
			t.TaskOriginInfo.ChatID,
			t.TaskOriginInfo.ChatName,
			t.TaskOriginInfo.SenderID,
			t.TaskOriginInfo.MessageID,
			t.TaskInfo.Name,
			s,
			t.TaskInfo.Description,
			t.TaskInfo.Message,
		)
	}
	taskports.TimerService.OnActivation = func(t *taskdomain.Task, s string) {
		session.OnAutoActivation(
			t.TaskOriginInfo.Origin,
			t.TaskOriginInfo.Channel,
			t.TaskOriginInfo.ChatID,
			t.TaskOriginInfo.ChatName,
			t.TaskOriginInfo.SenderID,
			t.TaskOriginInfo.MessageID,
			t.TaskInfo.Name,
			s,
			t.TaskInfo.Description,
			t.TaskInfo.Message,
		)
	}

	log.Info("Loading timers and cron jobs...")
	err = taskports.CronService.LoadCrons()
	if err != nil {
		log.Error("Fail to load cron jobs", "error", err)
	}
	err = taskports.TimerService.LoadTimers()
	if err != nil {
		log.Error("Fail to load cron jobs", "error", err)
	}

	log.Info("System ready and running...")

	log.Info("Press Ctrl+C to exit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Closing channels...")
	channelCtl.Close()

	log.Info("System stopped...")
}
