package server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
)

const InterfaceKey joiner.InterfaceKey = "server_http"
const PortInterfaceKey joiner.InterfaceKey = "server_http_port"

type Params map[string]string
type WorkerHTTP func(*auth.User, Params, *http.Request) (server.Response, error)

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type Operator interface {
	HandleEndpoint(key, serverPath string, endpoint Endpoint) error
	HandleFiles(key, serverPath string, staticPath StaticPath) error

	Start() error
}

//type Param struct {
//	Name  string
//	Left string
//}
//
//type Params []Param
//
//func (p Params) ByName(name string) string {
//	for i := range p {
//		if p[i].Name == name {
//			return p[i].Left
//		}
//	}
//	return ""
//}
//
//func (p Params) ByNum(num uint) string {
//	if int(num) >= len(p) {
//		return ""
//	}
//
//	return p[num].Left
//}

//func (p Info) AllExcept(names ...string) []string {
//	var values []string
//
//PARAM:
//	for _, param := range p {
//		for _, name := range names {
//			if param.Title == name {
//				continue PARAM
//			}
//			values = append(values, param.Left)
//		}
//	}
//	return values
//}
