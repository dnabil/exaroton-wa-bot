package main

import (
	"context"
	"errors"
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/handler"
	"exaroton-wa-bot/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/pflag"
)

var args *config.Args

func parseFlags() (err error) {
	args = &config.Args{}

	pflag.StringVarP(&args.Env, "env", "e", config.EnvProduction, "App's environment (production || development). Default: production")
	pflag.StringVarP(&args.CfgPath, "cfg", "c", "../config.yml", "Path to the config file (optional) e.g: ../config.yml")
	pflag.Parse()

	args.CfgPath, err = filepath.Abs(args.CfgPath)
	if err != nil {
		return errors.New("invalid config path: " + args.CfgPath)
	}

	return nil
}

func main() {
	err := parseFlags()
	if err != nil {
		slog.Error("failed to start app", config.KeyLogErr, err)
	}

	defer config.Recover(context.TODO(), args)

	run()
}

func run() {
	// init config
	cfg, err := config.InitConfig(args)
	if err != nil {
		slog.Error("failed to load config", config.KeyLogErr, err)
		os.Exit(1)
	}

	config.InitLogger(cfg)

	app := handler.NewWeb(cfg, service.New())

	// run server
	srvErrs := make(chan error, 1)
	go func() {
		defer config.Recover(context.TODO(), args)

		slog.Info("server started")
		srvErrs <- app.Run()
	}()

	// gracefulShutdown with its services
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	shutdown := gracefulShutdown(app.Router.Server)

	select {
	case err := <-srvErrs:
		shutdown(err)
	case sig := <-quit:
		shutdown(sig)
	}

	slog.Info("server shutdown")
}

func gracefulShutdown(srv *http.Server) func(reason interface{}) {
	return func(reason interface{}) {
		// put services that needs to be gracefully shutdown here...

		slog.Info("server shutting down:", "reason", reason)

		if err := srv.Shutdown(context.TODO()); err != nil {
			slog.Error("Error gracefully shutting down server:", config.KeyLogErr, err)
		}
	}
}
