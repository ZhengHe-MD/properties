package properties

import (
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshal__map(t *testing.T) {
	t.Run("marshal map", func(t *testing.T) {
		var m = map[string]string{
			"a": "hello",
			"b": "world",
		}

		expectedData1 := []byte(strings.Join([]string{
			"a=hello\n",
			"b=world\n",
		}, ""))

		expectedData2 := []byte(strings.Join([]string{
			"b=world\n",
			"a=hello\n",
		}, ""))

		data, err := Marshal(m)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(expectedData1, data) || reflect.DeepEqual(expectedData2, data))
	})

	t.Run("marshal map pointer", func(t *testing.T) {
		var m = map[string]string{
			"a": "hello",
			"b": "world",
		}

		expectedData1 := []byte(strings.Join([]string{
			"a=hello\n",
			"b=world\n",
		}, ""))

		expectedData2 := []byte(strings.Join([]string{
			"b=world\n",
			"a=hello\n",
		}, ""))

		data, err := Marshal(&m)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(expectedData1, data) || reflect.DeepEqual(expectedData2, data))
	})

	t.Run("marshal map complex usage", func(t *testing.T) {
		var m = map[string]string{
			"dailystudy.leveltext[0].level": "100000",
			"dailystudy.leveltext[0].ls":    "BR-100L",
			"dailystudy.leveltext[0].title": "完成Level A的学习，孩子可以",
			"dailystudy.leveltext[0].text":  "\"1. 认识26个字母及字母相关单词78个\n2. 阅读关于身体,动物等9个话题绘本70本\n3. 学习口语词汇500+、常用口语表达句式12个\"",
		}

		data, err := Marshal(&m)
		log.Println(string(data))
		assert.NoError(t, err, string(data))
	})
}

func TestMarshal__complex_usages(t *testing.T) {
	type A struct {
		SA1 string `properties:"sa1"`
		IA1 int    `properties:"ia1"`
	}

	type S struct {
		S1      string            `properties:"s1"`
		I1      int               `properties:"i1"`
		I2      int8              `properties:"i2"`
		I3      int16             `properties:"i3"`
		I4      int32             `properties:"i4"`
		I5      int64             `properties:"i5"`
		UI1     uint              `properties:"ui1"`
		UI2     uint8             `properties:"ui2"`
		UI3     uint16            `properties:"ui3"`
		UI4     uint32            `properties:"ui4"`
		UI5     uint64            `properties:"ui5"`
		F1      float32           `properties:"f1"`
		F2      float64           `properties:"f2"`
		B1      bool              `properties:"b1"`
		B2      bool              `properties:"b2"`
		Slice1  []int             `properties:"slice1"`
		Slice2  []*A              `properties:"slice2"`
		Map1    map[string]string `properties:"map1"`
		Map2    map[string]*A     `properties:"map2"`
		Pt1     *A                `properties:"pt1"`
		Struct1 A                 `properties:"st1"`
		Time    time.Time         `properties:"time"`
	}

	var s = S{
		S1:     "hello world",
		I1:     1,
		I2:     2,
		I3:     4,
		I4:     8,
		I5:     16,
		UI1:    1,
		UI2:    2,
		UI3:    4,
		UI4:    8,
		UI5:    16,
		F1:     3.1415,
		F2:     2.7187,
		B1:     true,
		B2:     false,
		Slice1: []int{1, 2, 3, 4},
		Slice2: []*A{
			{"hello", 1},
			{"world", 2},
		},
		Map1: map[string]string{
			"a": "haha",
		},
		Map2: map[string]*A{
			"a": {"hello", 1},
		},
		Pt1:     &A{"byebye", 3},
		Struct1: A{"morning", 4},
		Time:    time.Date(2021, 8, 30, 11, 11, 11, 11, time.UTC),
	}

	expectedLines := []string{
		"s1=hello world\n",
		"i1=1\n",
		"i2=2\n",
		"i3=4\n",
		"i4=8\n",
		"i5=16\n",
		"ui1=1\n",
		"ui2=2\n",
		"ui3=4\n",
		"ui4=8\n",
		"ui5=16\n",
		"f1=3.1415\n",
		"f2=2.7187\n",
		"b1=true\n",
		"b2=false\n",
		"slice1[0]=1\n",
		"slice1[1]=2\n",
		"slice1[2]=3\n",
		"slice1[3]=4\n",
		"slice2[0].sa1=hello\n",
		"slice2[0].ia1=1\n",
		"slice2[1].sa1=world\n",
		"slice2[1].ia1=2\n",
		"map1.a=haha\n",
		"map2.a.sa1=hello\n",
		"map2.a.ia1=1\n",
		"pt1.sa1=byebye\n",
		"pt1.ia1=3\n",
		"st1.sa1=morning\n",
		"st1.ia1=4\n",
		"time=2021-08-30T11:11:11.000000011Z",
	}

	expectedData := []byte(strings.Join(expectedLines, ""))
	data, err := Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)

	println(string(data))
}
