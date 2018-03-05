package builder

func (b *Builder) NewWhereQuery() *WhereQurey {
	return &WhereQurey{
		Condition: b.New(""),
	}
}

type WhereQurey struct {
	Condition *PlainQuery
}

func (q *WhereQurey) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "WHERE " + command
	}
	return command
}
func (q *WhereQurey) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
