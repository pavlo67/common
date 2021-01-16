package httplib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/workshop/common/data"

	"github.com/pkg/errors"
)

func RequestJSON(method, url string, data []byte, headers map[string]string) (data.Map, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Printf("%s", body)

	result := data.Map{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, errors.Wrapf(err, "can't unmarsal: %s", body)
	}

	return result, nil
}
