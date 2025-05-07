package repository

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/constants/errs"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type waClient struct {
	client iWhatsmeowClientWrapper

	qrSub     chan whatsmeow.QRChannelItem
	qrSubLock sync.RWMutex
}

// NewWAClient creates a new WhatsApp client (without connecting to an account).
//
// represents a single whatsapp device/account.
//
// defer Disconnect()
func NewWAClient(db *config.WhatsappDB) (*waClient, error) {
	deviceStore, err := db.Container.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, db.ClientLogger)
	return &waClient{
		client: &whatsmeowClientWrapper{client: client},
	}, nil
}

func (w *waClient) IsLoggedIn() bool {
	return w.client.IsLoggedIn()
}

// if isn't logged in, will return a channel (to send qr codes) that is closed/nil automatically.
//
// if already logged in, will return nil and nil error.
func (w *waClient) Login(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
	if w.IsLoggedIn() {
		return nil, errs.ErrWAAlreadyLoggedIn
	}

	// try to connect with saved session
	if w.client.GetLoggedInDeviceJID() != nil {
		err := w.client.Connect()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// check for existing qr session
	qrChan := w.getQRSub()
	if qrChan != nil {
		return *w.getQRSub(), nil
	}

	// create a new qr session
	newQRPub, _ := w.client.GetQRChannel(ctx)
	err := w.client.Connect()
	if err != nil {
		return nil, err
	}

	newQRSub := make(chan whatsmeow.QRChannelItem)
	w.setQRSub(&newQRSub)
	w.publishQR(newQRPub)

	return *w.getQRSub(), nil
}

func (w *waClient) Disconnect() {
	if w.client != nil {
		w.client.Disconnect()
	}
}

// starts a goroutine that publishes QR codes to the subscriber.
func (w *waClient) publishQR(pub <-chan whatsmeow.QRChannelItem) {
	if w.qrSub == nil {
		return
	}

	go func() {
		for {
			qr, ok := <-pub
			if !ok || qr.Event != whatsmeow.QRChannelEventCode {
				w.setQRSub(nil)
				return
			}

			subPtr := w.getQRSub()
			*subPtr <- qr
		}
	}()
}

// closes the qr channel if res is nil.
func (w *waClient) setQRSub(res *chan whatsmeow.QRChannelItem) {
	w.qrSubLock.Lock()
	defer w.qrSubLock.Unlock()

	if res == nil && w.qrSub != nil {
		close(w.qrSub)
		w.qrSub = nil
		return
	}

	w.qrSub = *res
}

func (w *waClient) getQRSub() *chan whatsmeow.QRChannelItem {
	w.qrSubLock.RLock()
	defer w.qrSubLock.RUnlock()

	if w.qrSub == nil {
		return nil
	}
	return &w.qrSub
}

// ================================
//
//	whatsmeow wrapper
//
// ================================
type iWhatsmeowClientWrapper interface {
	IsLoggedIn() bool
	Connect() error
	Disconnect()
	GetQRChannel(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)

	// to check wether the client already logged in
	GetLoggedInDeviceJID() *types.JID
}

var _ iWhatsmeowClientWrapper = &whatsmeowClientWrapper{}

type whatsmeowClientWrapper struct {
	client *whatsmeow.Client
}

func (w *whatsmeowClientWrapper) IsLoggedIn() bool {
	return w.client.IsLoggedIn()
}

func (w *whatsmeowClientWrapper) Connect() error {
	return w.client.Connect()
}

func (w *whatsmeowClientWrapper) Disconnect() {
	w.client.Disconnect()
}

func (w *whatsmeowClientWrapper) GetQRChannel(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
	return w.client.GetQRChannel(ctx)
}

func (w *whatsmeowClientWrapper) GetLoggedInDeviceJID() *types.JID {
	return w.client.Store.ID
}
