package users

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// to make sure userValidator implements everything from UserDB interface
var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
}
