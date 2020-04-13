package users

import (
	"github.com/jinzhu/gorm"
	"wiki/models"
)

// User is representation of the user DB table
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null; unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB defines the CRUD methods for the user table
type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) (string, error)
	Update(user *User) error
	Delete(id uint) error
}

// userGorm is the implementation of the UserDB interface
type userGorm struct {
	db *gorm.DB
}

// to make sure userGorm implements everything from UserDB interface
var _ UserDB = &userGorm{}

// ByID looks up the user by the provided ID.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// ByEmail looks up a user with the given email address and returns that users.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// ByRemember looks up a user by remember token and returns the user.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Create will create the provided user and fill in the data
// like the ID, CreatedAt and UpdatedAt fields.
func (ug *userGorm) Create(user *User) (string, error) {
	return user.PasswordHash, ug.db.Create(user).Error
}

// Update will update the provided user with all of the date
// in the provided user object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the user with provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// first executes a query from gorm.DB and writes data to dst by reference.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	return err
}
