package dto

import (
	"exaroton-wa-bot/internal/database/entity"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
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

type WhatsappJID struct {
	User       string
	RawAgent   uint8
	Device     uint16
	Integrator uint16
	Server     string
}

func NewWhatsappJID(jid types.JID) WhatsappJID {
	return WhatsappJID{
		User:       jid.User,
		RawAgent:   jid.RawAgent,
		Device:     jid.Device,
		Integrator: jid.Integrator,
		Server:     jid.Server,
	}
}

func (w *WhatsappJID) To() types.JID {
	return types.JID{
		User:       w.User,
		RawAgent:   w.RawAgent,
		Device:     w.Device,
		Integrator: w.Integrator,
		Server:     w.Server,
	}
}

type WhatsappMessage struct {
	Conversation *string
}

func (w *WhatsappMessage) To() *waE2E.Message {
	return &waE2E.Message{
		Conversation: w.Conversation,
	}
}

type WhatsappSendResponse struct {
	// The message timestamp returned by the server
	Timestamp time.Time

	// MessageID is the internal ID of a WhatsApp message
	ID string

	// The server-specified ID of the sent message. Only present for newsletter messages
	ServerID int

	// The identity the message was sent with (LID or PN)
	// This is currently not reliable in all cases
	Sender WhatsappJID
}

func NewWhatsappSendResponse(resp whatsmeow.SendResponse) WhatsappSendResponse {
	return WhatsappSendResponse{
		Timestamp: resp.Timestamp,
		ID:        resp.ID,
		ServerID:  resp.ServerID,
		Sender:    NewWhatsappJID(resp.Sender),
	}
}

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

type GetWhatsappGroupReq struct {
	Whitelist *bool `query:"whitelist"`
}

type WhitelistWhatsappGroupReq struct {
	User   string `json:"user"`
	Server string `json:"server"`
}

func (r *WhitelistWhatsappGroupReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.User, validation.Required),
		validation.Field(&r.Server, validation.Required),
	)
}

// UnwhitelistWhatsappGroupReq
type UnwhitelistWhatsappGroupReq struct {
	User   string `json:"user"`
	Server string `json:"server"`
}

func (r *UnwhitelistWhatsappGroupReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.User, validation.Required),
		validation.Field(&r.Server, validation.Required),
	)
}

// metadata for whatsapp whitelisted group
// dto for db
type WhatsappWhitelistedGroup struct {
	UserJID   string `db:"jid"`        // user_jid
	ServerJID string `db:"server_jid"` // server_jid
}

func NewWhatsappWhitelistedGroup(e *entity.WhatsappWhitelistedGroup) *WhatsappWhitelistedGroup {
	return &WhatsappWhitelistedGroup{
		UserJID:   e.JID,
		ServerJID: e.ServerJID,
	}
}
