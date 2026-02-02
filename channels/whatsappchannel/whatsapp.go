package whatsappchannel

import (
	"JGBot/log"
	"JGBot/tools"
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppCtl struct {
	dbLog     waLog.Logger
	ctx       context.Context
	container *sqlstore.Container

	clientLog waLog.Logger
	client    *whatsmeow.Client

	OnMsg func(msg *events.Message)
}

func NewWhatsAppCtl(dbName string) (*WhatsAppCtl, error) {
	ctl := &WhatsAppCtl{}
	if err := ctl.init(dbName); err != nil {
		ctl.Close()
		return nil, err
	}

	if err := ctl.connect(); err != nil {
		ctl.Close()
		return nil, err
	}

	if err := ctl.clientConnect(); err != nil {
		ctl.Close()
		return nil, err
	}

	return ctl, nil
}

func (ctl *WhatsAppCtl) Close() {
	if ctl.client != nil {
		ctl.client.Disconnect()
	}
	if ctl.container != nil {
		ctl.container.Close()
	}
}

func (ctl *WhatsAppCtl) Logout() error {
	err := ctl.client.Logout(ctl.ctx)
	if err != nil {
		return err
	}
	return ctl.clientConnect()
}

func (ctl *WhatsAppCtl) sendMsg(chat types.JID, message *waE2E.Message) error {
	_, err := ctl.client.SendMessage(
		ctl.ctx,
		chat,
		message,
	)
	return err
}

func (ctl *WhatsAppCtl) SendMessage(chat types.JID, message string) error {
	return ctl.sendMsg(
		chat,
		&waE2E.Message{
			Conversation: &message,
		},
	)
}

func (ctl *WhatsAppCtl) ReactMessage(chat types.JID, sender types.JID, messageID string, reaction string) error {
	return ctl.sendMsg(
		chat,
		ctl.client.BuildReaction(
			chat,
			sender,
			messageID,
			reaction,
		),
	)
}

func (ctl *WhatsAppCtl) GetGroupInfo(jit types.JID) (*types.GroupInfo, error) {
	return ctl.client.GetGroupInfo(ctl.ctx, jit)
}

func (ctl *WhatsAppCtl) GetContactInfo(jit types.JID) (types.ContactInfo, error) {
	return ctl.client.Store.Contacts.GetContact(ctl.ctx, jit)
}

func (ctl *WhatsAppCtl) init(dbName string) error {
	err := tools.CreateParentDir(dbName)
	if err != nil {
		return err
	}

	ctl.dbLog = waLog.Stdout("DATABASE", "ERROR", true)
	ctl.ctx = context.Background()
	container, err := sqlstore.New(
		ctl.ctx,
		"sqlite3",
		"file:"+dbName+"?_foreign_keys=on",
		ctl.dbLog,
	)
	if err != nil {
		return err
	}
	ctl.container = container
	return nil
}

func (ctl *WhatsAppCtl) connect() error {
	device, err := ctl.container.GetFirstDevice(ctl.ctx)
	if err != nil {
		return err
	}

	ctl.clientLog = waLog.Stdout("CLIENT", "ERROR", true)
	ctl.client = whatsmeow.NewClient(device, ctl.clientLog)
	ctl.client.AddEventHandler(ctl.handler)

	return nil
}

func (ctl *WhatsAppCtl) clientConnect() error {
	// Already logged in, just connect
	if ctl.client.Store.ID != nil {
		return ctl.client.Connect()
	}

	// No ID stored, new login
	qrChan, err := ctl.client.GetQRChannel(context.Background())
	if err != nil {
		return err
	}

	err = ctl.client.Connect()
	if err != nil {
		panic(err)
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println("QR CODE:", evt.Code)
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		} else {
			log.Info("LOGIN EVENT:", "event", evt.Event)
		}
	}

	return nil
}

func (ctl *WhatsAppCtl) handler(event any) {
	switch v := event.(type) {
	case *events.Message:
		if ctl.OnMsg != nil {
			ctl.OnMsg(v)
		}
	}
}
