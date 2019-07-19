# Properties 
[![CircleCI](https://circleci.com/gh/ZhengHe-MD/properties.svg?style=svg)](https://circleci.com/gh/ZhengHe-MD/properties)

a module for marshal/unmarshal .properties config file

[![GoDoc](https://godoc.org/github.com/ZhengHe-MD/properties?status.svg)](https://godoc.org/github.com/ZhengHe-MD/properties)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZhengHe-MD/properties)](https://goreportcard.com/report/github.com/ZhengHe-MD/properties)
![GitHub release](https://img.shields.io/github/release/ZhengHe-MD/properties.svg)

## API

1. Marshal

```go
func Marshal(v interface{}) ([]byte, error)
```

2. Unmarshal

```go
func Unmarshal(data []byte, v interface{}) error
```

3. UnmarshalKey

```go
func UnmarshalKey(key string, data []byte, v interface{}) error
```

4. UnmarshalKV

```go
func UnmarshalKV(kv map[string]string, v interface{}) error
```

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
	Name  string `properties:"name"`
	Phone string `properties:"phone"`
}

type Person struct {
	Name    string `properties:"name"`
	Age     int8   `properties:"age"`
	Email   string `properties:"email"`
	Bio     string `properties:"-"`
	Offline bool   `properties:"off"`

	EmergencyContact Contact   `properties:"emergency_contact"`
	AddressList      []Address `properties:"address_list"`
}

func main() {
	propsKV := map[string]string{
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
        "address_list[1].street": "Nanjing Road",
    }

    propsStr := `
    name=zhenghe
    age=18
    email=ranchardzheng@gmail.com
    bio=a boring guy
    offline=true
    emergency_contact.name=anonymous
    emergency_contact.phone=13333333333
    address_list[0].country=China
    address_list[0].city=Beijing
    address_list[0].street=Zhongguancun Street
    address_list[1].country=China
    address_list[1].city=Shanghai
    address_list[1].street=Nanjing Road
`

    var p1 Person
    var p2 Person
    var e1 Contact
    _ = properties.UnmarshalKV(propsKV, &p1)
    _ = properties.Unmarshal([]byte(propsStr), &p2)
    _ = properties.UnmarshalKey("emergency_contact", []byte(propsStr), &e1)

    fmt.Println(p1)
    fmt.Println(p2)
    fmt.Println(e1)

    var p3 = Person{
        Name: "zhenghe",
        Age: 18,
        Email: "ranchardzheng@gmail.com",
        Bio: "hahahaha",
        Offline: true,
        EmergencyContact: Contact{
            Name: "anonymous",
            Phone: "13333333333",
        },
        AddressList: []Address{
            {"China", "Beijing", "Zhongguancun Street"},
            {"China", "Shanghai", "Nanjing Road"},
        },
    }

    data, _ := properties.Marshal(p3)
    fmt.Println(string(data))
}
```

## Install

```sh
$ go get -u github.com/ZhengHe-MD/properties
```

