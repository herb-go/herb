package querybuilder

//DeleteClause delete clause  struct
type DeleteClause struct {
	// Builder query builder which create this query.
	Builder *Builder
	//TableName database table name
	TableName string
	// Prefix  query which insert between"DELETE" command and  table name.
	Prefix *PlainQuery
}

//NewDeleteClause create new delete clause with given table name.
func (b *Builder) NewDeleteClause(tableName string) *DeleteClause {
	return &DeleteClause{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
	}
}

//QueryCommand return query command
func (q *DeleteClause) QueryCommand() string {
	return q.Builder.LoadDriver().DeleteCommandBuilder(q)
}

// QueryArgs return query adts
func (q *DeleteClause) QueryArgs() []interface{} {
	return q.Builder.LoadDriver().DeleteArgBuilder(q)
}

// NewDelete create new delete query with given table name.s
func (b *Builder) NewDelete(TableName string) *Delete {
	return &Delete{
		Builder: b,
		Delete:  b.NewDeleteClause(TableName),
		Where:   b.NewWhereClause(),
		Other:   b.New(""),
	}
}

// Delete delete query
type Delete struct {
	// Builder query builder which create this query.
	Builder *Builder
	// Delete delete query
	Delete *DeleteClause
	// Where where query
	Where *WhereClause
	// Other  query after where
	Other *PlainQuery
}

// Query return plain query
func (d *Delete) Query() *PlainQuery {
	return d.Builder.Lines(d.Delete, d.Where, d.Other)
}
