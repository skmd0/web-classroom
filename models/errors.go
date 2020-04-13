package models

import (
	"strings"
)

const (
	// USER ERRORS

	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound modelError = "models: resource not found"

	// ErrPasswordInvalid is returned when an invalid password is provided
	ErrPasswordInvalid modelError = "models: provided password is invalid"

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

	// ErrPasswordHashRequired is returned when password hash is missing
	ErrPasswordHashRequired modelError = "models: password hash is required"

	// ErrRememberTooShort is returned when remember token is too short
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"

	// ErrRememberRequired is returned when a Create() or Update() is attempted
	// without a user remember token hash
	ErrRememberRequired privateError = "models: remember token hash is required"

	// ErrIdInvalid is returned when an invalid ID is provided to a method like Delete(
	ErrIdInvalid privateError = "models: ID provided was invalid"

	// POST ERRORS

	// ErrUserIdRequired is returned when userID is 0 or less
	ErrUserIdRequired privateError = "models: post UserID is invalid"

	// ErrTitleRequired is returned when post without title is submitted
	ErrTitleRequired modelError = "models: post title is required"

	// ErrContentRequired is returned when post without content is submitted
	ErrContentRequired modelError = "models: post content is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
