package connector

import (
	"time"

	"github.com/pavlo67/partes/crud"
	"github.com/pkg/errors"
)

var _ crud.Mapper = &MessageMapper{}

type MessageMapper struct {
	*Message
}

func (mm *MessageMapper) Fields() []crud.Field {
	return []crud.Field{
		{Key: "id", NotEmpty: true},
		{Key: "size", NotEmpty: true},
		{Key: "from", Creatable: true},
		{Key: "to", Creatable: true},
		{Key: "cc", Creatable: true},
		{Key: "subject", Creatable: true},
		{Key: "content_encoding", Creatable: true},
		{Key: "body", Creatable: true},
		// {Key: "files", Creatable: true},
	}
}

// TODO: test files

func (mm *MessageMapper) Write(data crud.StringMap) (interface{}, error) {
	message := &Message{
		From:            data["from"],
		To:              data["to"],
		CC:              data["cc"],
		Subject:         data["subject"],
		ContentEncoding: data["content_encoding"],
		Body:            data["body"],
	}

	var err error

	if data["date"] != "" {
		message.Date, err = time.Parse(time.RFC3339, data["date"])
	}

	return message, err
}

func (mm *MessageMapper) Read(mapped interface{}) (crud.StringMap, error) {
	var message Message
	var ok bool

	message, ok = mapped.(Message)
	if !ok {
		messagePtr, ok := mapped.(*Message)
		if !ok {
			return nil, errors.New("no connector.Message in mapped data")
		}
		if messagePtr == nil {
			return nil, errors.New("null pointer to connector.Message in mapped data")
		}
		message = *messagePtr
	}

	return crud.StringMap{
		"from":             message.From,
		"to":               message.To,
		"cc":               message.CC,
		"subject":          message.Subject,
		"content_encoding": message.ContentEncoding,
		"body":             string(message.Body),
		"date":             message.Date.Format(time.RFC3339),
	}, nil
}
