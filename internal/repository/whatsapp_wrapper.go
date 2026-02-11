package repository

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/constants/errs"
	"sync"
	"sync/atomic"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// represents a single whatsapp device/account.
type waClient struct {
	client iWhatsmeowClientWrapper

	qrSub     chan whatsmeow.QRChannelItem
	qrSubLock sync.RWMutex

	isSyncComplete atomic.Bool
}

// NewWAClient creates a new WhatsApp client (without connecting to an account).
//
// represents a single whatsapp device/account.
//
// defer Disconnect()
func NewWAClient(db *config.WhatsappDB) (*waClient, error) {
	deviceStore, err := db.Container.GetFirstDevice(context.TODO())
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, db.ClientLogger)
	waClient := &waClient{
		client: &whatsmeowClientWrapper{client: client},
	}

	waClient.client.RegisterEventHandler(getStateEventHandler(waClient))

	return waClient, nil
}

func getStateEventHandler(w *waClient) func(evt interface{}) {
	return func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			w.isSyncComplete.Store(true)
		case *events.OfflineSyncPreview:
			w.isSyncComplete.Store(false)
		case *events.OfflineSyncCompleted:
			w.isSyncComplete.Store(true)
		}
	}
}

func (w *waClient) IsSyncComplete(ctx context.Context) bool {
	return w.isSyncComplete.Load()
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

	// login with existing session
	if ok, err := w.LoginWithExistingSession(ctx); ok {
		return nil, err
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

func (w *waClient) LoginWithExistingSession(ctx context.Context) (bool, error) {
	if !w.IsLoggedIn() && w.client.GetLoggedInDeviceJID() != nil {
		err := w.client.Connect()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (w *waClient) Logout(ctx context.Context) error {
	return w.client.Logout(ctx)
}

func (w *waClient) Disconnect() {
	if w.client != nil {
		w.client.Disconnect()
	}
}

func (w *waClient) GetPhoneNumber() string {
	jid := w.client.GetLoggedInDeviceJID()
	if jid == nil {
		return ""
	}

	return jid.User
}

func (w *waClient) GetSelfLID() (*types.JID, error) {
	lid := w.client.GetLoggedInDeviceLID()

	if lid == nil {
		return nil, errors.New("Not logged in yet")
	}

	return lid, nil
}

func (w *waClient) GetGroups(ctx context.Context) ([]*types.GroupInfo, error) {
	return w.client.GetJoinedGroups(ctx)
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
	Logout(ctx context.Context) error

	// to check wether the client already logged in
	GetLoggedInDeviceJID() *types.JID
	GetLoggedInDeviceLID() *types.JID
	GetUserInfo(context.Context, []types.JID) (map[types.JID]types.UserInfo, error)
	GetJoinedGroups(ctx context.Context) ([]*types.GroupInfo, error)
	RegisterEventHandler(f func(any)) uint32
	UnregisterEventHandler(handlerID uint32) bool
}

var _ iWhatsmeowClientWrapper = &whatsmeowClientWrapper{}

type whatsmeowClientWrapper struct {
	client *whatsmeow.Client
}

func (w *whatsmeowClientWrapper) IsLoggedIn() bool {
	return w.client.IsLoggedIn() && w.GetLoggedInDeviceJID() != nil
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

func (w *whatsmeowClientWrapper) Logout(ctx context.Context) error {
	return w.client.Logout(ctx)
}

func (w *whatsmeowClientWrapper) GetLoggedInDeviceJID() *types.JID {
	return w.client.Store.ID
}

func (w *whatsmeowClientWrapper) GetLoggedInDeviceLID() *types.JID {
	if w.client.Store.LID == types.EmptyJID {
		return nil
	}

	return &w.client.Store.LID
}

func (w *whatsmeowClientWrapper) GetUserInfo(ctx context.Context, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	return w.client.GetUserInfo(ctx, jids)
}

func (w *whatsmeowClientWrapper) GetJoinedGroups(ctx context.Context) ([]*types.GroupInfo, error) {
	return w.client.GetJoinedGroups(ctx)
}

func (w *whatsmeowClientWrapper) RegisterEventHandler(f func(any)) uint32 {
	return w.client.AddEventHandler(f)
}

func (w *whatsmeowClientWrapper) UnregisterEventHandler(handlerID uint32) bool {
	return w.client.RemoveEventHandler(handlerID)
}
