package server_test

import (
	"deeler/integrationtest"
	"net/http"
	"testing"

	"github.com/matryer/is"
)

func TestServer_Start(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("starts the server and listens for requests", func(t *testing.T) {
		is := is.New(t)

		cleanup := integrationtest.CreateServer()
		defer cleanup()

		// Ping the server to check if it is running.
		resp, err := http.Get("http://localhost:8081/")
		is.NoErr(err)
		is.Equal(http.StatusNotFound, resp.StatusCode)
	})
}
