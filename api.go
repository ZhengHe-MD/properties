package properties

import "errors"

func Marshal(v interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func Unmarshal(data []byte, v interface{}) error {
	return errors.New("not implemented")
}

func UnmarshalKV(kv map[string]string, v interface{}) error {
	return unmarshalKV(kv, v)
}
