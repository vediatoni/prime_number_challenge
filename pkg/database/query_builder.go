package database

import (
	"fmt"
	"reflect"
	"strings"
)

// BuildInsertQuery - Builds a query string from a struct, accepts pointer as an argument
func BuildInsertQuery(inf interface{}) string {
	var query string
	v := reflect.ValueOf(inf).Elem()
	t := v.Type()

	query = "INSERT INTO prime_numbers ("
	for i := 0; i < v.NumField(); i++ {
		query += getCorrectQueryAppendantFromTag(t.Field(i).Tag.Get("json"))
	}
	query = query[:len(query)-1] + ") VALUES("
	for i := 0; i < v.NumField(); i++ {
		query += getCorrectQueryAppendantFromValue(v.Field(i))
	}
	query = query[:len(query)-1] + ");"
	return query
}

func getCorrectQueryAppendantFromTag(t string) string {
	tag := t
	if strings.Contains(tag, ",omitempty") {
		tag = strings.Replace(tag, ",omitempty", "", -1)
	}

	if tag != "" {
		return tag + ","
	}

	return ""
}

func getCorrectQueryAppendantFromValue(v reflect.Value) string {
	var val string

	// If the value is exported , we can just get the value
	if v.CanInterface() {
		val = fmt.Sprint(v.Interface())
	}

	if val != "" {
		return val + ","
	}

	return ""
}
