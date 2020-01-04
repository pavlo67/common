package types

import "github.com/pavlo67/workshop/common/identity"

type Type struct {
	Key      identity.Key
	Exemplar interface{}
}

const KeyString identity.Key = "/string"

var TypeString = Type{
	Key:      KeyString,
	Exemplar: "",
}

const KeyHRefImage identity.Key = "/href_image"
const KeyHRef identity.Key = "/href"
