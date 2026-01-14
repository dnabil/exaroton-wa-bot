package errs

import "errors"

// config errors
var (
	ErrCfgYmlPathNotSet      = errors.New(".yml path is not set")
	ErrWarnJwtDurationNotSet = errors.New("JWT exp duration is not set")
)

// general errors
var (
	ErrForbidden    = errors.New("Forbidden")
	ErrUnauthorized = errors.New("Unauthorized")

	ErrUserAlreadyLoggedIn = errors.New("User is already logged in")
	ErrUserNotLoggedIn     = errors.New("User is not logged in")
	ErrLoginFailed         = errors.New("Wrong credentials")
)

// whatsapp errors
var (
	ErrWANotLoggedIn     = errors.New("whatsapp is not logged in")
	ErrWAAlreadyLoggedIn = errors.New("Whatsapp is already logged in")

	ErrWAQRError             = errors.New("QR Pairing error")
	ErrWAQRTimeout           = errors.New("Timeout, please refresh and try again.")
	ErrWAQRIsPairing         = errors.New("Whatsapp is already pairing, please wait for it to finish.")
	ErrWAQRClientOutdated    = errors.New("Whatsapp client is outdated, please update this app.")
	ErrWAQREnableMultidevice = errors.New("Whatsapp is not enabled for multidevice, please enable it.")
)

// Game server specific errors
var (
	ErrGSIsDown        = errors.New("Game server might be down")
	ErrGSInvalidAPIKey = errors.New("Invalid API key")
)
