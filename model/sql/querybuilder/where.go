package querybuilder

// NewWhereClause create new where clause
func (b *Builder) NewWhereClause() *WhereClause {
	return &WhereClause{
		Builder:   b,
		Condition: b.New(""),
	}
}

// WhereClause where clause struct
type WhereClause struct {
	Builder *Builder
	// Condition where condition
	Condition *PlainQuery
}

// QueryCommand return query command
func (q *WhereClause) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "WHERE " + command
	}
	return command
}

// QueryArgs return query args
func (q *WhereClause) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
