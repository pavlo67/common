package inspector

import (
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var reCheckUTF8 = regexp.MustCompile(`(?sm)\p{Cyrillic}+`)

func Win1251ToUTF8(data []byte) ([]byte, error) {

	if reCheckUTF8.MatchString(string(data)) {
		return data, errors.New("it is not win1251 filelib")
	}
	rInUTF8 := transform.NewReader(bytes.NewReader(data), charmap.Windows1251.NewDecoder())
	newData, err := ioutil.ReadAll(rInUTF8)
	if err != nil {
		return nil, err
	}
	return newData, nil
}

var re1251Char = regexp.MustCompile(`(?ism)charset=windows-1251`)

func ChangeCharsetWin1251(data []byte) []byte {
	return []byte(re1251Char.ReplaceAllString(string(data), "charset=utf-8"))
}
