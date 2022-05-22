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

func EmailPassConfirm(lang, fromEmail, fromName, toEmail, toName, url string) (Event, error) {
	const tplName = "email-pass-confirm"

	var ev Event

	type Vars struct {
		Sender       string `json:"sender"`
		SenderName   string `json:"sender_name"`
		Receiver     string `json:"receiver"`
		ReceiverName string `json:"receiver_name"`
		URL          string `json:"url"`
	}

	vars := Vars{fromEmail, fromName, toEmail, toName, url}
	encodedVars, err := json.Marshal(&vars)
	if err != nil {
		return ev, errors.Wrap(err, "encode vars")
	}

	params := SendTemplate{
		Template: tplName,
		Language: lang,
		Vars:     encodedVars,
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
