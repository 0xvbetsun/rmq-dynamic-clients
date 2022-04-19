package ordered

import (
	"fmt"
	"reflect"
)

func ExampleNewMap() {
	m := NewMap()
	fmt.Println(reflect.TypeOf(m), m.Len())
	// Output: *ordered.Map 0
}

func ExampleMap_Len() {
	m := NewMap()
	fmt.Println(m.Len())
	m.Add("abc")
	fmt.Println(m.Len())
	// Output:
	// 0
	// 1
}

func ExampleMap_Add() {
	m := NewMap()
	fmt.Println(m.Add("abc"), m.Add("bcd"), m.Add("abc"))
	// Output: abc bcd abc
}

func ExampleMap_Get() {
	m := NewMap()
	m.Add("abc")
	fmt.Println(m.Get("abc"), m.Get("bcd"))
	// Output: abc <nil>
}

func ExampleMap_Keys() {
	m := NewMap()
	m.Add("abc")
	m.Add("bcd")
	fmt.Println(m.Keys())
	// Output: [abc bcd]
}

func ExampleMap_Delete() {
	m := NewMap()
	m.Add("abc")
	fmt.Println(m.Delete("abc"), m.Delete("abc"))
	// Output: true false
}
