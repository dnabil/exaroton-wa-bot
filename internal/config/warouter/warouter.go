package warouter

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/dto"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

type Context struct {
	context.Context
	iContext
	Message string
	Args    []string // Message[1:]

	PhoneNumber string // self
	Sender      dto.WhatsappJID
	Chat        dto.WhatsappJID
}

type iContext interface {
	SendMessage(ctx context.Context, to dto.WhatsappJID, message *dto.WhatsappMessage) (*dto.WhatsappSendResponse, error)
}

type HandlerFunc func(c *Context) error

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type WhatsappService interface {
	iContext
	RegisterEventHandler(f func(any)) uint32
	UnregisterEventHandler(handlerID uint32) bool
	GetPhoneNumber() string // self phone number
	GetSelfLID() *dto.WhatsappJID
	IsSyncComplete(ctx context.Context) bool
}

type Router struct {
	waSvc       WhatsappService
	cfg         *config.Cfg
	handlers    map[string]HandlerFunc
	middlewares []MiddlewareFunc // global middlewares

	ErrorHandlerFunc func(c *Context, err error) // nil if not set

	// event handler codes
	HandlerCodeCommandWA uint32
}

func NewRouter(cfg *config.Cfg, waService WhatsappService) *Router {
	return &Router{
		waSvc:    waService,
		cfg:      cfg,
		handlers: make(map[string]HandlerFunc),
	}
}

// Register a new event handler, returns handler code
func (r *Router) registerEventHandler(f func(any)) uint32 {
	return r.waSvc.RegisterEventHandler(f)
}

// Unregister an event handler by its code
func (r *Router) unregisterEventHandler(handlerID uint32) bool {
	return r.waSvc.UnregisterEventHandler(handlerID)
}

// Register a new handler with optional middlewares
func (r *Router) Register(cmd string, h HandlerFunc, mws ...MiddlewareFunc) {
	all := append(r.middlewares, mws...)

	for i := len(all) - 1; i >= 0; i-- {
		h = all[i](h)
	}

	r.handlers[cmd] = h
}

// Run registers the entry point function as an event handler and starts the event loop
func (r *Router) Run() {
	if r.HandlerCodeCommandWA == 0 {
		r.HandlerCodeCommandWA = r.registerEventHandler(r.entryPoint)
	}
}

// entryPoint is an entry point for any event from WhatsApp.
// Currently, it only handles *events.Message. It will create a new Context
// and call the handle function. If the handle function returns an error, it
// will call the ErrorHandlerFunc if it is not nil.
func (r *Router) entryPoint(evt any) {
	// skip all event if sync is not complete
	if !r.waSvc.IsSyncComplete(context.TODO()) {
		return
	}

	switch v := evt.(type) {
	case *events.Message:
		// is it really needed? skip messages from self
		// if v.Info.IsFromMe { return }

		var msg string
		var args []string

		switch v.Message != nil {
		case v.Message.Conversation != nil &&
			*v.Message.Conversation != "":
			msg = *v.Message.Conversation

		case v.Message.ExtendedTextMessage != nil &&
			v.Message.ExtendedTextMessage.Text != nil &&
			*v.Message.ExtendedTextMessage.Text != "":
			msg = *v.Message.ExtendedTextMessage.Text
		}

		parts := strings.Fields(msg)
		if len(parts) > 2 {
			args = parts[2:]
		}

		ctx := &Context{
			Context:     context.Background(),
			iContext:    r.waSvc,
			Message:     msg,
			Args:        args,
			PhoneNumber: r.waSvc.GetPhoneNumber(),
			Sender:      dto.NewWhatsappJID(v.Info.Sender),
			Chat:        dto.NewWhatsappJID(v.Info.Chat),
		}

		err := r.handleMsgEvent(ctx, msg)
		if err != nil && r.ErrorHandlerFunc != nil {
			r.ErrorHandlerFunc(ctx, err)
		}
	}
}

// handleMsgEvent handles incoming message events from WhatsApp.
// It will check if the message contains the self phone number and
// if it starts with a "/". If the checks pass, it will call
// the handle function associated with the command. If the handle
// function returns an error, it will call the ErrorHandlerFunc if it
// is not nil. If the checks fail, it will return nil without
// calling the ErrorHandlerFunc. If the handle function is not found,
// it will reply with "Unknown command".
func (r *Router) handleMsgEvent(c *Context, message string) error {
	parts := strings.Fields(message)

	if len(parts) < 2 {
		return nil
	}

	selfTagged, err := r.isSelfTagged(parts[0])
	if err != nil {
		return err
	}

	if !selfTagged {
		return nil
	}

	// command must start with a "/"
	cmd := parts[1]
	if !strings.HasPrefix(cmd, "/") {
		return nil
	}

	if h, ok := r.handlers[cmd]; ok {
		err := h(c)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (r *Router) isSelfTagged(tag string) (bool, error) {
	lid := r.waSvc.GetSelfLID()
	if lid != nil && lid.User != "" && strings.Contains(tag, lid.User) {
		return true, nil
	}

	phoneNumber := r.waSvc.GetPhoneNumber()
	if phoneNumber == "" {
		return false, errors.New("failed to get whatsapp self id")
	}

	return strings.Contains(tag, phoneNumber), nil
}

// Stop unregisters the entry point function as an event handler, stopping the event loop
func (r *Router) Stop() {
	if r.HandlerCodeCommandWA == 0 {
		_ = r.unregisterEventHandler(r.HandlerCodeCommandWA)
	}
}

// Middleware registration
func (r *Router) Use(mw MiddlewareFunc) {
	r.middlewares = append(r.middlewares, mw)
}
