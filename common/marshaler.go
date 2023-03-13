package common

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

const onSave = "on marshaler.Save()"

func Save(data interface{}, marshaler Marshaler, filename, prefix, indent string) error {
	var dataBytes []byte
	var err error

	if prefix == "" && indent == "" {
		dataBytes, err = marshaler.Marshal(data)
	} else {
		dataBytes, err = marshaler.MarshalIndent(data, prefix, indent)
	}
	if err != nil {
		return errors.Wrap(err, onSave)
	}

	if err = os.WriteFile(filename, dataBytes, 0644); err != nil {
		return errors.Wrap(err, onSave)
	}

	return nil
}

const onLoad = "on marshaler.Load()"

func Load(data interface{}, marshaler Marshaler, filename string) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, onLoad)
	}

	if err = marshaler.Unmarshal(dataBytes, data); err != nil {
		return errors.Wrap(err, onLoad)
	}

	return nil
}

var MarshalerYAML = MarshalerStruct{yaml.Marshal, yaml.Unmarshal, yamlMarshalIndent}
var MarshalerJSON = MarshalerStruct{json.Marshal, json.Unmarshal, json.MarshalIndent}

func yamlMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return yaml.Marshal(v)
}

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
	marshal       Marshal
	unmarshal     Unmarshal
	marshalIndent MarshalIndent
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
