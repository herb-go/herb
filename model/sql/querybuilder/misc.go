package querybuilder

import (
	"reflect"
	"strings"
)

// NewValueList create new value list with given data.
func (b *Builder) NewValueList(data ...interface{}) *PlainQuery {
	if len(data) == 0 {
		return b.New("")
	}
	var command = strings.Repeat("? , ", len(data))
	return b.New(command[:len(command)-3], data...)
}

//In return field in args query
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

// Equal return field equal to value query
func (b *Builder) Equal(field string, arg interface{}) *PlainQuery {
	return b.New(field+" = ?", arg)
}

// IsNull return field is null  query
func (b *Builder) IsNull(field string) *PlainQuery {
	return b.New(field + " IS NULL")
}

// IsNotNull return field is not null query
func (b *Builder) IsNotNull(field string) *PlainQuery {
	return b.New(field + " IS NOT NULL")
}

//Between return field between start and end query
func (b *Builder) Between(field string, start interface{}, end interface{}) *PlainQuery {
	return b.New(field+" BETWEEN ? AND ?", start, end)
}

//Search return search field query.
//  arg will be ecapced with EscapeSearch
func (b *Builder) Search(field string, arg string) *PlainQuery {
	if arg == "" || field == "" {
		return b.New("")
	}
	return b.New(field+" LIKE ?", "%"+b.EscapeSearch(arg)+"%")
}

// EscapeSearch escape search arg
func (b *Builder) EscapeSearch(arg string) string {
	arg = strings.Replace(arg, "\\", "\\\\", -1)
	arg = strings.Replace(arg, "_", "\\_", -1)
	arg = strings.Replace(arg, "%", "\\%", -1)
	return arg
}
