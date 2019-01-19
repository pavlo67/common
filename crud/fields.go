package crud

import (
	"strconv"
	"time"

	"github.com/pavlo67/punctum/basis"
	"github.com/pkg/errors"
)

type FieldType string

const TypeString FieldType = "str"
const TypeBytes FieldType = "bytes"
const TypeUInt FieldType = "uint"
const TypeInt FieldType = "int"
const TypeFloat FieldType = "float"
const TypeBool FieldType = "bool"
const TypeTime FieldType = "time"

type StringMap map[string]string

type Field struct {
	Key        string    `json:"key"`
	MaxLength  int       `json:"max_length,omitempty"`
	Type       FieldType `json:"type,omitempty"`
	Format     string    `json:"format,omitempty"`
	Creatable  bool      `json:"creatable,omitempty"`
	Updatable  bool      `json:"editable,omitempty"`
	Additable  bool      `json:"additable,omitempty"`
	Unique     bool      `json:"unique,omitempty"`
	NotEmpty   bool      `json:"not_empty,omitempty"`
	AutoUnique bool      `json:"auto_unique,omitempty"`
}

func FieldsMap(fields []Field) map[string]Field {
	fieldsMap := map[string]Field{}

	for _, f := range fields {
		fieldsMap[f.Key] = f
	}

	return fieldsMap
}

type NativePtrList []*interface{}

func StringMapToNativePtrList(data StringMap, fields []Field) (NativePtrList, error) {
	var nativeList NativePtrList
	for _, f := range fields {
		val, err := StringToNative(data[f.Key], f)
		if err != nil {
			return nil, errors.Wrapf(err, "can't crud.StringToNative('', %#v)", f)
		} else if val == nil {
			return nil, errors.Errorf("can't crud.StringToNative('', %#v): no value", f)
		}
		nativeList = append(nativeList, &val)
	}

	return nativeList, nil
}

func NativePtrListToNativeMap(nativePtrList NativePtrList, fields []Field) (NativeMap, error) {
	native := NativeMap{}
	for i, f := range fields {
		if nativePtrList[i] != nil {
			native[f.Key] = *nativePtrList[i]
		}
	}

	return native, nil
}

type NativeMap map[string]interface{}

func StringMapToNativeMap(data StringMap, fields []Field) (NativeMap, error) {
	var err error
	native := NativeMap{}
	for _, f := range fields {
		if str, ok := data[f.Key]; ok {
			native[f.Key], err = StringToNative(str, f)
			if err != nil {
				return nil, err
			}
		}
	}

	return native, nil
}

func NativeMapToStringMap(nativeMap NativeMap, fields []Field) (StringMap, error) {
	var err error
	data := StringMap{}
	for _, f := range fields {
		if val, ok := nativeMap[f.Key]; ok {
			data[f.Key], err = NativeToString(val, f)
			if err != nil {
				return nil, err
			}
		}
	}
	return data, nil
}

func StringToNative(str string, f Field) (interface{}, error) {
	switch f.Type {
	case TypeString:
		if str == "" {
			return *new(string), nil
		}
		return str, nil
	case TypeBytes:
		if str == "" {
			return *new([]byte), nil
		}
		return str, nil
	case TypeUInt:
		if str == "" {
			return *new(uint64), nil
		}
		return strconv.ParseUint(str, 10, 64)
	case TypeInt:
		if str == "" {
			return *new(int64), nil
		}
		return strconv.ParseInt(str, 10, 64)
	case TypeFloat:
		if str == "" {
			return *new(float64), nil
		}
		return strconv.ParseFloat(str, 64)
	case TypeBool:
		if str == "" {
			return *new(bool), nil
		}
		return strconv.ParseBool(str)
	case TypeTime:
		if str == "" {
			return *new(time.Time), nil
		}
		return time.Parse(time.RFC3339, str)
	}

	return nil, nil
}

func NativeToString(val interface{}, f Field) (string, error) {
	// TODO: check f.Type also

	if val == nil {
		return "", nil
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case time.Time:
		return v.Format(time.RFC3339), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil

	}

	return "", errors.Wrapf(basis.ErrWrongDataType, "on NativeToString: can't convert data %#v", val)
}

func NativeToBool(val interface{}, f Field) (bool, error) {
	if val == nil {
		return false, nil
	}

	switch v := val.(type) {
	case string:
		return strconv.ParseBool(v)
	case []byte:
		return strconv.ParseBool(string(v))
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	}

	return false, errors.Wrapf(basis.ErrWrongDataType, "on NativeToBool: can't convert data %#v (%T)", val, val)
}

func NativeToTime(val interface{}, f Field) (time.Time, error) {
	if val == nil {
		return time.Time{}, nil
	}

	switch v := val.(type) {
	case string:
		return time.Parse(time.RFC3339, v)
	case []byte:
		return time.Parse(time.RFC3339, string(v))
	case time.Time:
		return v, nil
	}

	return time.Time{}, errors.Wrapf(basis.ErrWrongDataType, "on NativeToTime: can't convert data %#v (%T)", val, val)
}
