package database

import (
	"fmt"
	"reflect"
	"strings"
)

// BuildInsertQuery - Builds a query string from a struct, accepts pointer as an argument
func BuildInsertQuery(table string, inf interface{}) string {
	var query string
	v := reflect.ValueOf(inf).Elem()
	t := v.Type()

	query = fmt.Sprintf("INSERT INTO %s (", table)
	for i := 0; i < v.NumField(); i++ {
		query += getCorrectQueryAppendantFromTag(t.Field(i).Tag.Get("json"), v.Field(i))
	}
	query = query[:len(query)-1] + ") VALUES ("
	for i := 0; i < v.NumField(); i++ {
		query += getCorrectQueryAppendantFromValue(v.Field(i))
	}
	query = query[:len(query)-1] + ")"
	return query
}

func getCorrectQueryAppendantFromTag(t string, v reflect.Value) string {
	var val string
	if v.CanInterface() {
		val = fmt.Sprint(v.Interface())
	}

	if val == "" {
		return ""
	}

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
		switch v.Type().String() {
		case "int":
			return val + ","
		case "string":
			return "'" + val + "',"
		case "bool":
			if val == "true" {
				return "TRUE,"
			}
			return "FALSE,"
		default:
			return val + ","
		}
	}

	return ""
}
