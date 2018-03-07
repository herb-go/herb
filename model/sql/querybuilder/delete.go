package querybuilder

type DeleteQuery struct {
	Builder   *Builder
	TableName string
	Prefix    *PlainQuery
	Alias     string
}

func (q *DeleteQuery) QueryCommand() string {
	var command = "DELETE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	if q.Alias != "" {
		command += " " + q.Alias
	}
	command += " FROM " + q.TableName
	if q.Alias != "" {
		command += " AS " + q.Alias
	}
	return command
}

func (q *DeleteQuery) QueryArgs() []interface{} {
	return q.Prefix.QueryArgs()
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
	Where   *WhereQurey
	Other   *PlainQuery
}

func (d *Delete) Query() *PlainQuery {
	return d.Builder.Lines(d.Delete, d.Where, d.Other)
}
