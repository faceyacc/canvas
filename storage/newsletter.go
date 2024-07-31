package storage

import (
	"context"
	"crypto/rand"
	"deeler/model"
	"fmt"
)

// SignupForNewsletter with the given email. Returns a unique token for the user to confirm their email address.
func (d *Database) SingupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	token, err := createSecretToken()
	if err != nil {
		return "", err
	}

	query := `insert into newsletter_subscribers (email, token)
			values ($1, $2)
			on conflict (email) do update set
				token = excluded.token,
				updated = now()`
	_, err = d.DB.ExecContext(ctx, query, email, token)

	return token, err
}

// Generates a unique token for the user to confirm their email address.
func createSecretToken() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", secret), nil
}
