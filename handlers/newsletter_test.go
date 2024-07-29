package handlers_test

import (
	"context"
	"deeler/handlers"
	"deeler/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

type signupMock struct {
	email model.Email
}

func (s *signupMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	s.email = email
	return "", nil
}

func TestNewsletterSignup(t *testing.T) {
	mux := chi.NewMux()
	s := &signupMock{}

	handlers.NewsletterSignup(mux, s)

	t.Run("sings up a valid email", func(t *testing.T) {
		is := is.New(t)

		code, _, _ := makePostRequest(mux, "/newsletter/signup", createFormHeader(), strings.NewReader("email=me%40example.com"))
		is.Equal(http.StatusFound, code)
		is.Equal(model.Email("me@example.com"), s.email)
	})

	t.Run("rejects an invalid email", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(mux, "/newsletter/signup", createFormHeader(), strings.NewReader("email=notemail"))
		is.Equal(http.StatusBadRequest, code)
	})

}

func makePostRequest(handler http.Handler, target string, header http.Header, body io.Reader) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header = header
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	result := res.Result()

	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	return result.StatusCode, result.Header, string(bodyBytes)

}

func createFormHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return header
}
