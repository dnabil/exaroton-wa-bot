package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mau.fi/whatsmeow/types"
)

// whatsapp qr event
const (
	// is sending a new qr code
	WhatsappQREventCode = "code"
	// pairing success
	WhatsappQREventSuccess = "success"

	// == errors ==
	// errors field must be set if event contains "error"

	// pairing error
	WhatsappQREventError = "error"
	// pairing timeout
	WhatsappQREventTimeout = "error-timeout"
	// client outdated
	WhatsappQREventClientOutdated = "error-client-outdated"
	// multidevice not enabled
	WhatsappQREventMultideviceNotEnabled = "error-multidevice-not-enabled"
	// is already pairing
	WhatsappQREventIsPairing = "error-is-pairing"
)

// whatsapp qr websocket response
type WhatsappQRWSRes struct {
	Code  string `json:"code"`
	Event string `json:"event"`
	Error string `json:"error"` // public error message

	Message string `json:"message"`
}

type WhatsappGroupInfo struct {
	JID          types.JID `json:"jid"`
	JIDUser      string    `json:"jid_user"`
	JIDServer    string    `json:"jid_server"`
	Name         string    `json:"name"`
	NameSetAt    time.Time `json:"name_set_at"`
	GroupCreated time.Time `json:"group_created"`

	ParticipantCount int `json:"participant_count"`

	// Group Parent
	IsParent                      bool   `json:"is_parent"`
	DefaultMembershipApprovalMode string `json:"default_membership_approval_mode"` // request_required
}

func NewWhatsappGroupInfo(g *types.GroupInfo) *WhatsappGroupInfo {
	return &WhatsappGroupInfo{
		JID:          g.JID,
		JIDUser:      g.JID.User,
		JIDServer:    g.JID.Server,
		Name:         g.Name,
		NameSetAt:    g.NameSetAt,
		GroupCreated: g.GroupCreated,

		ParticipantCount: len(g.Participants),

		IsParent:                      g.IsParent,
		DefaultMembershipApprovalMode: g.DefaultMembershipApprovalMode,
	}
}

// Whatsapp login page
type WhatsappLoginPageData struct {
	WSPath   string
	HomePath string
}

type SettingsWhatsappPageData struct {
	PhoneNumber string
}

type WhitelistWhatsappGroupReq struct {
	User   string `json:"user"`
	Server string `json:"server"`
}

type GetWhatsappGroupReq struct {
	Whitelist *bool `query:"whitelist"`
}

func (r *WhitelistWhatsappGroupReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.User, validation.Required),
		validation.Field(&r.Server, validation.Required),
	)
}
