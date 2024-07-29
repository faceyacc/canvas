package server

import (
	"context"
	"deeler/handlers"
	"deeler/model"
)

// setupRoutes registers the routes for the server.
func (s *Server) setupRoutes() {
	// Register routes
	handlers.Health(s.mux)
	handlers.FrontPage(s.mux)
	handlers.NewsletterSignup(s.mux, &signupMock{})
	handlers.NewsletterThankYou(s.mux)
}

// signupMock is a mock implementation of the signup interface.
type signupMock struct{}

// SignupForNewsletter is a mock implementation of the signup interface.
func (s signupMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	return "", nil
}
