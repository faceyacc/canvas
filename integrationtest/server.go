package integrationtest

import (
	"deeler/server"
	"net/http"
	"testing"
	"time"
)

// CreateServer tests server on port 8081, returning a cleanup function
// stops the server.
func CreateServer() func() {
	db, cleanupDB := CreateDatabase()
	s := server.New(server.Options{
		Host:     "localhost",
		Port:     8081,
		Database: db,
	})

	// start the server
	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	for {
		_, err := http.Get("http://localhost:8081/")
		if err == nil {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	// stop server
	return func() {
		if err := s.Stop(); err != nil {
			panic(err)
		}
		cleanupDB()
	}
}

// SkipIfShort skips the test if it is running in short mode.
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
}
