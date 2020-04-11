package controllers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
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
		LoginView: views.NewView("bulma", "users/login"),
		NewView:   views.NewView("bulma", "users/new"),
		us:        us,
	}
}

type Users struct {
	LoginView *views.View
	NewView   *views.View
	us        users.UserService
}

// Login is used to render the login form where user can login.
//
// GET /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	if err := u.LoginView.Render(w, nil); err != nil {
		panic(err)
	}
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
	form := LoginForm{}
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case users.ErrPasswordInvalid:
			http.Error(w, "Invalid password.", http.StatusForbidden)
		case users.ErrNotFound:
			http.Error(w, "Invalid email address.", http.StatusForbidden)
		default:
			http.Error(w, "Oops something went wrong.", http.StatusInternalServerError)
		}
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookie", http.StatusFound)
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

// CookieTest displays email cookie set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	type Alert struct {
		Level   string
		Message string
	}
	type Data struct {
		Alert Alert
		Yield interface{}
	}
	a := Alert{
		Level:   "is-success",
		Message: "Successfully rendered a dynamic alert!",
	}
	d := Data{
		Alert: a,
		Yield: "Hello world!",
	}
	err := u.NewView.Render(w, d)
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

	user := &users.User{
		Model:    gorm.Model{},
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	_, err := u.us.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookie", http.StatusFound)
}
