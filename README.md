# Properties

a module for marshal/unmarshal .properties config file

[![CircleCI](https://circleci.com/gh/ZhengHe-MD/properties.svg?style=svg)](https://circleci.com/gh/ZhengHe-MD/properties)
[![GoDoc](https://godoc.org/github.com/ZhengHe-MD/properties?status.svg)](https://godoc.org/github.com/ZhengHe-MD/properties)
![GitHub release](https://img.shields.io/github/release/ZhengHe-MD/properties.svg)

NOTE: this project is not production ready, apis may change if it's necessary

## Usages

```go
import (
	"fmt"
	"github.com/ZhengHe-MD/properties"
)

type Address struct {
	Country string `properties:"country"`
	City    string `properties:"city"`
	Street  string `properties:"street"`
}

type Contact struct {
	Name string `properties:"name"`
	Phone string `properties:"phone"`
}

type Person struct {
	Name string `properties:"name"`
	Age    int8 `properties:"age"`
	Email string `properties:"email"`
	Bio string `properties:"-"`
	Offline bool `properties:"off"`

	EmergencyContact Contact `properties:"emergency_contact"`
	AddressList []Address `properties:"address_list"`
}

func main() {
	props := map[string]string{
		"name": "zhenghe",
		"age": "18",
		"email": "ranchardzheng@gmail.com",
		"bio": "a boring guy",
		"offline": "true",
		"emergency_contact.name": "anonymous",
		"emergency_contact.phone": "13333333333",
		"address_list[0].country": "China",
		"address_list[0].city": "Beijing",
		"address_list[0].street": "Zhongguancun Street",
		"address_list[1].country": "China",
		"address_list[1].city": "Shanghai",
		"address_list[1].street": "Nanjing Street",
	}

	var p Person
	_ = properties.UnmarshalKV(props, &p)
	fmt.Println(p)
}
```

## API

1. UnmarshalKV

```go
func UnmarshalKV(kv map[string]string, v interface{}) error
```

2. Marshal (TODO)

```go
func Marshal(v interface{}) ([]byte, error)
```

3. Unmarshal (TODO)

```go
func Unmarshal(data []byte, v interface{}) error
```

## Install

```sh
$ go get -u github.com/ZhengHe-MD/properties
```

