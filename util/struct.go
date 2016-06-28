package util

import (
	"reflect"
	"sort"
)

// TagField represents field that with tag in struct
type TagField struct {
	Tag   string
	Value reflect.Value
}

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
//    fields := GetFieldsSortedByTag(u, "schema")
//    // fields =[{"age", 20}, {"name", "alice"}, {"passwd", "0000"}]
//    GetFieldsSortedByTag("string") // nil is return
//
func GetFieldsSortedByTag(v interface{}, tagKey string) (fields []TagField) {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct { // type not struct
		return nil
	}

	tags := make([]string, 0)           // tags, for sorting
	fieldMap := make(map[string]string) // key: tag, value: field name
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get(tagKey)
		if t == "" { // escape empty tag
			continue
		}
		fieldMap[t] = f.Name
		tags = append(tags, t)
	}
	sort.Strings(tags)

	value := reflect.ValueOf(v)
	fields = make([]TagField, 0)
	for _, t := range tags {
		fieldName := fieldMap[t]
		f := TagField{
			Tag:   t,
			Value: value.FieldByName(fieldName),
		}
		fields = append(fields, f)
	}

	return fields
}
