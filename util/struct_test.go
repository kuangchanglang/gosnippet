package util

import (
	"fmt"
	"reflect"
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
	tags, values := GetFieldsSortedByTag(1, "schema") // int
	if tags != nil {
		t.Errorf("tags %v!=nil", tags)
		return
	}
	if values != nil {
		t.Errorf("values %v!=nil", values)
		return
	}

	tags, values = GetFieldsSortedByTag("string", "schema") // string
	if tags != nil {
		t.Errorf("tags %v!=nil", tags)
		return
	}
	if values != nil {
		t.Errorf("values %v!=nil", values)
		return
	}

	u := User{"kcl", "123", "a", "c", "b"}

	tags, values = GetFieldsSortedByTag(u, "schema")
	fmt.Printf("tags: %v, values: %v\n", tags, values)

	expTags := []string{"a", "b", "name", "password"}
	expValues := []string{"a", "b", "kcl", "123"}

	if reflect.DeepEqual(tags, expTags) == false {
		t.Errorf("expect: %v, got: %v", expTags, tags)
		return
	}
	for i, _ := range values {
		if values[i].String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", values, expValues)
			return
		}
	}

	tags, values = GetFieldsSortedByTag(u, "json")
	fmt.Printf("tags: %v, values: %v\n", tags, values)

	expTags = []string{"1", "2", "3", "4"}
	expValues = []string{"c", "123", "kcl", "a"}

	if reflect.DeepEqual(tags, expTags) == false {
		t.Errorf("expect: %v, got: %v", expTags, tags)
		return
	}
	for i, _ := range values {
		if values[i].String() != expValues[i] {
			t.Errorf("expect: %v, got: %v", values, expValues)
			return
		}
	}

}
