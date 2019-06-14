package querybuilder

// HavingClause having clause struct
type HavingClause struct {
	Buidler   *Builder
	Condition *PlainQuery
}

//NewHavingClause create new having clause
func (b *Builder) NewHavingClause() *HavingClause {
	return &HavingClause{
		Buidler:   b,
		Condition: b.New(""),
	}
}

// QueryCommand return query command
func (q *HavingClause) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "Having " + command
	}
	return command
}

// QueryArgs return query args
func (q *HavingClause) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
