package querybuilder

// NewUpdateClause create new update clause with given table name.
func (b *Builder) NewUpdateClause(tableName string) *UpdateClause {
	return &UpdateClause{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

// UpdateClause update clause struct
type UpdateClause struct {
	Builder *Builder
	// Prefix prefix query
	Prefix *PlainQuery
	// TableName table name
	TableName string
	// Data update query data list
	Data []QueryData
}

// AddSelect add subquery selecto with given field name to update clause
func (q *UpdateClause) AddSelect(field string, Select *SelectQuery) *UpdateClause {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}

// AddFields add fields to update clause
func (q *UpdateClause) AddFields(m *Fields) *UpdateClause {
	for _, v := range *m {
		q.Add(v.Field, v.Data)
	}
	return q
}

// Add add data with given field name to update clause
func (q *UpdateClause) Add(field string, data interface{}) *UpdateClause {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		},
	)
	return q
}

// AddRaw add raw data to given field
// Raw data will not be esaped.
// Dont add unsafe data by this method.
func (q *UpdateClause) AddRaw(field string, raw string) *UpdateClause {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}

// QueryCommand return query command
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

// QueryArgs return query args
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

// NewUpdateQuery create new update query with given table name.
func (b *Builder) NewUpdateQuery(tableName string) *UpdateQuery {
	return &UpdateQuery{
		Builder: b,
		Update:  b.NewUpdateClause(tableName),
		Where:   b.NewWhereClause(),
		Other:   b.New(""),
	}
}

// UpdateQuery update query struct
type UpdateQuery struct {
	Builder *Builder
	Update  *UpdateClause
	Where   *WhereClause
	Other   *PlainQuery
}

// Query convert update query to plain query.
func (u *UpdateQuery) Query() *PlainQuery {
	return u.Builder.Lines(u.Update, u.Where, u.Other)
}

// QueryCommand return query command
func (u *UpdateQuery) QueryCommand() string {
	return u.Query().Command
}

// QueryArgs return query args
func (u *UpdateQuery) QueryArgs() []interface{} {
	return u.Query().Args
}
