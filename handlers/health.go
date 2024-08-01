package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type pinger interface {
	Ping(ctx context.Context) error
}

func Health(mux chi.Router, pinger pinger) {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pinger.Ping(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	})
}
