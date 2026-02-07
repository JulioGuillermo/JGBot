package whatsappchannel

import "go.mau.fi/whatsmeow/proto/waE2E"

func GetMsgContent(msg *waE2E.Message) string {
	switch {
	case msg.Conversation != nil:
		return msg.GetConversation()
	case msg.ExtendedTextMessage != nil:
		return msg.ExtendedTextMessage.GetText()
	case msg.ImageMessage != nil:
		return msg.ImageMessage.GetCaption()
	case msg.VideoMessage != nil:
		return msg.VideoMessage.GetCaption()
	}
	return ""
}
