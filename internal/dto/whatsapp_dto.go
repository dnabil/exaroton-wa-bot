package dto

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

// Whatsapp login page
type WhatsappLoginPageData struct {
	WSPath   string
	HomePath string
}

type SettingsWhatsappPageData struct {
	PhoneNumber string
}
