package entity

type WhatsappWhitelistedGroup struct {
	JID       string `gorm:"column:jid"`
	ServerJID string `gorm:"column:server_jid"`
}
