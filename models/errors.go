package models

import "strings"

const (
	// ErrIdInvalid is returned when an invalid ID is provided to a method like Delete(
	ErrIdInvalid modelError = "models: ID provided was invalid"

	// ErrEmailRequire is returned when there is not email set
	ErrEmailRequired modelError = "models: email address is required"

	// ErrEmailInvalid is returned when invalid email is provided
	ErrEmailInvalid modelError = "models: email address is not valid"

	// ErrEmailTaken is return when update or create is attempted with an email already in use
	ErrEmailTaken modelError = "models: email address already taken"

	// ErrPasswordTooShort is returned when password is shorter than 8 characters
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

	// ErrPasswordRequired is returned when Update() or Create() is called with no password
	ErrPasswordRequired modelError = "models: password is required"

	// ErrRememberTooShort is returned when remember token is too short
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"

	// ErrRememberRequired is returned when a Create() or Update() is attempted
	// without a user remember token hash
	ErrRememberRequired modelError = "models: remember token hash is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}
