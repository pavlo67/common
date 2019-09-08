package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"

	"github.com/pkg/errors"
)

func RequestJSON(method, url string, data []byte, headers map[string]string) (Map, error) {
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

	log.Printf("%s", body)

	result := Map{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, errors.Wrapf(err, "can't unmarsal: %s", body)
	}

	return result, nil
}
