package querybuilder

func (b *Builder) NewInsertClause(tableName string) *InsertClause {
	return &InsertClause{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type InsertClause struct {
	Builder   *Builder
	Prefix    *PlainQuery
	TableName string
	Data      []QueryData
	Select    *Select
}

func (q *InsertClause) WithSelect(s *Select) *InsertClause {
	q.Select = s
	return q
}

func (q *InsertClause) AddFields(m *Fields) *InsertClause {
	for _, v := range *m {
		q.Add(v.Field, v.Data)
	}
	return q
}
func (q *InsertClause) Add(field string, data interface{}) *InsertClause {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		})
	return q
}
func (q *InsertClause) AddRaw(field string, raw string) *InsertClause {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *InsertClause) AddSelect(field string, Select *Select) *InsertClause {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}
func (q *InsertClause) QueryCommand() string {
	var command = "INSERT"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " INTO " + q.TableName
	var values = ""
	var columns = ""
	for k := range q.Data {
		if q.Data[k].Raw == "" {
			values += "? , "
		} else {
			values += q.Data[k].Raw + " , "
		}
		columns += q.Data[k].Field + " , "
	}
	if len(q.Data) > 0 {
		columns = columns[:len(columns)-3]
		values = values[:len(values)-3]
	}
	command += " ("
	command += columns
	command += " )"
	if q.Select == nil {
		command += " VALUES ( "
		command += values
		command += " )"
	} else {
		command += "\n"
		command += q.Select.QueryCommand()
	}
	return command
}
func (q *InsertClause) QueryArgs() []interface{} {

	var result = []interface{}{}
	result = append(result, q.Prefix.QueryArgs()...)
	if q.Select == nil {
		var args = []interface{}{}
		for k := range q.Data {
			if q.Data[k].Data != nil {
				args = append(args, q.Data[k].Data...)
			}
		}
		result = append(result, args...)
	} else {
		args := q.Select.QueryArgs()
		result = append(result, args...)
	}
	return result
}

func (b *Builder) NewInsert(tableName string) *Insert {
	return &Insert{
		Builder: b,
		Insert:  b.NewInsertClause(tableName),
		Other:   b.New(""),
	}
}

type Insert struct {
	Builder *Builder
	Insert  *InsertClause
	Other   *PlainQuery
}

func (i *Insert) Query() *PlainQuery {
	return i.Builder.Lines(i.Insert, i.Other)
}
