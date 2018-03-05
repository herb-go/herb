package querybuilder

func (b *Builder) NewFromQuery() *FromQuery {
	return &FromQuery{
		Tables: [][2]string{},
	}

}

type FromQuery struct {
	Tables [][2]string
}

func (q *FromQuery) AddTable(t Table) *FromQuery {
	return q.AddAlias(t.Alias(), t.TableName())
}
func (q *FromQuery) AddAlias(alias string, tableName string) *FromQuery {
	q.Tables = append(q.Tables, [2]string{tableName, alias})
	return q
}

func (q *FromQuery) Add(tableName string) *FromQuery {
	q.Tables = append(q.Tables, [2]string{tableName, ""})
	return q
}

func (q *FromQuery) QueryCommand() string {
	var command = ""
	if len(q.Tables) == 0 {
		return command
	}
	command = "FROM "
	for k := range q.Tables {
		if q.Tables[k][1] == "" {
			command += q.Tables[k][0] + " , "
		} else {
			command += q.Tables[k][0] + " as " + q.Tables[k][1] + " , "
		}
	}
	return command[:len(command)-3]
}
func (q *FromQuery) QueryArgs() []interface{} {
	return []interface{}{}
}
