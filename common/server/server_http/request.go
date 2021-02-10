package server_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server"
)

//const OperatorJWTKey = "_operator"
const bodyLogLimit = 2048

const onRequest = "on server_http.Request()"

const ReAuthOpKey = "re_auth_operator"
const ReAuthSuffix = "re_auth_suffix"

type ResponseBinary struct {
	MIMEType string
	Data     []byte
}

func Request(serverURL string, ep EndpointSettled, requestData, responseData interface{}, options *crud.Options, l logger.Operator) error {
	client := &http.Client{}
	method := ep.Endpoint.Method

	var err error
	//for _, doReAuth := range reAuthTries {

	// start of single try

	var requestBody []byte
	var requestBodyReader io.Reader

	if requestData != nil {
		switch v := requestData.(type) {
		case []byte:
			requestBody = v
		case *[]byte:
			requestBody = *v
		case string:
			requestBody = []byte(v)
		case *string:
			requestBody = []byte(*v)
		default:
			if requestBody, err = json.Marshal(requestData); err != nil {
				return errors.Wrapf(err, onRequest+": can't marshal request responseData (%#v)", requestData)
			}
		}

		// must be checked for nil instead direct write
		// the external for GET requests expected nil body, but nil-requestData after json.Marshal return not empty responseData

		requestBodyReader = bytes.NewBuffer(requestBody)
	}

	req, err := http.NewRequest(method, serverURL, requestBodyReader)
	if err != nil || req == nil {
		logger.LogRequest(l, method, serverURL, nil, requestBody, nil, nil, err, 0)
		return fmt.Errorf("can't create request %s %s, got %#v, %s", method, serverURL, req, err)
	} else if req.Body != nil {
		defer Close(req.Body, client, nil)
	}

	if identity := options.GetIdentity(); identity != nil {
		if jwt := identity.GetCredsStr(auth.CredsJWT); jwt != "" {
			req.Header.Add("Authorization", jwt)
		} else if token := identity.GetCredsStr(auth.CredsToken); token != "" {
			req.Header.Add("Authorization", token)
		}

	}
	var responseBody []byte

	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer Close(resp.Body, client, nil)
	}

	if err != nil {
		var statusCode int
		var responseHeaders http.Header
		if resp != nil {
			statusCode = resp.StatusCode
			responseHeaders = resp.Header
			responseBody, _ = ioutil.ReadAll(resp.Body)
		}

		logger.LogRequest(l, method, serverURL, req.Header, requestBody, responseHeaders, responseBody, err, statusCode)
		return errors.Wrapf(err, "can't %s %s", method, serverURL)
	}

	responseBody, err = ioutil.ReadAll(resp.Body)
	logger.LogRequest(l, method, serverURL, req.Header, requestBody, resp.Header, responseBody, err, resp.StatusCode)
	if err != nil {
		return errors.Wrapf(err, "can't read body from %s %s", method, serverURL)
	}

	if resp.StatusCode == http.StatusUnauthorized { // && doReAuth
		//if identity.Token = reAuthJWT(*identity); identity.Token != "" {
		//	continue
		//}
	}

	if resp.StatusCode != http.StatusOK {
		// TODO!!! be careful writing server_http handlers, http.StatusOK is the only success code accepted here

		if len(responseBody) > bodyLogLimit {
			responseBody = responseBody[:bodyLogLimit]
		}

		var data common.Map
		if err = json.Unmarshal(responseBody, &data); err != nil {
			if len(responseBody) > bodyLogLimit {
				responseBody = responseBody[:bodyLogLimit]
			}
			return errors.Wrapf(err, "can't unmarshal body from %s %s: status = %d, body = %s", method, serverURL, resp.StatusCode, responseBody)
		}

		errCommon := fmt.Sprintf("can't %s %s: status = %d, body = %s", method, serverURL, resp.StatusCode, responseBody)
		if data["error"] != nil {
			data["error"] = errata.CommonError(data["error"], errCommon)
		} else {
			data["error"] = errCommon
		}
		errorKey := errata.Key(data.StringDefault(server.ErrorKey, ""))
		return errata.KeyableError(errorKey, data)
	}

	if responseData != nil {
		switch v := responseData.(type) {
		case *[]byte:
			if v != nil {
				*v = responseBody
			}
		case *string:
			if v != nil {
				*v = string(responseBody)
			}
		case *ResponseBinary:
			v.MIMEType = resp.Header.Get("Content-Type")
			v.Data = responseBody
		default:
			if err = json.Unmarshal(responseBody, responseData); err != nil {
				if len(responseBody) > bodyLogLimit {
					responseBody = responseBody[:bodyLogLimit]
				}
				return errors.Wrapf(err, "can't unmarshal body from %s %s: %s", method, serverURL, responseBody)
			}
		}
	}

	//	break // end of each try means the end of all tries if something other wasn't managed before
	//}

	return nil
}
