package builder

func (b *Builder) NewUpdateQuery(tableName string) *UpdateQuery {
	return &UpdateQuery{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type UpdateQuery struct {
	Builder   *Builder
	Prefix    *PlainQuery
	TableName string
	Alias     string
	Data      []QueryData
}

func (q *UpdateQuery) SetAlias(alias string) *UpdateQuery {
	q.Alias = alias
	return q
}
func (q *UpdateQuery) AddSelect(field string, Select *Select) *UpdateQuery {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}

func (q *UpdateQuery) AddFields(m Fields) *UpdateQuery {
	for k, v := range m {
		q.Add(k, v)
	}
	return q
}

func (q *UpdateQuery) Add(field string, data interface{}) *UpdateQuery {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		},
	)
	return q
}
func (q *UpdateQuery) AddRaw(field string, raw string) *UpdateQuery {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *UpdateQuery) QueryCommand() string {
	var command = "UPDATE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " " + q.TableName
	if q.Alias != "" {
		command += " AS " + q.Alias
	}
	command += " SET "
	var values = ""
	for k := range q.Data {
		values += q.Data[k].Field + " = "
		if q.Data[k].Raw == "" {
			values += "? , "
		} else {
			values += q.Data[k].Raw + " , "
		}
	}
	command += values[:len(values)-3]
	return command
}
func (q *UpdateQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		if q.Data[k].Data != nil {
			args = append(args, q.Data[k].Data...)
		}
	}
	var result = []interface{}{}
	result = append(result, q.Prefix.QueryArgs()...)
	result = append(result, args...)
	return result
}

func (b *Builder) NewUpdate(tableName string) *Update {
	return &Update{
		Builder: b,
		Update:  b.NewUpdateQuery(tableName),
		Where:   b.NewWhereQuery(),
		Other:   b.New(""),
	}
}

func (b *Builder) NewTableUpdate(t Table) *Update {
	update := b.NewUpdate(t.TableName())
	update.Update.SetAlias(t.Alias())
	return update
}

type Update struct {
	Builder *Builder
	Update  *UpdateQuery
	Where   *WhereQurey
	Other   *PlainQuery
}

func (u *Update) Query() *PlainQuery {
	return u.querybuilder.Lines(u.Update, u.Where, u.Other)
}
