package querybuilder

func (b *Builder) NewFromClause() *FromClause {
	return &FromClause{
		Tables: [][2]string{},
	}

}

type FromClause struct {
	Tables [][2]string
}

func (q *FromClause) AddAlias(alias string, tableName string) *FromClause {
	q.Tables = append(q.Tables, [2]string{tableName, alias})
	return q
}

func (q *FromClause) Add(tableName string) *FromClause {
	q.Tables = append(q.Tables, [2]string{tableName, ""})
	return q
}

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
func (q *FromClause) QueryArgs() []interface{} {
	return []interface{}{}
}
