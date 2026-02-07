package whatsappchannel

import (
	"JGBot/channels"
	"JGBot/channels/whatsappchannel/whatsappdb"
	"JGBot/formatter"
	"JGBot/log"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

type WhatsAppChannel struct {
	Ctl               *WhatsAppCtl
	autoEnableSession bool
	onMsg             channels.OnMessageHandler
}

func NewWhatsAppChannel() (*WhatsAppChannel, error) {
	err := whatsappdb.Migrate()
	if err != nil {
		return nil, err
	}

	conf := GetWhatsappConf()
	ch := &WhatsAppChannel{
		autoEnableSession: conf.AutoEnableSession,
	}

	ctl, err := NewWhatsAppCtl(conf.DBPath)
	if err != nil {
		return nil, err
	}

	ctl.OnMsg = ch.handler
	ch.Ctl = ctl

	return ch, nil
}

func (ch *WhatsAppChannel) handler(msg *events.Message) {
	sender, err := whatsappdb.ReceivedSender(&msg.Info.Sender, msg.Info.PushName)
	if err != nil {
		log.Error("Error receiving sender", "Error", err)
		return
	}

	var chatName string
	if msg.Info.IsGroup {
		group, err := ch.Ctl.GetGroupInfo(msg.Info.Chat)
		if err != nil {
			log.Error("Error getting group info", "Error", err)
			chatName = "Unknown group"
		} else {
			chatName = group.Name
		}
	} else {
		chatName = msg.Info.PushName
	}
	chat, err := whatsappdb.ReceivedChat(&msg.Info.Chat, chatName)
	if err != nil {
		log.Error("Error receiving chat", "Error", err)
		return
	}

	message, err := whatsappdb.ReceivedMessage(chat, sender, msg.Info.ID, msg.Message.GetConversation())
	if err != nil {
		log.Error("Error receiving message", "Error", err)
		return
	}

	if ch.onMsg != nil {
		ch.onMsg(
			ch.GetName(),
			fmt.Sprintf("%s:%s", ch.GetName(), chat.String()),
			chat.ID,
			chat.Name,
			sender.ID,
			sender.Name,
			message.ID,
			message.Text,
		)
	}
}

func (ch *WhatsAppChannel) GetName() string {
	return "WhatsApp"
}

func (ch *WhatsAppChannel) OnMessage(handler channels.OnMessageHandler) {
	ch.onMsg = handler
}

func (ch *WhatsAppChannel) Status(chatID uint, status channels.Status) error {
	chat, err := whatsappdb.GetChat(chatID)
	if err != nil {
		log.Error("Fail to find the chat", "chatID", chatID, "error", err)
		return err
	}

	err = ch.Ctl.Status(*chat.ToJID(), status)
	if err != nil {
		log.Error("Fail to send status to chat", "chatID", chatID, "error", err)
		return err
	}

	return nil
}

func (ch *WhatsAppChannel) SendMessage(chatID uint, message string) error {
	message = formatter.FormatMD2WhatsApp(message)

	chat, err := whatsappdb.GetChat(chatID)
	if err != nil {
		log.Error("Fail to find the chat", "chatID", chatID, "error", err)
		return err
	}

	err = ch.Ctl.SendMessage(*chat.ToJID(), message)
	if err != nil {
		log.Error("Fail to send message to chat", "chatID", chatID, "error", err)
		return err
	}

	return nil
}

func (ch *WhatsAppChannel) ReactMessage(chatID uint, messageID uint, reaction string) error {
	chat, err := whatsappdb.GetChat(chatID)
	if err != nil {
		log.Error("Fail to find the chat", "chatID", chatID, "error", err)
		return err
	}

	msg, err := whatsappdb.GetMessage(messageID)
	if err != nil {
		log.Error("Fail to find the message", "messageID", messageID, "error", err)
		return err
	}

	err = ch.Ctl.ReactMessage(*chat.ToJID(), *msg.Sender.ToJID(), msg.MessageID, reaction)
	if err != nil {
		log.Error("Fail to react message", "chatID", chatID, "messageID", messageID, "error", err)
		return err
	}
	return nil
}

func (ch *WhatsAppChannel) AutoEnableSession() bool {
	return ch.autoEnableSession
}

func (ch *WhatsAppChannel) Close() {
	ch.Ctl.Close()
}
