package users

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"wiki/hash"
	"wiki/rand"
)

var (
	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete()
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

type userValFunc func(*User) error

func runUserValFunc(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// to make sure userValidator implements everything from UserDB interface
var _ UserDB = &userValidator{}

// ByRemember will hash the remember token and then call ByRemember on the UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	rememberHash := uv.hmac.Hash(token)
	return uv.UserDB.ByRemember(rememberHash)
}

// Create will create the provided user and fill in the data
// like the ID, CreatedAt and UpdatedAt fields.
func (uv *userValidator) Create(user *User) (string, error) {
	if err := runUserValFunc(user, uv.bcryptPassword); err != nil {
		return "", err
	}
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return "", err
		}
		user.Remember = token
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return uv.UserDB.Create(user)
}

// Update will hash the remember token if it is provided
// and call the Update method on the UserDB layer
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFunc(user, uv.bcryptPassword); err != nil {
		return err
	}
	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with provided ID
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash the user's password with the predefined pepper, when
// the password field is not the empty string.
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + UserPwPepper)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	user.Password = ""
	return nil
}
