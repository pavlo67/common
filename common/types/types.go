package types

type Key string

type Type struct {
	Key      Key
	Exemplar interface{}
}

const KeyString Key = "string"

var TypeString = Type{
	Key:      KeyString,
	Exemplar: "",
}

const KeyHRefImage Key = "href_image"
const KeyHRef Key = "href"
