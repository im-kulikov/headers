package headers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

type CustomUnmarshaler interface {
	UnmarshalHeader([]string) error
}

const tagIdentifier = "header"

var ErrUnknownType = errors.New("unknown type")

func UnmarshalHeaders(v any, h http.Header) error {
	if len(h) == 0 {
		return nil
	}

	r := reflect.Indirect(reflect.ValueOf(v))
	t := reflect.TypeOf(v).Elem()

	for i := 0; i < r.NumField(); i++ {
		f := t.Field(i)

		if key, ok := f.Tag.Lookup(tagIdentifier); !ok {
			continue
		} else if key == "-" {
			continue
		} else if val := h.Values(key); len(val) == 0 {
			continue
		} else if err := setValueFromHeader(f, r.Field(i), val); err != nil {
			return err
		}
	}

	return nil
}

func setValueFromHeader(field reflect.StructField, value reflect.Value, values []string) error {
	if value.CanAddr() {
		if custom, ok := value.Addr().Interface().(CustomUnmarshaler); ok {
			return custom.UnmarshalHeader(values)
		}
	}

	switch field.Type.Kind() {
	default:
		return ErrUnknownType
	case reflect.Bool:
		val, err := strconv.ParseBool(values[0])
		if err != nil {
			return err
		}

		value.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return err
		}

		value.SetFloat(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}

		value.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return err
		}

		value.SetUint(val)
	case reflect.String:
		value.SetString(values[0])
	case reflect.Slice:
		count := len(values)
		slice := reflect.MakeSlice(reflect.TypeOf([]string{}), count, count)
		for i := 0; i < count; i++ {
			sliceItem := slice.Index(i)
			sliceItem.SetString(values[i])
		}
		value.Set(slice)
	}

	return nil
}
