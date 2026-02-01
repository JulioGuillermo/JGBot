package telegramchannel

import (
	"JGBot/channels"
	"JGBot/channels/telegramchannel/telegramdb"
	"JGBot/log"
	"fmt"

	"github.com/go-telegram/bot/models"
)

type TelegramChannel struct {
	Ctl   *TelegramCtl
	onMsg channels.OnMessageHandler
}

func NewTelegramChannel() (*TelegramChannel, error) {
	err := telegramdb.Migrate()
	if err != nil {
		return nil, err
	}

	ch := &TelegramChannel{}
	conf := GetTelegramConf()

	ctl, err := NewTelegramCtl(conf.Token)
	if err != nil {
		return nil, err
	}

	ctl.OnMsg = ch.handler
	ch.Ctl = ctl

	return ch, nil
}

func (ch *TelegramChannel) handler(msg *models.Message) {
	sender, err := telegramdb.ReceivedSender(
		msg.From.ID,
		msg.From.FirstName,
		msg.From.LastName,
		msg.From.Username,
	)
	if err != nil {
		log.Error("Error receiving sender", "error", err)
		return
	}

	var chatName string
	if msg.Chat.Title != "" {
		chatName = msg.Chat.Title
	} else if msg.Chat.FirstName != "" {
		chatName = msg.Chat.FirstName
	} else if msg.Chat.LastName != "" {
		chatName = msg.Chat.LastName
	} else if msg.Chat.Username != "" {
		chatName = msg.Chat.Username
	} else {
		chatName = sender.String()
	}
	chat, err := telegramdb.ReceivedChat(msg.Chat.ID, chatName)
	if err != nil {
		log.Error("Error receiving chat", "error", err)
		return
	}

	message, err := telegramdb.ReceivedMessage(chat, sender, msg.ID, msg.Text)
	if err != nil {
		log.Error("Error receiving message", "error", err)
		return
	}

	if ch.onMsg != nil {
		ch.onMsg(
			ch.GetName(),
			fmt.Sprintf("%s:%d", ch.GetName(), chat.ChatID),
			chat.ID,
			chat.ChatName,
			sender.ID,
			sender.String(),
			message.ID,
			message.Text,
		)
	}
}

func (ch *TelegramChannel) GetName() string {
	return "Telegram"
}

func (ch *TelegramChannel) OnMessage(handler channels.OnMessageHandler) {
	ch.onMsg = handler
}

func (ch *TelegramChannel) SendMessage(chatID uint, message string) error {
	chat, err := telegramdb.GetChat(chatID)
	if err != nil {
		log.Error("Fail to find the chat", "chatID", chatID, "error", err)
		return err
	}

	err = ch.Ctl.SendMessage(chat.ChatID, message)
	if err != nil {
		log.Error("Fail to send message to chat", "chatID", chatID, "error", err)
		return err
	}

	return nil
}

func (ch *TelegramChannel) ReactMessage(chatID uint, messageID uint, reaction string) error {
	chat, err := telegramdb.GetChat(chatID)
	if err != nil {
		log.Error("Fail to find the chat", "chatID", chatID, "error", err)
		return err
	}

	msg, err := telegramdb.GetMessage(messageID)
	if err != nil {
		log.Error("Fail to find the message", "messageID", messageID, "error", err)
		return err
	}

	err = ch.Ctl.ReactMessage(chat.ChatID, msg.MessageID, reaction)
	if err != nil {
		log.Error("Fail to react message", "chatID", chatID, "messageID", messageID, "error", err)
		return err
	}
	return nil
}

func (ch *TelegramChannel) Close() {
	ch.Ctl.Close()
}
