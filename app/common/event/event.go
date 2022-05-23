package event

import "github.com/mailru/easyjson"

// Event describes event
type Event struct {
	Subject string
	Payload easyjson.RawMessage
}

// List represents list of events
type List []Event

// Indicate there is no events
var None List
