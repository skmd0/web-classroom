package controllers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"wiki/models"
	"wiki/views"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us *models.UserService) *Users {
	return &Users{
		LoginView: views.NewView("bulma", "users/login"),
		NewView:   views.NewView("bulma", "users/new"),
		us:        us,
	}
}

type Users struct {
	LoginView *views.View
	NewView   *views.View
	us        *models.UserService
}

// Login is used to render the login form where user can login.
//
// GET /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	if err := u.LoginView.Render(w, nil); err != nil {
		panic(err)
	}
}

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	err := u.NewView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

type SignUpForm struct {
	Name           string `schema:"name"`
	Email          string `schema:"email"`
	Password       string `schema:"password"`
	RepeatPassword string `schema:"repeat_password"`
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	user := &models.User{
		Model:    gorm.Model{},
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err := u.us.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDB, err := u.us.ByID(user.ID)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, "Name:", userDB.Name)
	fmt.Fprintln(w, "Email:", userDB.Email)
	fmt.Fprintln(w, "Password:", userDB.PasswordHash)
}
