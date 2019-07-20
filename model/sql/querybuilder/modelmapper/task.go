package modelmapper

import "database/sql"
import "github.com/herb-go/herb/model/sql/querybuilder"

type Task interface {
	SetDB(querybuilder.DB)
	DB() querybuilder.DB
	OnSuccess(func() error)
	EmitSuccess() error
	OnPrepare(func() error)
	EmitPrepare() error
}

type CommonTask struct {
	db              querybuilder.DB
	successHandlers []func() error
	prepareHandlers []func() error
}

func (t *CommonTask) SetDB(db querybuilder.DB) {
	t.db = db
}

func (t *CommonTask) DB() querybuilder.DB {
	return t.db
}

func (t *CommonTask) OnSuccess(f func() error) {
	t.successHandlers = append(t.successHandlers, f)
}

func (t *CommonTask) EmitSuccess() error {
	var err error
	for k := range t.successHandlers {
		err = t.successHandlers[k]()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *CommonTask) OnPrepare(f func() error) {
	t.prepareHandlers = append(t.prepareHandlers, f)
}

func (t *CommonTask) EmitPrepare() error {
	var err error
	for k := range t.prepareHandlers {
		err = t.prepareHandlers[k]()
		if err != nil {
			return err
		}
	}
	return nil
}

type InsertTask struct {
	*querybuilder.InsertQuery
	CommonTask
}

func (t *InsertTask) Exec() (sql.Result, error) {
	err := t.EmitPrepare()
	if err != nil {
		return nil, err
	}
	r, err := t.InsertQuery.Query().Exec(t.db)
	if err != nil {
		return r, err
	}
	return r, t.EmitSuccess()
}

func NewInsertTask(q *querybuilder.InsertQuery, db querybuilder.DB) *InsertTask {
	t := &InsertTask{
		InsertQuery: q,
	}
	t.SetDB(db)
	return t
}

type UpdateTask struct {
	*querybuilder.UpdateQuery
	CommonTask
}

func (t *UpdateTask) Exec() (sql.Result, error) {
	err := t.EmitPrepare()
	if err != nil {
		return nil, err
	}
	r, err := t.UpdateQuery.Query().Exec(t.db)
	if err != nil {
		return r, err
	}
	return r, t.EmitSuccess()
}

func NewUpdateTask(q *querybuilder.UpdateQuery, db querybuilder.DB) *UpdateTask {
	t := &UpdateTask{
		UpdateQuery: q,
	}
	t.SetDB(db)
	return t
}

type DeleteTask struct {
	*querybuilder.DeleteQuery
	CommonTask
}

func (t *DeleteTask) Exec() (sql.Result, error) {
	err := t.EmitPrepare()
	if err != nil {
		return nil, err
	}
	r, err := t.DeleteQuery.Query().Exec(t.db)
	if err != nil {
		return r, err
	}
	return r, t.EmitSuccess()
}

func NewDeleteTask(q *querybuilder.DeleteQuery, db querybuilder.DB) *DeleteTask {
	t := &DeleteTask{
		DeleteQuery: q,
	}
	t.SetDB(db)
	return t
}

type SelectTask struct {
	*querybuilder.SelectQuery
	CommonTask
}

func (t *SelectTask) QueryRow() *sql.Row {
	return t.SelectQuery.Query().QueryRow(t.db)
}
func (t *SelectTask) QueryRows() (*sql.Rows, error) {
	return t.SelectQuery.Query().QueryRows(t.db)
}

func NewSelectTask(q *querybuilder.SelectQuery, db querybuilder.DB) *SelectTask {
	t := &SelectTask{
		SelectQuery: q,
	}
	t.SetDB(db)
	return t
}

func (t *SelectTask) ByField(fieldName string, fieldValue interface{}) *SelectTask {
	t.Where.Condition = t.Builder.Equal(fieldName, fieldValue)
	return t
}

func (t *SelectTask) ByFields(fieldsmap map[string]interface{}) *SelectTask {
	for k, v := range fieldsmap {
		t.Where.Condition.And(t.Builder.Equal(k, v))
	}
	return t
}

func (t *SelectTask) QueryRowToFields(fields *querybuilder.Fields) error {
	row := t.QueryRow()
	err := t.Select.Result().
		BindFields(fields).
		ScanFrom(row)
	return err
}
func (t *SelectTask) FindAllTo(rs ...Result) error {
	r := Results(rs)
	rows, err := t.QueryRows()
	if err != nil {
		return r.OnFinish(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = t.Select.Result().
			BindFields(r.Fields()).
			ScanFrom(rows)
		if err != nil {
			return r.OnFinish(err)
		}
		err = r.OnFinish(err)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *SelectTask) FindTo(rs ...Result) error {
	r := Results(rs)
	row := t.QueryRow()
	var sr = t.Select.Result()
	sr.BindFields(r.Fields())
	err := sr.ScanFrom(row)
	return r.OnFinish(err)
}

type Results []Result

func (rs *Results) Append(r ...Result) {
	*rs = append(*rs, r...)
}
func (rs *Results) Fields() *querybuilder.Fields {
	fields := querybuilder.NewFields()
	for _, v := range *rs {
		*fields = append(*fields, *v.Fields()...)
	}
	return fields
}
func (rs *Results) OnFinish(err error) error {
	for _, v := range *rs {
		err = v.OnFinish(err)
	}
	return err
}

type Result interface {
	Fields() *querybuilder.Fields
	OnFinish(error) error
}
