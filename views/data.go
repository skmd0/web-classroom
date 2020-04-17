package views

import (
	"log"
	"wiki/models/users"
)

const (
	AlertLvlError   = "is-danger"
	AlertLvlWarning = "is-warning"
	AlertLvlInfo    = "is-info"
	AlertLvlSuccess = "is-success"

	// AlertMsgGeneric is displayed when any unpredictable error is encountered by our backend
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// Alert is data used to render Bulma notifications
type Alert struct {
	Level   string
	Message string
}

type Page struct {
	Title string
	URL   string
}

type Breadcrumbs struct {
	Pages []Page
	// LastPageKey is used to apply active-link css class when rendering template
	LastPageKey string
}

// Data is a top level wrapper of data that views render
type Data struct {
	Alert       *Alert
	Breadcrumbs Breadcrumbs
	User        *users.User
	Yield       interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertSuccess(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlSuccess,
		Message: msg,
	}
}

type PublicError interface {
	error

	// Returns the public error message
	Public() string
}
