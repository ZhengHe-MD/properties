package properties

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func unmarshalKV(kv map[string]string, v interface{}) error {
	p := &props{kv: kv}
	return p.unmarshal(v)
}

type props struct {
	kv map[string]string
}

func propsFromBytes(data []byte) (*props, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))

	var kv = map[string]string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip empty line
		if len(line) == 0 {
			continue
		}
		// skip comment line
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, InvalidPropBytes
		}
		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		kv[k] = v
	}

	return &props{kv}, nil
}

func (p *props) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	// NOTE: must be non-nil pointer to struct
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Type().Kind() != reflect.Struct {
		return InvalidUnmarshalError
	}

	return p.value("", rv)
}

func (p *props) value(key string, v reflect.Value) (err error) {
	switch v.Kind() {
	default:
		err = p.valueBasicType(key, v)
	case reflect.Ptr:
		err = p.value(key, v.Elem())
	case reflect.Struct:
		err = p.valueStruct(key, v)
	case reflect.Map:
		err = p.valueMap(key, v)
	case reflect.Slice:
		err = p.valueSlice(key, v)
	}

	return err
}

func (p *props) valueStruct(key string, v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		vf, tf := v.Field(i), v.Type().Field(i)

		if vf.Kind() == reflect.Ptr {
			vf.Set(reflect.New(tf.Type.Elem()))
		}

		// TODO: support opts
		kk, _ := parseTag(tf.Tag.Get(tagName))

		if kk == "-" {
			continue
		}

		if key != "" {
			kk = fmt.Sprintf("%s.%s", key, kk)
		}

		if err := p.value(kk, vf); err != nil {
			return nil
		}
	}
	return nil
}

// valueBasicType deal with int, float, bool, string
func (p *props) valueBasicType(key string, v reflect.Value) error {
	s, ok := p.get(key)
	// NOTE: if key not found, just skip over it.
	if !ok {
		return nil
	}

	switch v.Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		uiv, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(uiv).Convert(v.Type()))
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		iv, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(iv).Convert(v.Type()))
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fv, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(fv).Convert(v.Type()))
	case reflect.String:
		v.Set(reflect.ValueOf(s))
	case reflect.Bool:
		bv, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(bv))
	default:
		return UnsupportedTypeError
	}

	return nil
}

func (p *props) valueMap(key string, v reflect.Value) (err error) {
	m := reflect.MakeMap(v.Type())
	pp := p.subprops(key)
	for kk := range pp.kv {
		mv := reflect.New(v.Type().Elem())
		mk := strings.Split(kk, ".")[0]
		err = pp.value(mk, mv)
		if err != nil {
			return
		}
		m.SetMapIndex(reflect.ValueOf(mk), mv.Elem())
	}
	v.Set(m)
	return
}

func (p *props) valueSlice(key string, v reflect.Value) (err error) {
	var spp = map[string]*props{}
	var sepp = map[string]*props{}

	i := 0
	for {
		sk := fmt.Sprintf("%s[%d]", key, i)
		if !p.hasKeyPrefix(sk) {
			break
		}

		if pp := p.subprops(sk); !pp.isEmpty() {
			spp[sk] = pp
		}

		if epp := p.exactSubprops(sk); !epp.isEmpty() {
			sepp[sk] = epp
		}

		i += 1
	}

	slice := reflect.MakeSlice(v.Type(), 0, len(spp))

	for ii := 0; ii < len(spp); ii++ {
		sk := fmt.Sprintf("%s[%d]", key, ii)
		pp := spp[sk]

		var ev reflect.Value
		if v.Type().Elem().Kind() == reflect.Ptr {
			ev = reflect.New(v.Type().Elem().Elem())
			if err := pp.value("", ev); err != nil {
				return err
			}
			slice = reflect.Append(slice, ev)
		} else {
			ev = reflect.New(v.Type().Elem())
			if err := pp.value("", ev); err != nil {
				return err
			}
			slice = reflect.Append(slice, ev.Elem())
		}
	}

	for ii := 0; ii < len(sepp); ii++ {
		sk := fmt.Sprintf("%s[%d]", key, ii)
		epp := sepp[sk]

		ev := reflect.New(v.Type().Elem())
		err := epp.value(sk, ev)
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, ev.Elem())
	}

	v.Set(slice)
	return nil
}

func (p *props) subprops(prefix string) *props {
	var kv = map[string]string{}

	for k, v := range p.kv {
		if strings.HasPrefix(k, prefix+".") {
			kv[k[len(prefix)+1:]] = v
		}
	}

	return &props{kv}
}

func (p *props) exactSubprops(name string) *props {
	var kv = map[string]string{}

	for k, v := range p.kv {
		if k == name {
			kv[k] = v
		}
	}

	return &props{kv}
}

func (p *props) isEmpty() bool {
	return len(p.kv) == 0
}

func (p *props) get(k string) (string, bool) {
	v, ok := p.kv[k]
	return v, ok
}

func (p *props) hasKeyPrefix(prefix string) bool {
	for k := range p.kv {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}
