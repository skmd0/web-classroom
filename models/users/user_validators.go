package users

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"wiki/hash"
	"wiki/rand"
)

var (
	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete()
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrEmailRequire is returned when there is not email set
	ErrEmailRequired = errors.New("models: email address is required")

	// ErrEmailInvalid is returned when invalid email is provided
	ErrEmailInvalid = errors.New("models: email address is not valid")
)

// validator function type signature
type userValFunc func(*User) error

// runUserValFunc loops through all the validator and returns err if any fail
func runUserValFunc(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// to make sure userValidator implements everything from UserDB interface
var _ UserDB = &userValidator{}

// ByEmail will normalize the email address before calling ByEmail on the userDB
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{Email: email}
	if err := runUserValFunc(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the remember token and then call ByRemember on the UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{Remember: token}
	if err := runUserValFunc(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will create the provided user and fill in the data
// like the ID, CreatedAt and UpdatedAt fields.
func (uv *userValidator) Create(user *User) (string, error) {
	err := runUserValFunc(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat)
	if err != nil {
		return "", err
	}
	return uv.UserDB.Create(user)
}

// Update will hash the remember token if it is provided
// and call the Update method on the UserDB layer
func (uv *userValidator) Update(user *User) error {
	err := runUserValFunc(user,
		uv.normalizeEmail,
		uv.bcryptPassword,
		uv.hmacRemember,
		uv.requireEmail,
		uv.emailFormat)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with provided ID
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFunc(&user, uv.idGreaterThanZero)
	if err != nil {
		return err
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

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

func (uv *userValidator) idGreaterThanZero(user *User) error {
	if user.ID == 0 {
		return ErrInvalidID
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}
