package views

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

// Data is a top level wrapper of data that views render
type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
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

type PublicError interface {
	error

	// Returns the public error message
	Public() string
}
