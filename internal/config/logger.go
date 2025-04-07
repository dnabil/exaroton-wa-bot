package config

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
)

func InitLogger(cfg *Cfg) {
	isJsonLog := cfg.Bool(keyIsJsonLog)

	opt := &slog.HandlerOptions{
		Level: strToSlogLevel(cfg.String(keyLogLevel)),
	}

	if isJsonLog {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opt)))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opt)))
	}
}

// Centalizing err log implementation, use this on every err/panic.
//
// errStack indicates that app panicked. just nil if for normal errs.
func ErrLog(ctx context.Context, args *Args, err error, errStack []byte, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String(KeyLogErr, err.Error()))

	// if panic happens, add id & print trace on debug
	if errStack != nil {
		panicAttr := slog.String(keyLogPanicID, uuid.New().String())
		attrs = append(attrs, panicAttr)

		slog.Default().LogAttrs(ctx, slog.LevelDebug, "PANIC TRACE",
			slog.String(KeyLogErrStack, string(errStack)),
			panicAttr,
		)
	}

	slog.Default().LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR", attrs...)
}

// Converts a string to slog.Level type.
// If the string is not recognized, it returns slog.LevelWarn by default.
func strToSlogLevel(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "ERROR":
		return slog.LevelError
	case "WARN":
		fallthrough
	default:
		return slog.LevelWarn
	}
}
