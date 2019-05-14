package querybuilder

import (
	"strconv"
	"strings"
)

type PostgreSQLBuilderDriver struct {
	EmptyBuilderDriver
}

func (d PostgreSQLBuilderDriver) ConvertQuery(q Query) (string, []interface{}) {
	cmd := q.QueryCommand()
	arg := q.QueryArgs()
	converted := ""
	var sum = 1
	for len(cmd) != 0 {
		index := strings.Index(cmd, "?")
		if index < 0 {
			converted += cmd
			break
		}
		if index == 0 {
			converted += "$" + strconv.Itoa(sum)
			sum++
			cmd = cmd[1:]
			continue
		}
		if cmd[index-1:index] == "\\" {
			converted += cmd[:index-1]
			converted += "q"
			cmd = cmd[index+1:]
			continue
		}
		converted += cmd[:index]
		converted += "$" + strconv.Itoa(sum)
		sum++
		cmd = cmd[index+1:]
	}
	return converted, arg
}

var PostgreSQLDriver = &PostgreSQLBuilderDriver{}

func init() {
	RegisterDriver("postgres", PostgreSQLDriver)
}
