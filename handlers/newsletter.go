package handlers

import (
	"context"
	"deeler/model"
	"deeler/views"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type signup interface {
	SignupForNewsletter(ctx context.Context, email model.Email) (string, error)
}

// Handler to have user signup for the newsletter.
func NewsletterSignup(mux chi.Router, s signup) {
	mux.Post("/newsletter/signup", func(w http.ResponseWriter, r *http.Request) {
		email := model.Email(r.FormValue("email"))
		if !email.IsValid() {
			http.Error(w, "Email not valid", http.StatusBadRequest)
			return
		}

		if _, err := s.SignupForNewsletter(r.Context(), email); err != nil {
			http.Error(w, "Dman, we hit a snag, refresh page to try again", http.StatusBadGateway)
			return
		}

		http.Redirect(w, r, "/newsletter/thankyou", http.StatusFound)
	})
}

// Handler to serve Thank You page
func NewsletterThankYou(mux chi.Router) {
	mux.Get("/newsletter/thankyou", func(w http.ResponseWriter, r *http.Request) {
		_ = views.NewsletterThankYouPage("/newsletter/thanks").Render(w)
	})
}
