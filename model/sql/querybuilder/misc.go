package builder

import (
	"reflect"
	"strings"
)

func (b *Builder) NewValueList(data ...interface{}) *PlainQuery {
	if len(data) == 0 {
		return b.New("")
	}
	var command = strings.Repeat("? , ", len(data))
	return b.New(command[:len(command)-3], data...)
}

func (b *Builder) In(field string, args interface{}) *PlainQuery {
	var argsvalue = reflect.ValueOf(args)
	var data = make([]interface{}, argsvalue.Len())
	for k := range data {
		data[k] = argsvalue.Index(k).Interface()
	}
	var query = b.NewValueList(data...)
	query.Command = field + " IN ( " + query.Command + " )"
	return query
}

func (b *Builder) Equal(field string, arg interface{}) *PlainQuery {
	return b.New(field+" = ?", arg)
}
func (b *Builder) Search(field string, arg string) *PlainQuery {
	if arg == "" || field == "" {
		return b.New("")
	}
	return b.New(field+" LIKE ?", "%"+b.EscapeSearch(arg)+"%")
}

func (b *Builder) EscapeSearch(command string) string {
	command = strings.Replace(command, "\\", "\\\\", -1)
	command = strings.Replace(command, "_", "\\_", -1)
	command = strings.Replace(command, "%", "\\%", -1)
	return command
}
