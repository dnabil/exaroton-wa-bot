package main

import (
	"context"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/handler"
	"exaroton-wa-bot/internal/handler/wahandler"
	"exaroton-wa-bot/internal/repository"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/sync/errgroup"
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

	// whatsapp client
	waClient, err := repository.NewWAClient(waDb)
	if err != nil {
		slog.Error("failed to create whatsapp bot client", "error", err)
		os.Exit(1)
	}

	repo, err := repository.New(db, waClient)
	if err != nil {
		slog.Error("failed to create repo", config.KeyLogErr, err)
		os.Exit(1)
	}

	service := service.New(cfg, db, repo)
	handler := handler.NewWeb(cfg, service)
	waHandler := wahandler.NewWAHandler(
		cfg,
		repo.WhatsappRepo,
		service.AuthService,
		service.ServerSettingsService,
	)

	port, err := strconv.Atoi(cfg.String(config.KeyPort))
	if err != nil {
		slog.Error(fmt.Sprintf("bad port configuration: %s", cfg.String(config.KeyPort)))
		os.Exit(1)
	}

	g, ctx := errgroup.WithContext(ctx)

	autoWaLogin := cfg.Bool(config.KeyAutoWhatsappLogin)
	if autoWaLogin {
		slog.Info("auto login whatsapp during startup is enabled")
		isWaLoggedIn, err := waClient.LoginWithExistingSession(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("error while logging in whatsapp during startup: %s", err.Error()))
			os.Exit(1)
		}

		slog.Info(fmt.Sprintf("whatsapp login status: %t", isWaLoggedIn))
	}

	g.Go(func() error {
		slog.Info("server started")
		waHandler.Run()
		return handler.RunHTTP(port)
	})

	// graceful shutdown
	shutdown := getGracefulShutdown(handler.Router.Server, db, waDb, repo.WhatsappRepo, waHandler)

	// run server
	srvErrs := make(chan error, 1)
	go func() {
		defer config.Recover(context.TODO(), args)
		srvErrs <- g.Wait()
	}()

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
	waHandler *wahandler.WaHandler,
) func(reason interface{}) {
	return func(reason interface{}) {
		// put services that needs to be gracefully shutdown here...
		slog.Info("Server shutting down:", "reason", reason)

		if waHandler != nil {
			waHandler.Stop()
		}

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
