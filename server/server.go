// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
	log     *zap.Logger
}

type Options struct {
	Host string
	Port int
	Log  *zap.Logger
}

// Takes Options as params and returns a new Server instance.
func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		address: address,
		mux:     mux,
		log:     opts.Log,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Method to spin up the server.
func (s *Server) Start() error {
	s.setupRoutes()

	// fmt.Println("Starting on", s.address)
	s.log.Info("Starting server ", zap.String("address", s.address))

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Method to stop the server.
func (s *Server) Stop() error {
	s.log.Info("Stopping server... ")

	// Stop server gracefully within the timeout period.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}
	return nil
}
