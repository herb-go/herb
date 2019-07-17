package querybuilder

import "database/sql"

type Task interface {
	SetDB(DB)
	DB() DB
	OnSuccess(func() error)
	EmitSuccess() error
	OnPrepare(func() error)
	EmitPrepare() error
}

type CommonTask struct {
	db              DB
	successHandlers []func() error
	prepareHandlers []func() error
}

func (t *CommonTask) SetDB(db DB) {
	t.db = db
}

func (t *CommonTask) DB() DB {
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
	*Insert
	CommonTask
}

func (t *InsertTask) Exec() (sql.Result, error) {
	err := t.EmitPrepare()
	if err != nil {
		return nil, err
	}
	r, err := t.Insert.Query().Exec(t.db)
	if err != nil {
		return r, err
	}
	return r, t.EmitSuccess()
}

func NewInsertTask(q *Insert, db DB) *InsertTask {
	t := &InsertTask{
		Insert: q,
	}
	t.SetDB(db)
	return t
}

type UpdateTask struct {
	*Update
	CommonTask
}

func (t *UpdateTask) Exec() (sql.Result, error) {
	err := t.EmitPrepare()
	if err != nil {
		return nil, err
	}
	r, err := t.Update.Query().Exec(t.db)
	if err != nil {
		return r, err
	}
	return r, t.EmitSuccess()
}

func NewUpdateTask(q *Update, db DB) *UpdateTask {
	t := &UpdateTask{
		Update: q,
	}
	t.SetDB(db)
	return t
}

type DeleteTask struct {
	*Delete
	CommonTask
}

func (t *DeleteTask) Exec() (sql.Result, error) {
	err := t.EmitPrepare()
	if err != nil {
		return nil, err
	}
	r, err := t.Delete.Query().Exec(t.db)
	if err != nil {
		return r, err
	}
	return r, t.EmitSuccess()
}

func NewDeleteTask(q *Delete, db DB) *DeleteTask {
	t := &DeleteTask{
		Delete: q,
	}
	t.SetDB(db)
	return t
}

type SelectTask struct {
	*Select
	CommonTask
}

func (t *SelectTask) QueryRow() *sql.Row {
	return t.Select.Query().QueryRow(t.db)
}
func (t *SelectTask) QueryRows() (*sql.Rows, error) {
	return t.Select.Query().QueryRows(t.db)
}

func NewSelectTask(q *Select, db DB) *SelectTask {
	t := &SelectTask{
		Select: q,
	}
	t.SetDB(db)
	return t
}
