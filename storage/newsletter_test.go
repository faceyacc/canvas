package storage_test

import (
	"context"
	"deeler/integrationtest"
	"testing"

	"github.com/matryer/is"
)

func TestDatabase_SignupForNewsletter(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("signs up", func(t *testing.T) {
		is := is.New(t)

		// Create database
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		// Singup for newsletter with test email
		expectedToken, err := db.SignupForNewsletter(context.Background(), "me@example.com")
		is.NoErr(err)
		is.Equal(64, len(expectedToken))

		// Check that the email and token was inserted into the database
		var email, token string
		err = db.DB.QueryRow(`select email, token from newsletter_subscribers`).Scan(&email, &token)
		is.NoErr(err)
		is.Equal("me@example.com", email)
		is.Equal(expectedToken, token)

		// Test with same email
		expectedToken2, err := db.SignupForNewsletter(context.Background(), "me@example.com")
		is.NoErr(err)
		is.True(expectedToken != expectedToken2)

		err = db.DB.QueryRow(`select email, token from newsletter_subscribers`).Scan(&email, &token)
		is.NoErr(err)
		is.Equal("me@example.com", email)
		is.Equal(expectedToken2, token)
	})
}
