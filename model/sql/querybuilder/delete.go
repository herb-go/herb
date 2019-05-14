package querybuilder

//DeleteQuery delete query  struct
type DeleteQuery struct {
	// Builder query builder which create this query.
	Builder *Builder
	//TableName database table name
	TableName string
	// Prefix  query which insert between"DELETE" command and  table name.
	Prefix *PlainQuery
	alias  string
}

// NewDeleteQuery create new delete statement with given table name.
func (b *Builder) NewDeleteQuery(tableName string) *DeleteQuery {
	return &DeleteQuery{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
	}
}

//QueryCommand return query command
func (q *DeleteQuery) QueryCommand() string {
	return q.Builder.LoadDriver().DeleteCommandBuilder(q)
}

// QueryArgs return query adts
func (q *DeleteQuery) QueryArgs() []interface{} {
	return q.Builder.LoadDriver().DeleteArgBuilder(q)
}

// NewDelete create new delete query with given table name.s
func (b *Builder) NewDelete(TableName string) *Delete {
	return &Delete{
		Builder: b,
		Delete:  b.NewDeleteQuery(TableName),
		Where:   b.NewWhereQuery(),
		Other:   b.New(""),
	}
}

// Delete delete query
type Delete struct {
	// Builder query builder which create this query.
	Builder *Builder
	// Delete delete query
	Delete *DeleteQuery
	// Where where query
	Where *WhereQuery
	// Other  query after where
	Other *PlainQuery
}

// Query return plain query
func (d *Delete) Query() *PlainQuery {
	return d.Builder.Lines(d.Delete, d.Where, d.Other)
}
