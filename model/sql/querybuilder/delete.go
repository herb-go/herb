package querybuilder

type DeleteQuery struct {
	Builder   *Builder
	TableName string
	Prefix    *PlainQuery
	alias     string
}

func (q *DeleteQuery) QueryCommand() string {
	return q.Builder.LoadDriver().DeleteCommandBuilder(q)
}
func (q *DeleteQuery) SetAlias(alias string) *DeleteQuery {
	q.alias = alias
	return q
}
func (q *DeleteQuery) QueryArgs() []interface{} {
	return q.Builder.LoadDriver().DeleteArgBuilder(q)
}

func (b *Builder) NewDeleteQuery(tableName string) *DeleteQuery {
	return &DeleteQuery{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
	}
}

func (b *Builder) NewDelete(TableName string) *Delete {
	return &Delete{
		Builder: b,
		Delete:  b.NewDeleteQuery(TableName),
		Where:   b.NewWhereQuery(),
		Other:   b.New(""),
	}
}

type Delete struct {
	Builder *Builder
	Delete  *DeleteQuery
	Where   *WhereQuery
	Other   *PlainQuery
}

func (d *Delete) Query() *PlainQuery {
	return d.Builder.Lines(d.Delete, d.Where, d.Other)
}
