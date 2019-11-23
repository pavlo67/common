package common

type Getter interface {
	Get(interface{}) error
}
