package query

import (
	"reflect"
	"strings"
)

func NewValueList(data ...interface{}) *PlainQuery {
	if len(data) == 0 {
		return New("")
	}
	var command = strings.Repeat("? , ", len(data))
	return New(command[:len(command)-3], data...)
}

func In(field string, args interface{}) *PlainQuery {
	var argsvalue = reflect.ValueOf(args)
	var data = make([]interface{}, argsvalue.Len())
	for k := range data {
		data[k] = argsvalue.Index(k).Interface()
	}
	var query = NewValueList(data...)
	query.Command = field + " IN ( " + query.Command + " )"
	return query
}

func Equal(field string, arg interface{}) *PlainQuery {
	return New(field+" = ?", arg)
}
func Search(field string, arg string) *PlainQuery {
	if arg == "" || field == "" {
		return New("")
	}
	return New(field+" LIKE ?", "%"+EscapeSearch(arg)+"%")
}

func EscapeSearch(command string) string {
	command = strings.Replace(command, "\\", "\\\\", -1)
	command = strings.Replace(command, "_", "\\_", -1)
	command = strings.Replace(command, "%", "\\%", -1)
	return command
}
