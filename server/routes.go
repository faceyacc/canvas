package server

import "deeler/handlers"

func (s *Server) setupRoutes() {
	// Register routes
	handlers.Health(s.mux)
}
