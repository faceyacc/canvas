package server

import "deeler/handlers"

// setupRoutes registers the routes for the server.
func (s *Server) setupRoutes() {
	// Register routes
	handlers.Health(s.mux)
	handlers.FrontPage(s.mux)
}
