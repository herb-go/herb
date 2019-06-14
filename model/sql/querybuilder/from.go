package querybuilder

// NewFromClause create new form clause
func (b *Builder) NewFromClause() *FromClause {
	return &FromClause{
		Buidler: b,
		Tables:  [][2]string{},
	}

}

//FromClause from caluse struuct
type FromClause struct {
	Buidler *Builder
	Tables  [][2]string
}

// AddAlias add table to from clause with given table name and alias
func (q *FromClause) AddAlias(alias string, tableName string) *FromClause {
	q.Tables = append(q.Tables, [2]string{tableName, alias})
	return q
}

// Add add table to form clause
func (q *FromClause) Add(tableName string) *FromClause {
	q.Tables = append(q.Tables, [2]string{tableName, ""})
	return q
}

// QueryCommand return query command.
func (q *FromClause) QueryCommand() string {
	var command = ""
	command = "FROM "
	for k := range q.Tables {
		if q.Tables[k][1] == "" {
			command += q.Tables[k][0] + " , "
		} else {
			command += q.Tables[k][0] + " AS " + q.Tables[k][1] + " , "
		}
	}
	if len(q.Tables) > 0 {
		command = command[:len(command)-3]
	}
	return command
}

// QueryArgs return query args.
func (q *FromClause) QueryArgs() []interface{} {
	return []interface{}{}
}
