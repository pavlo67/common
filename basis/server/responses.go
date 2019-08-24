package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Response struct {
	Status   int
	Data     []byte
	MIMEType string
	FileName string
}

func ResponseREST(status int, data interface{}) (Response, error) {
	if data == nil {
		return Response{Status: status}, nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return Response{Status: http.StatusInternalServerError}, errors.Wrapf(err, "can't marshal data (%#v)", data)
	}

	return Response{Status: status, Data: jsonBytes}, nil
}
