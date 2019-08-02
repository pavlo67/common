package associatio

import "github.com/pavlo67/constructor/server/server_http"

type InterfaceHTTP interface {
	Interface
	WorkerHTTP(EndpointKey) server_http.WorkerHTTP
	CallHTTP(EndpointKey) Call
}

func JoinHTTP(serverOp server_http.Operator, interfaceHTTP InterfaceHTTP) (entryPoint string, err error) {
	return "", nil
}

func InterfaceHTTPForTest(entryPoint string, interfaceHTTP InterfaceHTTP) (Interface, error) {
	return nil, nil
}