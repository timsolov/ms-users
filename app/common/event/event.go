package event

import "github.com/mailru/easyjson"

// Event describes event
type Event struct {
	Subject string
	Payload easyjson.RawMessage
}
