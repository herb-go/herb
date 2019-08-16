package querybuilder

import (
	"strconv"
	"strings"

	driver "github.com/lib/pq"
)

// PostgreSQLBuilderDriver posgresql builder driver struct
type PostgreSQLBuilderDriver struct {
	EmptyBuilderDriver
}

//IsDuplicate check if error is Is duplicate error.
func (d *PostgreSQLBuilderDriver) IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*driver.Error)
	if ok == false {
		return false
	}
	return e.Code == "23505"
}

//ConvertQuery  convert query to command and args
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

// PostgreSQLDriver postgre sql driver
var PostgreSQLDriver = &PostgreSQLBuilderDriver{}

func init() {
	RegisterDriver("postgres", PostgreSQLDriver)
}
