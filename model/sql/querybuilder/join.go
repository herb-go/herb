package querybuilder

//JoinData join clause struct
type JoinData struct {
	Builder *Builder
	// Type join type
	Type string
	// TableJoin table struct.
	// Table[0]:table name
	//Table[1]:table alias
	Table [2]string
	// Condition join condition
	Condition *PlainQuery
}

// Alias add joined table with given alias.
func (d *JoinData) Alias(alias string, tableName string) *JoinData {
	d.Table[0] = tableName
	d.Table[1] = alias
	return d
}

// On add On condition
func (d *JoinData) On(condition *PlainQuery) *JoinData {
	d.Condition = condition
	return d
}

//OnEqual On add On condition which given fieds equal
func (d *JoinData) OnEqual(field1 string, field2 string) *JoinData {
	d.Condition = d.Builder.New(field1 + "=" + field2)
	return d
}

// QueryCommand return query command
func (d *JoinData) QueryCommand() string {
	var command = d.Type + " JOIN "
	command += d.Table[indexTableName]
	if d.Table[indexAlias] != "" {
		command += " AS " + d.Table[indexAlias]
	}
	command += " ON " + d.Condition.QueryCommand()
	return command
}

// QueryArgs return query args
func (d *JoinData) QueryArgs() []interface{} {
	if d.Condition != nil {
		return d.Condition.QueryArgs()
	}
	return []interface{}{}
}

// NewJoinClause create new join clause
func (b *Builder) NewJoinClause() *JoinClause {
	return &JoinClause{
		Builder: b,
		Data:    []*JoinData{},
	}
}

// JoinClause query struct
type JoinClause struct {
	Builder *Builder
	// Data join data list
	Data []*JoinData
}

func (q *JoinClause) join(jointype string) *JoinData {
	data := &JoinData{
		Type:      jointype,
		Table:     [2]string{},
		Builder:   q.Builder,
		Condition: nil,
	}
	q.Data = append(q.Data, data)
	return data
}

// InnerJoin set type of join clause  to INNER
func (q *JoinClause) InnerJoin() *JoinData {
	return q.join("INNER")
}

// LeftJoin set type of join clause  to LEFT
func (q *JoinClause) LeftJoin() *JoinData {
	return q.join("LEFT")
}

// RightJoin set type of join clause  to RIGHT
func (q *JoinClause) RightJoin() *JoinData {
	return q.join("RIGHT")
}

// QueryCommand return query command
func (q *JoinClause) QueryCommand() string {
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

// QueryArgs return query args
func (q *JoinClause) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		a := q.Data[k].QueryArgs()
		if a != nil {
			args = append(args, a...)
		}
	}
	return args
}
