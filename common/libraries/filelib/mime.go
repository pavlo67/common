package filelib

import (
	"net/http"
	"os"
	"regexp"

	"github.com/pavlo67/common/common/libraries/strlib"
)

var reTXTFile = regexp.MustCompile(`(?i)\.txt$`)
var reRTFFile = regexp.MustCompile(`(?i)\.rtf$`)
var reWin1251 = regexp.MustCompile(`(?ims)content="text/html;\s*charset=windows-1251"`)

var reMIMETextHTML = regexp.MustCompile(`(?i)text/html`)
var reMIMETextPlain = regexp.MustCompile(`(?i)text/plain`)

func MIME(filename string, fileHeader []byte) (string, error) {
	if reRTFFile.MatchString(filename) {
		return "text/richtext", nil
	}

	var mimeType string
	if reTXTFile.MatchString(filename) {
		mimeType = "text/plain"
	}

	if fileHeader == nil {
		OpenFile, err := os.Open(filename)
		defer OpenFile.Close()
		if err != nil {
			return "", err
		}
		fileHeader = make([]byte, 512)
		_, err = OpenFile.Read(fileHeader)
		if err != nil {
			return "", err
		}
	}
	if len(fileHeader) == 0 {
		return mimeType, nil
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(fileHeader)
		// win1251-кодування ця хрєнь не визначає як слід
	}

	if reMIMETextPlain.MatchString(mimeType) {
		// TODO: remove this kostyle!!!
		if strlib.ReCheckUTF8.MatchString(string(fileHeader)) {
			return "text/plain; charset=utf-8", nil
		}
		return "text/plain; charset=windows-1251", nil
	} else if reMIMETextHTML.MatchString(mimeType) {
		// TODO: remove this kostyle!!!
		if reWin1251.MatchString(string(fileHeader)) {
			return "text/html; charset=windows-1251", nil
		} else if strlib.ReCheckUTF8.MatchString(string(fileHeader)) {
			return "text/html; charset=utf-8", nil
		}
	}

	if mimeType == "application/octet-stream" {
		return "", nil
	}
	return mimeType, nil
}

func ImageMIME(ext string) string {

	var iMIME string

	switch ext {
	case ".jpg", ".jpeg":
		iMIME = "image/jpeg"
	case ".gif":
		iMIME = "image/gif"
	case ".png":
		iMIME = "image/png"
	case ".tiff":
		iMIME = "image/tiff"
	case ".bmp":
		iMIME = "image/bmp"
	}

	return iMIME
}

////We read 512 bytes from the file already so we reset the offset back to 0
//OpenFile.Seek(0, 0)
////Send the headers
//w.Header().Set("Contentus-Disposition", "attachment; filename="+FileStat.Nick())
//io.Copy(w, OpenFile) //'Copy' the file to the client
