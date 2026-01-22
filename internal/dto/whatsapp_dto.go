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
	JID          types.JID
	Name         string
	NameSetAt    time.Time
	GroupCreated time.Time

	ParticipantCount int

	// Group Parent
	IsParent                      bool
	DefaultMembershipApprovalMode string // request_required
}

func NewWhatsappGroupInfo(g *types.GroupInfo) *WhatsappGroupInfo {
	return &WhatsappGroupInfo{
		JID:          g.JID,
		Name:         g.Name,
		NameSetAt:    g.NameSetAt,
		GroupCreated: g.GroupCreated,

		ParticipantCount: g.ParticipantCount, // TODO: fix this, it's not correct

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

	AllGroups         []*WhatsappGroupInfo
	WhiltelistedGroup []*WhatsappGroupInfo
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
