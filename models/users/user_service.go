package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"wiki/hash"
)

const (
	UserPwPepper  = "my-secret-pepper-string123!"
	hmacSecretKey = "my-secret-hmac-key"
)

var (
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
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{uv}, nil
}

// Authenticate is used to authenticate a user with the provided
// email address and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	userPassHash := []byte(foundUser.PasswordHash)
	inputPass := []byte(password + UserPwPepper)
	if err = passwordHashesMatch(userPassHash, inputPass); err != nil {
		return nil, err
	}
	return foundUser, nil
}

// passHashesMatch compares the bcrypt hashes
func passwordHashesMatch(userPassHash, inputPass []byte) error {
	err := bcrypt.CompareHashAndPassword(userPassHash, inputPass)
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return ErrInvalidPassword
		default:
			return err
		}
	}
	return nil
}
