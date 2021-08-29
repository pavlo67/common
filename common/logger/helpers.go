package logger

import (
	"fmt"
	"net/http"
	"time"
)

func LogTimestamp(l Operator) {
	l.Info("\n\ncreated at: " + time.Now().Format(time.RFC3339))
}

const MaxLoggedDataLength = 2048

func LogRequest(l Operator, method, path string, reqHeaders http.Header, reqBody []byte, respHeaders http.Header, respBody []byte, bodyErr error, status int) {
	if len(reqBody) > MaxLoggedDataLength {
		reqBody = reqBody[:MaxLoggedDataLength]
	}
	if len(respBody) > MaxLoggedDataLength {
		respBody = respBody[:MaxLoggedDataLength]
	}

	data := fmt.Sprintf("\n%s %s %s\nheaders: %#v\nbody: %s\nresponse: %d %s\nheaders: %#v\n", method, path, time.Now().Format(time.RFC3339), reqHeaders, reqBody, status,
		respBody, respHeaders)
	if bodyErr != nil {
		data = fmt.Sprintf("\nERROR: %s", bodyErr) + data
		l.Error(data)
	} else if status != http.StatusOK {
		l.Error(data)
	} else {
		l.Info(data)
	}
}
