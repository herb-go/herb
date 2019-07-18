package modelmapper

import "database/sql"
import "github.com/herb-go/herb/model/sql/querybuilder"

// DBTable database table interface
type DBTable interface {
	querybuilder.DB
	DB() *sql.DB
	TableName() string
	Alias() string
	SetAlias(string)
	Driver() string
}

// ModelMapper database table mapper
type ModelMapper struct {
	DBTable
}

//QueryBuilder return querybuilder of  table
func (m *ModelMapper) QueryBuilder() *querybuilder.Builder {
	b := querybuilder.New()
	b.Driver = m.DBTable.Driver()
	return b
}

//FieldAlias return field name with table alias.
func (m *ModelMapper) FieldAlias(field string) string {
	a := m.Alias()
	if a != "" {
		field = a + "." + field
	}
	return field
}

//NewSelectQuery : create  select query for table
func (m *ModelMapper) NewSelectQuery() *querybuilder.SelectQuery {
	Select := m.QueryBuilder().NewSelectQuery()
	alias := m.Alias()
	if alias != "" {
		Select.From.AddAlias(alias, m.TableName())
	} else {
		Select.From.Add(m.TableName())
	}
	return Select
}

func (m *ModelMapper) Select() *SelectTask {
	return NewSelectTask(m.NewSelectQuery(), m)
}

//NewInsertQuery : new insert query for table node
func (m *ModelMapper) NewInsertQuery() *querybuilder.InsertQuery {
	Insert := m.QueryBuilder().NewInsertQuery(m.TableName())
	return Insert

}

func (m *ModelMapper) Insert() *InsertTask {
	return NewInsertTask(m.NewInsertQuery(), m)
}

//NewUpdateQuery : new update query for table
func (m *ModelMapper) NewUpdateQuery() *querybuilder.UpdateQuery {
	Update := m.QueryBuilder().NewUpdateQuery(m.TableName())
	return Update
}

func (m *ModelMapper) Update() *UpdateTask {
	return NewUpdateTask(m.NewUpdateQuery(), m)
}

//NewDeleteQuery : build delete query for table node
func (m *ModelMapper) NewDeleteQuery() *querybuilder.DeleteQuery {
	Delete := m.QueryBuilder().NewDeleteQuery(m.TableName())
	return Delete
}

func (m *ModelMapper) Delete() *DeleteTask {
	return NewDeleteTask(m.NewDeleteQuery(), m)
}

//NewCountQuery : build count select query for table
func (m *ModelMapper) NewCountQuery() *querybuilder.SelectQuery {
	Select := m.NewSelectQuery()
	Select.Select.Add("count(*)")
	return Select
}

//Count : count  from table  by given select m.QueryBuilder().
func (m *ModelMapper) Count(Select *querybuilder.SelectQuery) (int, error) {
	var result int
	row := Select.QueryRow(m)
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// New create new query table with given database table
func New(dbtable DBTable) *ModelMapper {
	return &ModelMapper{
		DBTable: dbtable,
	}
}
