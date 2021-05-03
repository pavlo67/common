package config

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	// MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var MarshalerYAML = MarshalerStruct{yaml.Marshal, yaml.Unmarshal}
var MarshalerJSON = MarshalerStruct{json.Marshal, json.Unmarshal}

// var MarshalerJSON5 = MarshalerStruct{json5.Marshal, json5.Unmarshal}

//// var ConvertorXML = MarshalerStruct{xml.Marshal, xml.marshalIndent, xml.unmarshal}
//
//func MarshalerCustom(Marshal Marshal, marshalIndent marshalIndent, unmarshal unmarshal) (Marshaler, error) {
//	var errs common.Errors
//	if Marshal == nil {
//		errs = append(errs, errors.New("Marshal method is nil"))
//	}
//	if marshalIndent == nil {
//		errs = append(errs, errors.New("marshalIndent method is nil"))
//	}
//	if unmarshal == nil {
//		errs = append(errs, errors.New("unmarshal method is nil"))
//	}
//	if errs != nil {
//		return nil, errs.Err()
//	}
//
//	return &MarshalerStruct{Marshal, marshalIndent, unmarshal}, nil
//}

// MarshalerStruct ----------------------------------------------------------------------------------------

type Marshal func(v interface{}) ([]byte, error)
type MarshalIndent func(v interface{}, prefix, indent string) ([]byte, error)
type Unmarshal func(data []byte, v interface{}) error

var _ Marshaler = &MarshalerStruct{}

type MarshalerStruct struct {
	marshal   Marshal
	unmarshal Unmarshal
	// marshalIndent MarshalIndent
}

func (cs MarshalerStruct) Marshal(v interface{}) ([]byte, error) {
	return cs.marshal(v)
}

//func (cs MarshalerStruct) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
//	return cs.marshalIndent(v, prefix, indent)
//}

func (cs MarshalerStruct) Unmarshal(data []byte, v interface{}) error {
	return cs.unmarshal(data, v)
}
