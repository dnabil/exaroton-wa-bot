package config

// implements waLog (go.mau.fi/whatsmeow/util/log) waLog.Logger interface
//
// this uses slog for logging

import (
	"context"
	"fmt"
	"log/slog"

	waLog "go.mau.fi/whatsmeow/util/log"
)

// used in whatsapp_db.go
func initWhatsappLogger(cfg *Cfg) (DBLogger, ClientLogger waLog.Logger) {
	return newWASlogger(
			newSlogger(cfg, keyWADBLogLevel),
			"DB",
		),
		newWASlogger(
			newSlogger(cfg, keyWAClientLogLevel),
			"Client",
		)
}

// ================================
//
//	whatsapp logger wrapper
//
// ================================
type WhatsappSlogger struct {
	module string
	logger *slog.Logger
}

// slog logger for whatsmeow/whatsapp
func newWASlogger(slogger *slog.Logger, module string) waLog.Logger {
	return &WhatsappSlogger{
		module: module,
		logger: slogger,
	}
}

func (l *WhatsappSlogger) outputf(level slog.Level, msg string, args ...any) {
	logAttrs := []any{slog.String("module", "WhatsApp:"+l.module)}
	if len(args) > 0 {
		logAttrs = append(logAttrs, args...)
	}

	l.logger.Log(context.Background(), level, msg, logAttrs...)
}

func (l *WhatsappSlogger) Errorf(msg string, args ...interface{}) {
	l.outputf(slog.LevelError, msg, args...)
}

func (l *WhatsappSlogger) Warnf(msg string, args ...interface{}) {
	l.outputf(slog.LevelWarn, msg, args...)
}

func (l *WhatsappSlogger) Infof(msg string, args ...interface{}) {
	l.outputf(slog.LevelInfo, msg, args...)
}

func (l *WhatsappSlogger) Debugf(msg string, args ...interface{}) {
	l.outputf(slog.LevelDebug, msg, args...)
}

func (l *WhatsappSlogger) Sub(mod string) waLog.Logger {
	return &WhatsappSlogger{
		module: fmt.Sprintf("%s/%s", l.module, mod),
		logger: l.logger,
	}
}
