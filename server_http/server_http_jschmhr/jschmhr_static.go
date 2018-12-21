package server_http_jschmhr

import (
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"

	"io"
	"os"
)

var reHTMLExt = regexp.MustCompile(`\.html?$`)

func (s *server_http_jschmhr) HandleFile(serverPath, localPath string, mimeType *string) error {
	l.Info("FILES: "+localPath, "\t-->", serverPath)

	// TODO: check localPath

	fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if mimeType != nil {
			w.Header().Set("Content-Type", *mimeType)
			OpenFile, err := os.Open(localPath + "/" + p.ByName("filepath"))
			defer OpenFile.Close()
			if err != nil {
				l.Error(err)
			} else {
				io.Copy(w, OpenFile)
			}
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	return nil
}

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.URL.Path, nil)
// if err != nil {
//	l.Error("can't read MIMEType for file: ", localPath+"/"+r.URL.Path, err)
// }

func (s *server_http_jschmhr) HandleString(serverPath, str string, mimeType *string) {
	s.handleFunc("GET", serverPath, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		if mimeType != nil {
			// "application/javascript"
			w.Header().Set("Content-Type", *mimeType)
		}
		w.Write([]byte(str))
	})
}

// it should be performed in .HandleFile()
//
//func (s *serverhttp_jschmhr) PrepareDirWithMIME(frontPath, path string) string {
//	files, err := ioutil.ReadDir(path)
//	if err != nil {
//		l.Error("can't ioutil.ReadDir(%s): %s", path, err)
//		return ""
//	}
//
//	var htmlFront string
//
//	for _, f := range files {
//		if f.IsDir() {
//			continue
//		}
//		name := f.Name()
//		ext := strings.ToLower(filepath.Ext(name))
//		if ext != ".js" && ext != ".css" {
//			continue
//		}
//
//		bytes, err := ioutil.ReadFile(path + name)
//		if err != nil {
//			l.Error(errors.Wrapf(err, "can't read JS file: %s", path+name))
//			continue
//		}
//
//		switch ext {
//		case ".js":
//			s.Handle("GET", frontPath+name, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//				w.Header().Set("Content-Type", "application/javascript")
//				w.Write(bytes)
//			})
//			htmlFront += `<script type="text/javascript" src="` + frontPath + name + `" type="text/javascript"></script>` + "\n"
//
//		case ".css":
//			s.Handle("GET", frontPath+name, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//				w.Header().Set("Content-Type", "text/css")
//				w.Write(bytes)
//			})
//			htmlFront += `<link rel="stylesheet" type="text/css" href="` + frontPath + name + `">` + "\n"
//
//		case ".png", ".gif", ".tiff":
//			s.Handle("GET", frontPath+name, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//				w.Header().Set("Content-Type", "image/"+ext[1:])
//				w.Write(bytes)
//			})
//
//		case ".jpg", ".jpeg":
//			s.Handle("GET", frontPath+name, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//				w.Header().Set("Content-Type", "image/jpeg")
//				w.Write(bytes)
//			})
//
//		default:
//			s.Handle("GET", frontPath+name, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//				w.Header().Set("Content-Type", "text/plain")
//				w.Write(bytes)
//			})
//
//		}
//
//	}
//
//	return htmlFront
//}
