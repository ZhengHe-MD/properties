package properties

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalKV__int(t *testing.T) {
	type S struct {
		A int    `properties:"a"`
		B int8   `properties:"b"`
		C int16  `properties:"c"`
		D int32  `properties:"d"`
		E int64  `properties:"e"`
		F uint   `properties:"f"`
		G uint8  `properties:"g"`
		H uint16 `properties:"h"`
		I uint32 `properties:"i"`
		J uint64 `properties:"j"`

		K int `properties:"k"`
	}

	var want = S{
		A: -1,
		B: 2,
		C: 512,
		D: 1024,
		E: 4096,
		F: 0,
		G: 2,
		H: 4,
		I: 8,
		J: 16,
		K: 0,
	}

	var input = map[string]string{
		"a": "-1",
		"b": "2",
		"c": "512",
		"d": "1024",
		"e": "4096",
		"f": "0",
		"g": "2",
		"h": "4",
		"i": "8",
		"j": "16",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__string(t *testing.T) {
	type S struct {
		A string `properties:"a"`
		B string `properties:"b"`
	}

	var want = S{
		A: "hello",
		B: "",
	}

	var input = map[string]string{
		"a": "hello",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__float(t *testing.T) {
	type S struct {
		A float32 `properties:"a"`
		B float64 `properties:"b"`

		C float32 `properties:"c"`
		D float32 `properties:"d"`
	}

	var want = S{
		A: 3.1415,
		B: 2.7187,
		C: 0,
		D: 0,
	}

	var input = map[string]string{
		"a": "3.1415",
		"b": "2.7187",
		"c": "0",
		"d": "0",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__bool(t *testing.T) {
	type S struct {
		A bool `properties:"a"`
		B bool `properties:"b"`
	}

	var want = S{
		A: true,
		B: false,
	}

	var input = map[string]string{
		"a": "true",
		"b": "0",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__map(t *testing.T) {
	type S struct {
		A map[string]string  `properties:"a"`
		B map[string]int     `properties:"b"`
		C map[string]float64 `properties:"c"`
		D map[string]string  `properties:"d"`
	}

	var want = S{
		A: map[string]string{
			"a": "hello",
			"b": "world",
		},
		B: map[string]int{
			"a": 1,
			"b": 2,
		},
		C: map[string]float64{
			"a": 3.1415,
			"b": 2.7187,
		},
		D: map[string]string{},
	}

	var input = map[string]string{
		"a.a": "hello",
		"a.b": "world",
		"b.a": "1",
		"b.b": "2",
		"c.a": "3.1415",
		"c.b": "2.7187",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__slice(t *testing.T) {
	type S struct {
		A []string `properties:"a"`
		B []int    `properties:"b"`

		C []bool `properties:"c"`
	}

	var want = S{
		A: []string{"hello", "world"},
		B: []int{1, 2, 3, 4, 5, 6, 7, 8},
		C: []bool{},
	}

	var input = map[string]string{
		"a[0]": "hello",
		"a[1]": "world",
		"b[0]": "1",
		"b[1]": "2",
		"b[2]": "3",
		"b[3]": "4",
		"b[4]": "5",
		"b[5]": "6",
		"b[6]": "7",
		"b[7]": "8",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__nested_struct(t *testing.T) {
	type AA struct {
		A string `properties:"a"`
		B int    `properties:"b"`
	}

	type BB struct {
		A bool   `properties:"a"`
		B string `properties:"b"`
	}

	type S struct {
		A AA `properties:"a"`
		B BB `properties:"b"`
	}

	var want = S{
		A: AA{
			A: "hello",
			B: 3,
		},
		B: BB{
			A: true,
			B: "world",
		},
	}

	var input = map[string]string{
		"a.a": "hello",
		"a.b": "3",
		"b.a": "true",
		"b.b": "world",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__nested_struct_pointer(t *testing.T) {
	type AA struct {
		A string `properties:"a"`
		B int    `properties:"b"`
	}

	type BB struct {
		A bool   `properties:"a"`
		B string `properties:"b"`
	}

	type S struct {
		A *AA `properties:"a"`
		B *BB `properties:"b"`
	}

	var want = S{
		A: &AA{
			A: "hello",
			B: 3,
		},
		B: &BB{
			A: true,
			B: "world",
		},
	}

	var input = map[string]string{
		"a.a": "hello",
		"a.b": "3",
		"b.a": "true",
		"b.b": "world",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__array_of_struct(t *testing.T) {
	type A struct {
		A string `properties:"a"`
		B int    `properties:"b"`
	}

	type S struct {
		AS []A `properties:"as"`
	}

	var want = S{
		AS: []A{
			{"hello", 1},
			{"world", 2},
		},
	}

	var input = map[string]string{
		"as[0].a": "hello",
		"as[0].b": "1",
		"as[1].a": "world",
		"as[1].b": "2",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__array_of_struct_pointers(t *testing.T) {
	type A struct {
		A string `properties:"a"`
		B int    `properties:"b"`
	}

	type S struct {
		AS []*A `properties:"as"`
	}

	var want = S{
		AS: []*A{
			{"hello", 1},
			{"world", 2},
		},
	}

	var input = map[string]string{
		"as[0].a": "hello",
		"as[0].b": "1",
		"as[1].a": "world",
		"as[1].b": "2",
	}

	var given S
	assert.NoError(t, unmarshalKV(input, &given))
	assert.Equal(t, want, given)
}

func TestUnmarshalKV__string_to_struct_map(t *testing.T) {
	type A struct {
		PA string `properties:"pa"`
	}

	type B struct {
		M map[string]A `properties:"m"`
	}

	var expectedB = B{
		M: map[string]A{
			"k1": {PA: "pa1"},
			"k2": {PA: "pa2"},
		},
	}

	var input = map[string]string{
		"m.k1.pa": "pa1",
		"m.k2.pa": "pa2",
	}

	var c B
	assert.NoError(t, unmarshalKV(input, &c))
	assert.Equal(t, expectedB, c)
}

func TestUnmarshalKV__string_to_struct_pointer_map(t *testing.T) {
	type A struct {
		PA string `properties:"pa"`
	}

	type B struct {
		M map[string]*A `properties:"m"`
	}

	var expectedB = B{
		M: map[string]*A{
			"k1": {PA: "pa1"},
			"k2": {PA: "pa2"},
		},
	}

	var input = map[string]string{
		"m.k1.pa": "pa1",
		"m.k2.pa": "pa2",
	}

	var b B
	assert.NoError(t, unmarshalKV(input, &b))
	assert.Equal(t, expectedB, b)
}

func TestPropsFromBytes(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		input := []byte(`
			a.a=hello
			a.b=world
			b[0].a=1
			b[0].b=2
			c=3.1415
			d=2.7187
			e={"a": 3, "b": "ha=ha=haha"}
		`)
		want := map[string]string{
			"a.a":    "hello",
			"a.b":    "world",
			"b[0].a": "1",
			"b[0].b": "2",
			"c":      "3.1415",
			"d":      "2.7187",
			"e":      "{\"a\": 3, \"b\": \"ha=ha=haha\"}",
		}

		p, err := propsFromBytes(input, "")
		assert.NoError(t, err)
		assert.Equal(t, want, p.kv)
	})

	t.Run("with comment and empty lines", func(t *testing.T) {
		input := []byte(`
			# comment 1
			a.a=hello
			a.b=world

			# comment 2
			b[0].a=1
			b[0].b=2

			# comment 3
			c=3.1415
			d=2.7187
		`)
		want := map[string]string{
			"a.a":    "hello",
			"a.b":    "world",
			"b[0].a": "1",
			"b[0].b": "2",
			"c":      "3.1415",
			"d":      "2.7187",
		}

		p, err := propsFromBytes(input, "")
		assert.NoError(t, err)
		assert.Equal(t, want, p.kv)
	})

	t.Run("with prefix", func(t *testing.T) {
		input := []byte(`
			a.a=hello
			a.b=world
		`)

		want := map[string]string{
			"a": "hello",
			"b": "world",
		}

		p, err := propsFromBytes(input, "a.")
		assert.NoError(t, err)
		assert.Equal(t, want, p.kv)
	})
}

func TestUnmarshalKey(t *testing.T) {
	type A struct {
		A string `properties:"a"`
		B int    `properties:"b"`
	}

	input := []byte(`
		a.a=hello
		a.b=1
		a=bye
		b=2
	`)

	var want = A{
		A: "hello",
		B: 1,
	}

	var given A
	err := UnmarshalKey("a", input, &given)
	assert.NoError(t, err)
	assert.Equal(t, want, given)
}
