package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anirudhlakkaraju/polychat/config/basicauth"
	"github.com/anirudhlakkaraju/polychat/config/props"
	"github.com/anirudhlakkaraju/polychat/config/secrets"
	"github.com/anirudhlakkaraju/polychat/log"
	"github.com/anirudhlakkaraju/polychat/server"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use default json handler as startup logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	logger := slog.Default()

	// Initialize all secrets into os env vars
	must(logger, secrets.Init())

	// Initialize all properties
	must(logger, props.Init("./config/props"))

	// Set gin mode
	setGinMode()

	// Initialize logger
	must(logger, log.Init())

	// Initialize app basic auth
	must(logger, basicauth.Init(ctx))

	// Initilize custom server
	must(logger, server.Init(ctx))

	// Create and run server
	sc := server.GetServerConfig()
	router := sc.CreateServer(ctx)

	// Run server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- router.Run(":" + sc.Port)
	}()

	logger.Info("app up and running", "port", sc.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("Shutting down server...")
		cancel()
	case err := <-serverErr:
		logger.Error("Server error", "err", err)
		cancel()
	}
}

func must(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error("Error during app startup", slog.Any("error", err))
		panic(err)
	}
}

func setGinMode() {
	env := os.Getenv("configEnvironment")
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
