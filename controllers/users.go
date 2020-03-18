package controllers

import (
	"fmt"
	"net/http"
	"wiki/views"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers() *Users {
	return &Users{
		LoginView: views.NewView("bulma", "users/login"),
		NewView:   views.NewView("bulma", "users/new"),
	}
}

type Users struct {
	LoginView *views.View
	NewView   *views.View
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

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, "Email:", form.Email)
	fmt.Fprintln(w, "Password:", form.Password)
}
