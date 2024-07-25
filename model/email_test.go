package model_test

import (
	"deeler/model"
	"testing"

	"github.com/matryer/is"
)

func TestEmailIsValid(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"me@example.com", true},
		{"dhc..3@example.com", false},
		{"@example.com", false},
		{"me@", false},
		{"@", false},
		{"", false},
	}

	t.Run("valid email address", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.address, func(t *testing.T) {
				is := is.New(t)
				e := model.Email(test.address)
				is.Equal(test.valid, e.IsValid())
			})
		}
	})
}
