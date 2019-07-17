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

// Mapper database table mapper
type Mapper struct {
	DBTable
}

//QueryBuilder return querybuilder of  table
func (t *Mapper) QueryBuilder() *querybuilder.Builder {
	b := querybuilder.New()
	b.Driver = t.DBTable.Driver()
	return b
}

//FieldAlias return field name with table alias.
func (t *Mapper) FieldAlias(field string) string {
	a := t.Alias()
	if a != "" {
		field = a + "." + field
	}
	return field
}

//NewSelect : create  select query for table
func (t *Mapper) NewSelect() *querybuilder.Select {
	Select := t.QueryBuilder().NewSelect()
	alias := t.Alias()
	if alias != "" {
		Select.From.AddAlias(alias, t.TableName())
	} else {
		Select.From.Add(t.TableName())
	}
	return Select
}

//NewInsert : new insert query for table node
func (t *Mapper) NewInsert() *querybuilder.Insert {
	Insert := t.QueryBuilder().NewInsert(t.TableName())
	return Insert

}

//NewUpdate : new update query for table
func (t *Mapper) NewUpdate() *querybuilder.Update {
	Update := t.QueryBuilder().NewUpdate(t.TableName())
	return Update
}

//NewDelete : build delete query for table node
func (t *Mapper) NewDelete() *querybuilder.Delete {
	Delete := t.QueryBuilder().NewDelete(t.TableName())
	return Delete
}

//BuildCount : build count select query for table
func (t *Mapper) BuildCount() *querybuilder.Select {
	Select := t.NewSelect()
	Select.Select.Add("count(*)")
	return Select
}

//Count : count  from table  by given select t.QueryBuilder().
func (t *Mapper) Count(Select *querybuilder.Select) (int, error) {
	var result int
	row := Select.QueryRow(t)
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// New create new query table with given database table
func New(dbtable DBTable) *Mapper {
	return &Mapper{
		DBTable: dbtable,
	}
}
