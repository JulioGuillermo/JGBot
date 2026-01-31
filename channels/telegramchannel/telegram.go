package telegramchannel

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramCtl struct {
	ctx    context.Context
	client *bot.Bot
	OnMsg  func(*models.Message)
}

func NewTelegramCtl(botToken string) (*TelegramCtl, error) {
	ctl := &TelegramCtl{}

	ctl.ctx = context.Background()

	client, err := bot.New(
		botToken,
		bot.WithDefaultHandler(ctl.handler),
		// bot.WithAllowedUpdates(
		// 	bot.AllowedUpdates{
		// 		"message",
		// 		"channel_post",
		// 		"business_message",
		// 	},
		// ),
	)
	if err != nil {
		return nil, err
	}

	ctl.client = client
	go ctl.client.Start(ctl.ctx)
	return ctl, nil
}

func (ctl *TelegramCtl) SendMessage(chatID int64, message string) error {
	_, err := ctl.client.SendMessage(ctl.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
	return err
}

func (ctl *TelegramCtl) ReactMessage(chatID int64, messageID int, reaction string) error {
	_, err := ctl.client.SetMessageReaction(ctl.ctx, &bot.SetMessageReactionParams{
		ChatID:    chatID,
		MessageID: messageID,
		Reaction: []models.ReactionType{
			{
				Type: models.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &models.ReactionTypeEmoji{
					Type:  models.ReactionTypeTypeEmoji,
					Emoji: reaction,
				},
			},
		},
	})
	return err
}

func (ctl *TelegramCtl) Close() {
	if ctl.client != nil {
		ctl.client.Close(ctl.ctx)
	}
}

func (ctl *TelegramCtl) handler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	if ctl.OnMsg != nil {
		ctl.OnMsg(update.Message)
	}
	// fmt.Println("New message", update.Message.Chat.ID)
	// fmt.Println(update.Message.Chat.Title)
	// fmt.Println(update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.Username)
	// fmt.Println(update.Message.Text)
	// ctl.SendMessage(update.Message.Chat.ID, "Hello "+update.Message.From.FirstName)
	// err := ctl.ReactMessage(update.Message.Chat.ID, update.Message.ID, "👍")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
