package util

import (
	"reflect"
	"sort"
)

// GetFieldsSortedByTag returns tags and its' corresponding values sorted
// by given tagKey in alphabet order.
// (nil, nil) is return if v is not struct
// otherwise, len(tags)==len(values) is guaranteed
//
// Example:
//    type User struct{
//			Passwd string `schema:"passwd"`
//		    Name   string `schema:"name"`
//          Age    int    `schema:"age"`
//    }
//    u := User{"0000", "alice", 20}
//    tags, values := GetFieldsSortedByTag(u, "schema")
//    // tags=["age", "name", "passwd"], values=[20, "alice", "0000]
//    GetFieldsSortedByTag("string") // nil, nil is return
//
func GetFieldsSortedByTag(v interface{}, tagKey string) (tags []string, values []reflect.Value) {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct { // type not struct
		return nil, nil
	}

	tags = make([]string, 0)          // tags, for sorting
	fields := make(map[string]string) // key: tag, value: field name
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get(tagKey)
		if t == "" { // escape empty tag
			continue
		}
		fields[t] = f.Name
		tags = append(tags, t)
	}
	sort.Strings(tags)

	value := reflect.ValueOf(v)
	values = make([]reflect.Value, 0)
	for _, t := range tags {
		fieldName := fields[t]
		values = append(values, value.FieldByName(fieldName))
	}

	return tags, values
}
