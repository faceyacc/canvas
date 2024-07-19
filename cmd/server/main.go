// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"context"
	"deeler/server"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Used for logging and error reporting.
var release string

func main() {
	os.Exit(start())
}

func start() int {
	logEnv := getStringOrDefault("LOG_ENV", "development") // look for LOG_ENV in the environment, default to "development"
	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("Error setting up the logger: ", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		_ = log.Sync() // flushes buffer and makes sure all logs are written before the program exits
	}()
	host := getStringOrDefault("HOST", "localhost")
	port := getIntOrDefault("PORT", 8080)

	// Create a server instance.
	s := server.New(server.Options{
		Host: host,
		Port: port,
		Log:  log,
	})

	// Return context if singal terminated or singal interrupt is received or.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errorGroup, ctx := errgroup.WithContext(ctx)

	errorGroup.Go(func() error {
		if err := s.Start(); err != nil {
			log.Info("Error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	// Listen for interrupt signal.
	<-ctx.Done()

	errorGroup.Go(func() error {
		if err := s.Stop(); err != nil {
			log.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	// Wait for all goroutines to finish.
	if err := errorGroup.Wait(); err != nil {
		return 1
	}
	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

// Helper functions to get an environment variable or return a default value
func getIntOrDefault(name string, defaultV int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}

	vAsInt, err := strconv.Atoi(v)
	if err != nil {
		return defaultV
	}
	return vAsInt
}

func getStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}
