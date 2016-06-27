package util

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
	fields := GetFieldsSortedByTag(1, "schema") // int
	if fields != nil {
		t.Errorf("tags %v!=nil", fields)
		return
	}

	fields = GetFieldsSortedByTag("string", "schema") // string
	if fields != nil {
		t.Errorf("tags %v!=nil", fields)
		return
	}

	u := User{"kcl", "123", "a", "c", "b"}

	fields = GetFieldsSortedByTag(u, "schema")
	fmt.Printf("fields: %v\n", fields)

	expTags := []string{"a", "b", "name", "password"}
	expValues := []string{"a", "b", "kcl", "123"}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", f.Tag, expTags[i])
			return
		}
		if f.Field.String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", f.Field.String(), expValues[i])
			return
		}
	}

	fields = GetFieldsSortedByTag(u, "json")
	fmt.Printf("fields: %v\n", fields)

	expTags = []string{"1", "2", "3", "4"}
	expValues = []string{"c", "123", "kcl", "a"}

	for i, f := range fields {
		if f.Tag != expTags[i] {
			t.Errorf("expect: %v, got: %v", f.Tag, expTags[i])
			return
		}
		if f.Field.String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", f.Field.String(), expValues[i])
			return
		}
	}

}
