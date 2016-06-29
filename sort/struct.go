package sort

// package sort provides primitives for sorting method in gas system

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

// TagField represents field value along with tag in any struct
type TagField struct {
	Tag   string
	Value reflect.Value
}

// TagFields wraps TagField slice, and provides Encode method
// TagFields can be sort by tag in alphabet order
type TagFields []TagField

// Len for sort
func (t TagFields) Len() int {
	return len(t)
}

// Swap for sort
func (t TagFields) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less for sort
func (t TagFields) Less(i, j int) bool {
	return strings.Compare(t[i].Tag, t[j].Tag) < 0
}

// Encode encodes TagFields into string by joining [tag=value, tag=value, ...]
// seperated by given sep.
// Example: given sep=&, and t=[TagField{"tag1", "val1"}, TagField{"tag2", "val2}, ...]
// "tag1=val1&tag2=val2&...&tagn=valn" is return
func (t TagFields) Encode(sep string) string {
	var buffer bytes.Buffer
	for i, f := range t {
		if i != 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(f.Tag)
		buffer.WriteString("=")
		buffer.WriteString(fmt.Sprintf("%v", f.Value))
	}
	return buffer.String()
}

// EncodeValOnly encodes TagFields into string by joining [value, value, ...]
// seperated by given sep
// Example: given sep=&, and t=[TagField{"tag1", "val1"}, TagField{"tag2", "val2}, ...]
// "val1&val2&...&valn" is return
func (t TagFields) EncodeValOnly(sep string) string {
	var buffer bytes.Buffer
	for i, f := range t {
		if i != 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(fmt.Sprintf("%v", f.Value))
	}
	return buffer.String()
}

// GetFields calls GetFieldsByTag with param tagKey="sort"
func GetFields(v interface{}) (fields TagFields) {
	return GetFieldsByTag(v, "sort")
}

// GetFieldsByTag returns tags and its' corresponding values sorted
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
func GetFieldsByTag(v interface{}, tagKey string) (fields TagFields) {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct { // type not struct
		return nil
	}

	obj := reflect.ValueOf(v)
	var tfs []TagField
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get(tagKey)
		if t == "" || t == "-" { // escape empty tag and "-" tag
			continue
		}
		name, opts := parseTag(t)
		val := obj.FieldByName(f.Name)
		if opts.Contains("omitempty") && isEmptyValue(val) {
			continue
		}

		tf := TagField{
			Tag:   name,
			Value: val,
		}
		tfs = append(tfs, tf)
	}
	fields = TagFields(tfs)
	sort.Sort(fields)
	return fields
}

// tagOptions is the string following a comma in a struct field's "url" tag, or
// the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}
