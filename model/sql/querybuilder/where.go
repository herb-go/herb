package querybuilder

func (b *Builder) NewWhereQuery() *WhereQuery {
	return &WhereQuery{
		Condition: b.New(""),
	}
}

type WhereQuery struct {
	Condition *PlainQuery
}

func (q *WhereQuery) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "WHERE " + command
	}
	return command
}
func (q *WhereQuery) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
