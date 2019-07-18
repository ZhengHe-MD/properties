package properties

import (
	"fmt"
	"reflect"
)

func toPropLineBytes(key, val string) []byte {
	return []byte(fmt.Sprintf("%s=%s\n", key, val))
}

func marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Struct && (rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct) {
		return nil, InvalidMarshalError
	}

	return devalue("", rv)
}

func devalue(key string, v reflect.Value) ([]byte, error) {
	var data []byte
	switch v.Kind() {
	case reflect.Ptr:
		return devalue(key, v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			vf, tf := v.Field(i), v.Type().Field(i)

			kk, _ := parseTag(tf.Tag.Get(tagName))

			if kk == "-" {
				continue
			}

			if key != "" {
				kk = fmt.Sprintf("%s.%s", key, kk)
			}

			d, err := devalue(kk, vf)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
	case reflect.Map:
		for _, kk := range v.MapKeys() {
			vv := v.MapIndex(kk)
			d, err := devalue(fmt.Sprintf("%s.%v", key, kk.Interface()), vv)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			vv := v.Index(i)
			d, err := devalue(fmt.Sprintf("%s[%d]", key, i), vv)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
	case reflect.String: fallthrough
	case reflect.Int: fallthrough
	case reflect.Int8: fallthrough
	case reflect.Int16: fallthrough
	case reflect.Int32: fallthrough
	case reflect.Int64: fallthrough
	case reflect.Float32: fallthrough
	case reflect.Float64: fallthrough
	case reflect.Bool: fallthrough
	case reflect.Uint: fallthrough
	case reflect.Uint8: fallthrough
	case reflect.Uint16: fallthrough
	case reflect.Uint32: fallthrough
	case reflect.Uint64:
		return toPropLineBytes(key, fmt.Sprint(v.Interface())), nil
	}
	return data, nil
}