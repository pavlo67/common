package common

import (
	"encoding/json"
	"errors"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var MarshalerJSON = MarshalerStruct{json.Marshal, json.MarshalIndent, json.Unmarshal}
var MarshalerJSON5 = MarshalerStruct{json5.Marshal, json5.MarshalIndent, json5.Unmarshal}

// var ConvertorXML = MarshalerStruct{xml.marshal, xml.marshalIndent, xml.unmarshal}

func MarshalerCustom(marshal marshal, marshalIndent marshalIndent, unmarshal unmarshal) (Marshaler, error) {
	var errs Errors
	if marshal == nil {
		errs = append(errs, errors.New("marshal method is nil"))
	}
	if marshalIndent == nil {
		errs = append(errs, errors.New("marshalIndent method is nil"))
	}
	if unmarshal == nil {
		errs = append(errs, errors.New("unmarshal method is nil"))
	}
	if errs != nil {
		return nil, errs.Err()
	}

	return &MarshalerStruct{marshal, marshalIndent, unmarshal}, nil
}

// MarshalerStruct ----------------------------------------------------------------------------------------

type marshal func(v interface{}) ([]byte, error)
type marshalIndent func(v interface{}, prefix, indent string) ([]byte, error)
type unmarshal func(data []byte, v interface{}) error

var _ Marshaler = &MarshalerStruct{}

type MarshalerStruct struct {
	marshal
	marshalIndent
	unmarshal
}

func (cs MarshalerStruct) Marshal(v interface{}) ([]byte, error) {
	return cs.marshal(v)
}

func (cs MarshalerStruct) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return cs.marshalIndent(v, prefix, indent)
}

func (cs MarshalerStruct) Unmarshal(data []byte, v interface{}) error {
	return cs.unmarshal(data, v)
}
