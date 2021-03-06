package event

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const SendTemplateSubject = "email.SendTemplate"

type SendTemplate struct {
	Template string `json:"template"` // name of template to use
	Language string `json:"language"` // "en"  | "de"
	Vars     []byte `json:"vars"`     // marshaled json object with variables such as subject, receiver, send, sender_name and others to use in template
}

func EmailSendTemplate(tplName, lang, fromEmail, fromName, toEmail string, vars map[string]any) (Event, error) {
	var ev Event

	vars["sender"] = fromEmail
	vars["sender_name"] = fromName
	vars["receiver"] = toEmail
	vars["receiver_name"] = ""

	marshaledVars, err := json.Marshal(&vars)
	if err != nil {
		return ev, errors.Wrap(err, "encode vars")
	}

	params := SendTemplate{
		Template: tplName,
		Language: lang,
		Vars:     marshaledVars,
	}

	encodedParams, err := json.Marshal(&params)
	if err != nil {
		return ev, errors.Wrap(err, "encode params")
	}

	return Event{
		Subject: SendTemplateSubject,
		Payload: encodedParams,
	}, nil
}
