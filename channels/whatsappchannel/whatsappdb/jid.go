package whatsappdb

import (
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type JID struct {
	gorm.Model

	// This is for whatsmeow JID
	User       string
	RawAgent   uint8
	Device     uint16
	Integrator uint16
	Server     string
}

func (j *JID) ToJID() *types.JID {
	return &types.JID{
		User:       j.User,
		RawAgent:   j.RawAgent,
		Device:     j.Device,
		Integrator: j.Integrator,
		Server:     j.Server,
	}
}

func (j *JID) FromJID(jid *types.JID) {
	j.User = jid.User
	j.RawAgent = jid.RawAgent
	j.Device = jid.Device
	j.Integrator = jid.Integrator
	j.Server = jid.Server
}

func (j *JID) String() string {
	return j.ToJID().String()
}
