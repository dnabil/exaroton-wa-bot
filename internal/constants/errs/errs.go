package errs

import "errors"

// config errors
var (
	ErrCfgYmlPathNotSet      = errors.New(".yml path is not set")
	ErrWarnJwtDurationNotSet = errors.New("JWT exp duration is not set")
)

// general errors
var (
	ErrForbidden       = errors.New("You do not have permission to perform this action")
	ErrUnauthorized    = errors.New("Unauthorized")
	ErrAlreadyReported = errors.New("Already reported") // 208, unique case

	ErrUserAlreadyLoggedIn = errors.New("User is already logged in")
	ErrUserNotLoggedIn     = errors.New("User is not logged in")
	ErrLoginFailed         = errors.New("Wrong credentials")

	ErrWAGroupNotWhitelisted = errors.New("Whatsapp group is not whitelisted")
	ErrInvalidCommandPrefix  = errors.New("Invalid command prefix")
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
	ErrGSIsDown                = errors.New("Game server might be down")
	ErrGSInvalidAPIKey         = errors.New("Invalid API key")
	ErrGSEmptyAPIKey           = errors.New("API key is empty")
	ErrServerNotFound          = errors.New("Server not found")
	ErrServerIsAlreadyStopping = errors.New("Server is already stopped/stopping")
)

// command error
var (
	ErrCommandNotFound   = errors.New("Command not found")
	ErrCommandMissingArg = errors.New("Missing argument")
	ErrCommandInvalidArg = errors.New("Invalid argument")
)
