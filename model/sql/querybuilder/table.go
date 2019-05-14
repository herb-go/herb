package querybuilder

import "database/sql"

type DBTable interface {
	DB
	DB() *sql.DB
	TableName() string
	Alias() string
	SetAlias(string)
	Driver() string
}

type Table struct {
	DBTable
}

//QueryBuilder return querybuilder of  table
func (t *Table) QueryBuilder() *Builder {
	b := NewBuilder()
	b.Driver = t.DBTable.Driver()
	return b
}

//FieldAlias return field name with table alias.
func (t *Table) FieldAlias(field string) string {
	a := t.Alias()
	if a != "" {
		field = a + "." + field
	}
	return field
}

//NewSelect : create  select query for table
func (t *Table) NewSelect() *Select {
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
func (t *Table) NewInsert() *Insert {
	Insert := t.QueryBuilder().NewInsert(t.TableName())
	return Insert

}

//NewUpdate : new update query for table
func (t *Table) NewUpdate() *Update {
	Update := t.QueryBuilder().NewUpdate(t.TableName())
	return Update
}

//NewDelete : build delete query for table node
func (t *Table) NewDelete() *Delete {
	Delete := t.QueryBuilder().NewDelete(t.TableName())
	return Delete
}

//BuildCount : build count select query for table
func (t *Table) BuildCount() *Select {
	Select := t.NewSelect()
	Select.Select.Add("count(*)")
	return Select
}

//Count : count  from table  by given select t.QueryBuilder().
func (t *Table) Count(Select *Select) (int, error) {
	var result int
	row := Select.QueryRow(t)
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func NewTable(dbtable DBTable) *Table {
	return &Table{
		DBTable: dbtable,
	}
}
