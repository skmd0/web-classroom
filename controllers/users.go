package controllers

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"wiki/models"
	"wiki/models/users"
	"wiki/rand"
	"wiki/views"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us users.UserService) *Users {
	return &Users{
		LoginView:   views.NewView("bulma", "users/login"),
		NewView:     views.NewView("bulma", "users/new"),
		ProfileView: views.NewView("bulma", "users/profile"),
		us:          us,
	}
}

type Users struct {
	LoginView   *views.View
	NewView     *views.View
	ProfileView *views.View
	us          users.UserService
}

// Login is used to render the login form where user can login.
//
// GET /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginUser is used to verify the provided email address and password
// and then log the user in if the credentials are correct.
//
// POST /login
func (u *Users) LoginUser(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	user.Password = form.Password
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *users.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
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
	var vd views.Data
	var form SignUpForm
	if err := ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := &users.User{
		Model:    gorm.Model{},
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	_, err := u.us.Create(user)
	if err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// /GET /profile
func (u *Users) Profile(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	user := context.User(r.Context())
	vd.Yield = user
	u.ProfileView.Render(w, r, vd)
}

// /GET /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

}
