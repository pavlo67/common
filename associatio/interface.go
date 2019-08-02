package associatio

type EndpointKey string

type Call func (EndpointKey, interface{}) (interface{}, error)

type Interface interface {
	Name() string
	Endpoints() map[EndpointKey]Endpoint
	Call (EndpointKey, interface{}) (interface{}, error)
}

