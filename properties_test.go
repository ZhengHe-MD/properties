package properties

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalMap__int(t *testing.T) {
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

	var expected = S{
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

	var s S
	assert.NoError(t, unmarshalKV(input, &s))
	assert.Equal(t, expected, s)
}

func TestUnmarshalMap__string(t *testing.T) {
	type S struct {
		A string `properties:"a"`
		B string `properties:"b"`
	}

	var expected = S{
		A: "hello",
		B: "",
	}

	var input = map[string]string{
		"a": "hello",
	}

	var s S
	assert.NoError(t, unmarshalKV(input, &s))
	assert.Equal(t, expected, s)
}

func TestUnmarshalMap__float(t *testing.T) {
	type S struct {
		A float32 `properties:"a"`
		B float64 `properties:"b"`

		C float32 `properties:"c"`
		D float32 `properties:"d"`
	}

	var expected = S{
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

	var s S
	assert.NoError(t, unmarshalKV(input, &s))
	assert.Equal(t, expected, s)
}

func TestUnmarshalMap__bool(t *testing.T) {
	type S struct {
		A bool `properties:"a"`
		B bool `properties:"b"`
	}

	var expected = S{
		A: true,
		B: false,
	}

	var input = map[string]string{
		"a": "true",
		"b": "0",
	}

	var s S
	assert.NoError(t, unmarshalKV(input, &s))
	assert.Equal(t, expected, s)
}

func TestUnmarshalMap__map(t *testing.T) {
	type S struct {
		A map[string]string  `properties:"a"`
		B map[string]int     `properties:"b"`
		C map[string]float64 `properties:"c"`

		D map[string]string `properties:"d"`
	}

	var expected = S{
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

	var s S
	assert.NoError(t, unmarshalKV(input, &s))
	assert.Equal(t, expected, s)
}

func TestUnmarshalMap__slice(t *testing.T) {
	type S struct {
		A []string `properties:"a"`
		B []int    `properties:"b"`

		C []bool `properties:"c"`
	}

	var expected = S{
		A: []string{"hello", "world"},
		B: []int{1, 2},
		C: []bool{},
	}

	var input = map[string]string{
		"a[0]": "hello",
		"a[1]": "world",
		"b[0]": "1",
		"b[1]": "2",
	}

	var s S
	assert.NoError(t, unmarshalKV(input, &s))
	assert.Equal(t, expected, s)
}

func TestUnmarshalMap__complex_usages(t *testing.T) {
	t.Run("nested struct", func(t *testing.T) {
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

		var expected = S{
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

		var s S
		assert.NoError(t, unmarshalKV(input, &s))
		assert.Equal(t, expected, s)
	})

	t.Run("array of struct", func(t *testing.T) {
		type A struct {
			A string `properties:"a"`
			B int    `properties:"b"`
		}

		type S struct {
			AS []A `properties:"as"`
		}

		var expected = S{
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

		var s S
		assert.NoError(t, unmarshalKV(input, &s))
		assert.Equal(t, expected, s)
	})
}
