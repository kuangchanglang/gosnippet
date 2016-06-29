package sort

import (
	"fmt"
	"testing"
)

type User struct {
	Name   string `schema:"name" json:"3"`
	Passwd string `schema:"password" json:"2"`
	A      string `schema:"a" json:"4"`
	C      string `json:"1"`
	b      string `schema:"b"`
}

func TestGetFieldSortedByTag(t *testing.T) {
	fields := GetFieldsByTag(1, "schema") // int
	if fields != nil {
		t.Errorf("tags %v!=nil", fields)
		return
	}

	fields = GetFieldsByTag("string", "schema") // string
	if fields != nil {
		t.Errorf("tags %v!=nil", fields)
		return
	}

	u := User{"kcl", "123", "a", "c", "b"}

	fields = GetFieldsByTag(u, "schema")
	fmt.Printf("fields: %v\n", fields)

	expTags := []string{"a", "b", "name", "password"}
	expValues := []string{"a", "b", "kcl", "123"}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", f.Tag, expTags[i])
			return
		}
		if f.Value.String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", f.Value.String(), expValues[i])
			return
		}
	}

	fields = GetFieldsByTag(u, "json")
	fmt.Printf("fields: %v\n", fields)

	expTags = []string{"1", "2", "3", "4"}
	expValues = []string{"c", "123", "kcl", "a"}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", f.Tag, expTags[i])
			return
		}
		if f.Value.String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", f.Value.String(), expValues[i])
			return
		}
	}
}

type Foo struct {
	A int     `schema:"aa"`
	B float64 `schema:"00"`
	C string  `schema:"cc,omitempty" sort:"haha"`
	D string  `schema:"1"`
	E int     `schema:"e,omitempty"`
	F int     `schema:"f,omitempty"`
}

var foo = Foo{3, 8.33, "", "d", 0, 9}

func TestGetFields(t *testing.T) {
	// get tagkey "sort" as default
	fields := GetFields(foo)
	if fields.Len() != 1 {
		t.Errorf("fields length expect 1, got:%v", fields.Len())
		return
	}
}

func TestEncode(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.Encode("&")
	fmt.Println(str)
	expect := "00=8.33&1=d&aa=3&f=9"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}

func TestEncodeValOnly(t *testing.T) {
	fields := GetFieldsByTag(foo, "schema")
	str := fields.EncodeValOnly("")
	expect := "8.33d39"
	if str != expect {
		t.Errorf("expect: %v, got: %v", expect, str)
		return
	}
}

func BenchmarkGetFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFieldsByTag(foo, "schema")
	}
}
