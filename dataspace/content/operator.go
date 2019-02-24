package content

type Type string

type Item interface {
	Type() Type
	Key() string
	Set(interface{}) error
	Refresh() error
	String() string
}
