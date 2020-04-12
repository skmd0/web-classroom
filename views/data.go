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
