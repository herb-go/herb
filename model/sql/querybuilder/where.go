package querybuilder

func (b *Builder) NewWhereClause() *WhereClause {
	return &WhereClause{
		Condition: b.New(""),
	}
}

type WhereClause struct {
	Condition *PlainQuery
}

func (q *WhereClause) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "WHERE " + command
	}
	return command
}
func (q *WhereClause) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
