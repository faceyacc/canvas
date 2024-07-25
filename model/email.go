package model

import "regexp"

var emailAddressMatcher = regexp.MustCompile(
	`^` + // Start of string
		`(?P<local>[a-zA-Z0-9.!#$%&'*+/=?^_\x60{|}~-]+)` + // Matches local part of the address
		`@` + // Matches the @ symbol
		`(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)` + // Matches the domain part of the address
		`$`, // End of string
)

type Email string

func (e Email) IsValid() bool {
	return emailAddressMatcher.MatchString(string(e))
}

func (e Email) String() string {
	return string(e)
}
