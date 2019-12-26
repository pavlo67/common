package httplib

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func DownloadFile(url, pathToLoad string, fileIndex int, perm os.FileMode) (fileName, fileType string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// TODO!!!
	fileType = "html"
	fileName = pathToLoad + strconv.Itoa(fileIndex) + "." + fileType

	err = ioutil.WriteFile(fileName, body, perm)
	if err != nil {
		return "", "", err
	}

	return fileName, fileType, nil
}
