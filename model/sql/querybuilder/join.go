package builder

import "strings"

type JoinData struct {
	Type         string
	Table        [2]string
	Condition    *PlainQuery
	UsingColnums []string
}

func (d *JoinData) Alias(alias string, tableName string) *JoinData {
	d.Table[0] = tableName
	d.Table[1] = alias
	return d
}

func (d *JoinData) AliasTable(t Table) *JoinData {
	d.Table[0] = t.TableName()
	d.Table[1] = t.TableName()
	return d
}

func (d *JoinData) On(condition *PlainQuery) *JoinData {
	d.Condition = condition
	return d
}
func (d *JoinData) Using(colnums ...string) *JoinData {
	d.UsingColnums = colnums
	return d
}

func (d *JoinData) QueryCommand() string {
	var command = d.Type + " Join "
	command += d.Table[indexTableName]
	if d.Table[indexAlias] != "" {
		command += " AS " + d.Table[indexAlias]
	}
	if len(d.UsingColnums) == 0 {
		command += " ON " + d.Condition.QueryCommand()
	} else {
		command += " USING (" + strings.Join(d.UsingColnums, " , ") + ")"
	}
	return command
}
func (q *JoinData) QueryArgs() []interface{} {
	if q.Condition != nil && len(q.UsingColnums) == 0 {
		return q.Condition.QueryArgs()
	}
	return []interface{}{}
}

func (b *Builder) NewJoinQuery() *JoinQuery {
	return &JoinQuery{
		Builder: b,
		Data:    []*JoinData{},
	}
}

type JoinQuery struct {
	Builder *Builder
	Data    []*JoinData
}

func (q *JoinQuery) join(jointype string) *JoinData {
	data := &JoinData{
		Type:         jointype,
		Table:        [2]string{},
		Condition:    nil,
		UsingColnums: []string{},
	}
	q.Data = append(q.Data, data)
	return data
}

func (q *JoinQuery) InnerJoin() *JoinData {
	return q.join("INNER")
}
func (q *JoinQuery) LeftJoin() *JoinData {
	return q.join("LEFT")
}
func (q *JoinQuery) RightJoin() *JoinData {
	return q.join("RIGHT")
}

func (q *JoinQuery) QueryCommand() string {
	var command = ""
	for k := range q.Data {
		c := q.Data[k].QueryCommand()
		if c != "" {
			command += c + "\n"
		}
	}
	if command != "" {
		command = command[:len(command)-1]
	}
	return command
}
func (q *JoinQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		a := q.Data[k].QueryArgs()
		if a != nil {
			args = append(args, a...)
		}
	}
	return args
}
