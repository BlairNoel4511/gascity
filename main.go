// Package main is the entry point for the gascity application.
// gascity is a fork of gastownhall/gascity, providing gas price
// monitoring and analysis tooling for EVM-compatible networks.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gascity/gascity/internal/config"
	"github.com/gascity/gascity/internal/server"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"
	// Commit is the git commit hash set at build time.
	Commit = "none"
	// BuildDate is the build timestamp set at build time.
	BuildDate = "unknown"
)

func main() {
	// Default to Info level; use DEBUG=1 env var to enable verbose output
	// during local development without having to change the code each time.
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") == "1" {
		logLevel = slog.LevelDebug
	} else if os.Getenv("ENV") == "production" {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("starting gascity",
		"version", Version,
		"commit", Commit,
		"build_date", BuildDate,
	)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	srv, err := server.New(cfg, logger)
	if err != nil {
		slog.Error("failed to initialise server", "error", err)
		os.Exit(1)
	}

	if err := srv.Run(ctx); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}

	fmt.Println("gascity shut down gracefully")
}
