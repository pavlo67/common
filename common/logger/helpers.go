package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Log(l Operator, data string, isError bool, message string) {
	data = message + ": " + data
	if l == nil {
		if isError {
			data = "ERROR: " + data
		}
		log.Print(data)
	} else if isError {
		l.Error(data)
	} else {
		l.Info(data)
	}
}

func LogIntoFile(logfile string, l Operator, data, message string) {

	data += "\n"
	Log(l, data, false, message)

	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			Log(l, err.Error(), true, message)
		}
		defer f.Close()
		if _, err := f.WriteString(data); err != nil {
			Log(l, err.Error(), true, message)
		}
	}
}

func FatalIntoFile(logfile string, l Operator, data, message string) {

	data += "\n"
	Log(l, data, true, message)

	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			Log(l, err.Error(), true, message)
		}
		defer f.Close()
		if _, err := f.WriteString(data); err != nil {
			Log(l, err.Error(), true, message)
		}
	}

	os.Exit(1)
}

const onLogTimestamp = "on logger.LogTimestamp()"

func LogTimestamp(logfile string, l Operator) {
	message := "\n\ncreated at: " + time.Now().Format(time.RFC3339)

	Log(l, message, false, "")

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Log(l, err.Error(), true, onLogTimestamp)
	}
	defer f.Close()
	if _, err := f.WriteString(message); err != nil {
		Log(l, err.Error(), true, onLogTimestamp)
	}
}

// const MaxLoggedDataLength = 2048
const MaxLoggedDataLength = 20000

func LogRequest(logfile string, l Operator, method, path string, requestHeaders http.Header, requestBody []byte, responseHeaders http.Header, responseBody []byte, bodyErr error,
	status int) {
	if len(requestBody) > MaxLoggedDataLength {
		requestBody = requestBody[:MaxLoggedDataLength]
	}
	if len(responseBody) > MaxLoggedDataLength {
		responseBody = responseBody[:MaxLoggedDataLength]
	}

	data := fmt.Sprintf("\n\n%s %s %s in %s\nheaders: %#v\nbody: %s\nresponse: %d %s %#v", method, path, time.Now().Format(time.RFC3339), logfile, requestHeaders, requestBody,
		status,
		responseBody, responseHeaders)
	if bodyErr != nil {
		data = fmt.Sprintf("\nERROR: %s", bodyErr) + data
	}

	LogIntoFile(logfile, l, data, "")
}
