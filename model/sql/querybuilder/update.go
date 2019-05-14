package querybuilder

func (b *Builder) NewUpdateClause(tableName string) *UpdateClause {
	return &UpdateClause{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type UpdateClause struct {
	Builder   *Builder
	Prefix    *PlainQuery
	TableName string
	Data      []QueryData
}

func (q *UpdateClause) AddSelect(field string, Select *Select) *UpdateClause {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}

func (q *UpdateClause) AddFields(m *Fields) *UpdateClause {
	for _, v := range *m {
		q.Add(v.Field, v.Data)
	}
	return q
}

func (q *UpdateClause) Add(field string, data interface{}) *UpdateClause {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		},
	)
	return q
}
func (q *UpdateClause) AddRaw(field string, raw string) *UpdateClause {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *UpdateClause) QueryCommand() string {
	var command = "UPDATE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " " + q.TableName
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
	if len(q.Data) > 0 {
		values = values[:len(values)-3]
	}
	command += values
	return command
}
func (q *UpdateClause) QueryArgs() []interface{} {
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
		Update:  b.NewUpdateClause(tableName),
		Where:   b.NewWhereClause(),
		Other:   b.New(""),
	}
}

type Update struct {
	Builder *Builder
	Update  *UpdateClause
	Where   *WhereClause
	Other   *PlainQuery
}

func (u *Update) Query() *PlainQuery {
	return u.Builder.Lines(u.Update, u.Where, u.Other)
}

func (u *Update) QueryCommand() string {
	return u.Query().Command
}

func (u *Update) QueryArgs() []interface{} {
	return u.Query().Args
}
