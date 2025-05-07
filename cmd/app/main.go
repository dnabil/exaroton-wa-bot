package main

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/handler"
	"exaroton-wa-bot/internal/repository"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"gorm.io/gorm"
)

var args *config.Args

func main() {
	var err error
	args, err = config.ParseFlags()
	if err != nil {
		slog.Error("failed to start app", config.KeyLogErr, err)
		os.Exit(1)
	}

	run()
}

func run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// init config
	cfg, err := config.InitConfig(args)
	if err != nil {
		slog.Error("failed to load config", config.KeyLogErr, err)
		os.Exit(1)
	}

	// init default logger
	config.InitAppLogger(cfg)

	defer config.Recover(ctx, args)

	// app db
	db, err := config.InitDB(ctx, cfg)
	if err != nil {
		slog.Error("failed to init db", config.KeyLogErr, err)
		os.Exit(1)
	}

	// whatsapp db
	waDb, err := config.InitWhatsappDB(cfg)
	if err != nil {
		slog.Error("failed to get sql db", config.KeyLogErr, err)
		os.Exit(1)
	}

	repo, err := repository.New(db, waDb)
	if err != nil {
		slog.Error("failed to create repo", config.KeyLogErr, err)
		os.Exit(1)
	}

	service := service.New(cfg, db, repo)
	handler := handler.NewWeb(cfg, service)

	port, err := strconv.Atoi(cfg.String(config.KeyPort))
	if err != nil {
		slog.Error(fmt.Sprintf("bad port configuration: %s", cfg.String(config.KeyPort)))
		os.Exit(1)
	}

	// run server
	srvErrs := make(chan error, 1)
	go func() {
		defer config.Recover(context.TODO(), args)

		slog.Info("server started")
		srvErrs <- handler.RunHTTP(port)
	}()

	// graceful shutdown
	shutdown := getGracefulShutdown(handler.Router.Server, db, waDb, repo.WhatsappRepo)

	select {
	case err := <-srvErrs:
		shutdown(err)
	case <-ctx.Done():
		shutdown("Termination/interrupt signal received")
	}

	slog.Info("server shutdown")
}

func getGracefulShutdown(
	srv *http.Server,
	gormDB *gorm.DB,
	waDb *config.WhatsappDB,
	whatsappRepo repository.IWhatsappRepo,
) func(reason interface{}) {
	return func(reason interface{}) {
		// put services that needs to be gracefully shutdown here...
		slog.Info("Server shutting down:", "reason", reason)

		// whatsapp client
		if whatsappRepo != nil {
			whatsappRepo.Disconnect()
		}

		// waDb
		slog.Info("Disconnecting whatsapp database...")
		if err := waDb.Container.Close(); err != nil {
			slog.Error("Failed to close whatsapp database", config.KeyLogErr, err)
		} else {
			slog.Info("Whatsapp database disconnected")
		}

		// db
		slog.Info("Disconnecting database...")
		db, err := gormDB.DB()
		if err != nil {
			slog.Error("Can't get database instance", config.KeyLogErr, err)
		} else {
			if err = db.Close(); err != nil {
				slog.Error("Failed to close database", config.KeyLogErr, err)
			} else {
				slog.Info("database disconnected")
			}
		}

		if err := srv.Shutdown(context.TODO()); err != nil {
			slog.Error("Error gracefully shutting down server:", config.KeyLogErr, err)
		}
	}
}
