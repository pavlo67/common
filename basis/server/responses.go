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

func ResponseRESTError(status int, err error) (Response, error) {
	if err == nil {
		err = errors.Errorf("unknown error with status %d", status)
	}
	if status == 0 || status == http.StatusOK {
		status = http.StatusInternalServerError
	}

	return Response{Status: status}, err
}

func ResponseRESTOk(data interface{}) (Response, error) {
	if data == nil {
		return Response{Status: http.StatusOK}, nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return Response{Status: http.StatusInternalServerError}, errors.Wrapf(err, "can't marshal data (%#v)", data)
	}

	return Response{Status: http.StatusOK, Data: jsonBytes}, nil
}
