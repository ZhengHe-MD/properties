package properties

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	InvalidUnmarshalError = errors.New("invalid unmarshal")
	UnsupportedTypeError  = errors.New("unsupported type")
)

func unmarshalKV(kv map[string]string, v interface{}) error {
	p := &props{kv: kv}
	return p.unmarshal(v)
}

type props struct {
	kv map[string]string
}

func (p *props) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	// NOTE: must be non-nil pointer to struct
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Type().Kind() != reflect.Struct {
		return InvalidUnmarshalError
	}

	return p.value("", rv)
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

func (p *props) value(key string, v reflect.Value) error {
	t := v.Type()
	k := v.Kind()
	switch k {
	case reflect.Ptr:
		return p.value(key, v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv, ft := v.Field(i), t.Field(i)

			if !fv.CanSet() {
				log.Printf("%v cannot be set", fv)
				continue
			}

			// TODO: support opts
			kk, _ := parseTag(ft.Tag.Get(tagName))

			if kk == "-" {
				continue
			}

			if key != "" {
				kk = fmt.Sprintf("%s.%s", key, kk)
			}

			if err := p.value(kk, fv); err != nil {
				return nil
			}
		}
	case reflect.String:
		s, ok := p.get(key)

		// use zero value
		if !ok {
			return nil
		}

		v.Set(reflect.ValueOf(s))
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uint:
		s, ok := p.get(key)
		if !ok {
			return nil
		}

		ival, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil
		}

		v.Set(reflect.ValueOf(ival).Convert(t))
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		s, ok := p.get(key)
		if !ok {
			return nil
		}

		ival, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}

		v.Set(reflect.ValueOf(ival).Convert(t))
	case reflect.Bool:
		s, ok := p.get(key)
		if !ok {
			return nil
		}

		bval, err := strconv.ParseBool(s)
		if err != nil {
			return nil
		}

		v.Set(reflect.ValueOf(bval).Convert(t))
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		s, ok := p.get(key)
		if !ok {
			return nil
		}

		fval, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil
		}

		v.Set(reflect.ValueOf(fval).Convert(t))
	case reflect.Map:
		vt := t.Elem()
		mm := reflect.MakeMap(t)
		pp := p.subprops(key)
		for kk := range pp.kv {
			mv := reflect.New(vt)
			mk := strings.Split(kk, ".")[0]
			err := pp.value(mk, mv)
			if err != nil {
				return err
			}
			mm.SetMapIndex(reflect.ValueOf(mk), mv.Elem())
		}
		v.Set(mm)
	case reflect.Slice:
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

		slice := reflect.MakeSlice(t, 0, len(spp))
		for _, pp := range spp {
			ev := reflect.New(t.Elem())
			err := pp.value("", ev)
			if err != nil {
				return err
			}
			slice = reflect.Append(slice, ev.Elem())
		}

		for sk, epp := range sepp {
			ev := reflect.New(t.Elem())
			err := epp.value(sk, ev)
			if err != nil {
				return err
			}
			slice = reflect.Append(slice, ev.Elem())
		}

		v.Set(slice)
	default:
		return UnsupportedTypeError
	}

	return nil
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
