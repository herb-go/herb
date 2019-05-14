package querybuilder

func (b *Builder) NewInsertQuery(tableName string) *InsertQuery {
	return &InsertQuery{
		Builder:   b,
		Prefix:    b.New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type InsertQuery struct {
	Builder   *Builder
	Prefix    *PlainQuery
	TableName string
	Data      []QueryData
	Select    *Select
}

func (q *InsertQuery) WithSelect(s *Select) *InsertQuery {
	q.Select = s
	return q
}

func (q *InsertQuery) AddFields(m *Fields) *InsertQuery {
	for _, v := range *m {
		q.Add(v.Field, v.Data)
	}
	return q
}
func (q *InsertQuery) Add(field string, data interface{}) *InsertQuery {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		})
	return q
}
func (q *InsertQuery) AddRaw(field string, raw string) *InsertQuery {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *InsertQuery) AddSelect(field string, Select *Select) *InsertQuery {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}
func (q *InsertQuery) QueryCommand() string {
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
func (q *InsertQuery) QueryArgs() []interface{} {

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
		Insert:  b.NewInsertQuery(tableName),
		Other:   b.New(""),
	}
}

type Insert struct {
	Builder *Builder
	Insert  *InsertQuery
	Other   *PlainQuery
}

func (i *Insert) Query() *PlainQuery {
	return i.Builder.Lines(i.Insert, i.Other)
}
