package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete()
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password is provided
	ErrInvalidPassword = errors.New("models: provided password is invalid")
)

// UserService is a set of methods used to work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and password
	// are correct.
	Authenticate(email, password string) (*User, error)

	// UserDB is embedded interface with DB methods
	UserDB
}

type userService struct {
	UserDB
}

// to make sure userValidator implements everything from UserDB interface
var _ UserService = &userService{}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

// Authenticate is used to authenticate a user with the provided
// email address and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+UserPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}
