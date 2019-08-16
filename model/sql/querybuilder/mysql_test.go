package querybuilder_test

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/querybuilder"
)

func TestMysqlIsDuplicate(t *testing.T) {
	var err error

	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(MysqlConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	builder := querybuilder.New()
	builder.Driver = DB.Driver()

	truncatequery := builder.TruncateTableQuery("testtable1")
	truncatequery.MustExec(DB)

	_, err = builder.Exec(DB, builder.TruncateTableQuery("testtable2"))
	if err != nil {
		t.Fatal(err)
	}
	insertquery := builder.NewInsertQuery("testtable1")
	fields := querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = builder.NewInsertQuery("testtable1")
	fields = querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if !builder.IsDuplicate(err) {
		t.Fatal(err)
	}
}
func TestMysql(t *testing.T) {
	type Result struct {
		ID   string
		Body string
	}
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(MysqlConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	builder := querybuilder.New()
	builder.Driver = "mysql"

	truncatequery := builder.TruncateTableQuery("testtable1")
	truncatequery.MustExec(DB)

	_, err = builder.Exec(DB, builder.TruncateTableQuery("testtable2"))
	if err != nil {
		t.Fatal(err)
	}
	var count int

	fields := querybuilder.NewFields()
	fields.Set(builder.CountField(), &count)
	countquery := builder.NewSelectQuery()
	countquery.Select.AddFields(fields)
	countquery.From.Add("testtable1")
	r := countquery.Query().QueryRow(DB)
	err = countquery.Result().BindFields(fields).ScanFrom(r)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal(err)
	}
	insertquery := builder.NewInsertQuery("testtable1")
	fields = querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	insertquery.Other = builder.New("ON DUPLICATE KEY UPDATE body= ?", "testbodydup")
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	result := &Result{}
	fields = querybuilder.NewFields()
	fields.Set("id", &result.ID).
		Set("body", &result.Body)
	selectquery := builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.AddFields(fields)
	row := selectquery.QueryRow(DB)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid" && result.Body != "testbody" {
		t.Fatal(result)
	}
	insertquery = builder.NewInsertQuery("testtable1")
	fields = querybuilder.NewFields()
	fields.Set("id", "testid2").Set("body", "testbody2")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}

	fields = querybuilder.NewFields()
	fields.Set(builder.CountField(), &count)
	countquery = builder.NewSelectQuery()
	countquery.Select.AddFields(fields)
	countquery.From.Add("testtable1")
	r = countquery.Query().QueryRow(DB)
	err = countquery.Result().BindFields(fields).ScanFrom(r)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal(err)
	}

	fields = querybuilder.NewFields()
	fields.Set("id", nil).
		Set("body", nil)
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.AddFields(fields)
	selectquery.OrderBy.Add("id", false)
	rows, err := selectquery.QueryRows(DB)
	if err != nil {
		t.Fatal(err)
	}
	results := []*Result{}
	for rows.Next() {
		var result = &Result{}
		fields = querybuilder.NewFields()
		fields.Set("id", &result.ID).
			Set("body", &result.Body)
		err = selectquery.Result().BindFields(fields).ScanFrom(rows)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, result)
	}
	if len(results) != 2 {
		t.Fatal(result)
	}
	if results[0].ID != "testid2" || results[0].Body != "testbody2" {
		t.Fatal(*results[0], *results[1])
	}
	if results[1].ID != "testid" || results[1].Body != "testbody" {
		t.Fatal(*results[0], *results[1])
	}
	//limit
	fields.Set("id", nil).
		Set("body", nil)
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.AddFields(fields)
	selectquery.OrderBy.Add("id", false)
	selectquery.Limit.SetLimit(1)
	selectquery.Limit.SetOffset(1)
	rows, err = selectquery.Query().QueryRows(DB)
	if err != nil {
		t.Fatal(err)
	}
	results = []*Result{}
	for rows.Next() {
		var result = &Result{}
		fields = querybuilder.NewFields()
		fields.Set("id", &result.ID).
			Set("body", &result.Body)
		err = selectquery.Result().BindFields(fields).ScanFrom(rows)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, result)
	}
	if len(results) != 1 {
		t.Fatal(result)
	}
	if results[0].ID != "testid" || results[0].Body != "testbody" {
		t.Fatal(*results[0])
	}
	fields = querybuilder.NewFields()
	fields.Set("id", "testid2").Set("body", "testbody2updated")
	updatequery := builder.NewUpdateQuery("testtable1")
	updatequery.Update.AddFields(fields)
	updatequery.Where.Condition = builder.Equal("id", "testid2")
	_, err = updatequery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	result = &Result{}
	fields.Set("id", &result.ID).Set("body", &result.Body)
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.AddFields(fields)
	selectquery.Where.Condition = builder.Equal("id", "testid2")
	row = selectquery.QueryRow(DB)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid2" || result.Body != "testbody2updated" {
		t.Fatal(result)
	}

	deletequery := builder.NewDeleteQuery("testtable1")
	deletequery.Where.Condition = builder.Equal("id", "testid2")
	_, err = deletequery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}

	fields.Set("id", nil).
		Set("body", nil)
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.AddFields(fields)
	rows, err = selectquery.QueryRows(DB)
	if err != nil {
		t.Fatal(err)
	}
	results = []*Result{}
	for rows.Next() {
		var result = &Result{}
		fields = querybuilder.NewFields()
		fields.Set("id", &result.ID).
			Set("body", &result.Body)
		err = selectquery.Result().BindFields(fields).ScanFrom(rows)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, result)
	}
	if len(results) != 1 {
		t.Fatal(result)
	}
	if results[0].ID != "testid" || results[0].Body != "testbody" {
		t.Fatal(*results[0])
	}
}
func TestJoin(t *testing.T) {
	type Result struct {
		ID    string
		Body  string
		Body2 string
	}
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(MysqlConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	builder := querybuilder.New()
	builder.Driver = "mysql"

	truncatequery := builder.TruncateTableQuery("testtable1")
	truncatequery.MustExec(DB)

	_, err = builder.Exec(DB, builder.TruncateTableQuery("testtable2"))
	if err != nil {
		t.Fatal(err)
	}

	insertquery := builder.NewInsertQuery("testtable1")
	fields := querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = builder.NewInsertQuery("testtable2")
	fields = querybuilder.NewFields()

	fields.Set("id", "testid").Set("body2", "testbody2")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	result := &Result{}
	fields = querybuilder.NewFields()
	fields.Set("t1.id", &result.ID).Set("t1.body", &result.Body).Set("t2.body2", &result.Body2)
	selectquery := builder.NewSelectQuery()
	selectquery.From.AddAlias("t1", "testtable1")
	selectquery.Select.AddFields(fields)
	selectquery.Join.LeftJoin().Alias("t2", "testtable2").On(builder.New("t1.id = t2.id"))
	row := selectquery.QueryRow(DB)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid" || result.Body != "testbody" || result.Body2 != "testbody2" {
		t.Fatal(*result)
	}

	result = &Result{}
	fields = querybuilder.NewFields()
	fields.Set("t1.id", &result.ID).Set("t1.body", &result.Body).Set("t2.body2", &result.Body2)
	selectquery = builder.NewSelectQuery()
	selectquery.From.AddAlias("t1", "testtable1")
	selectquery.Select.AddFields(fields)
	selectquery.Join.InnerJoin().Alias("t2", "testtable2").On(builder.New("t1.id = t2.id"))
	row = selectquery.QueryRow(DB)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid" || result.Body != "testbody" || result.Body2 != "testbody2" {
		t.Fatal(*result)
	}

	result = &Result{}
	fields = querybuilder.NewFields()
	fields.Set("t1.id", &result.ID).Set("t1.body", &result.Body).Set("t2.body2", &result.Body2)
	selectquery = builder.NewSelectQuery()
	selectquery.From.AddAlias("t1", "testtable1")
	selectquery.Select.AddFields(fields)
	selectquery.Join.RightJoin().Alias("t2", "testtable2").On(builder.New("t1.id = t2.id"))
	row = selectquery.QueryRow(DB)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid" || result.Body != "testbody" || result.Body2 != "testbody2" {
		t.Fatal(*result)
	}
}

func TestSubquery(t *testing.T) {
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(MysqlConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	builder := querybuilder.New()
	builder.Driver = "mysql"
	truncatequery := builder.TruncateTableQuery("testtable1")
	truncatequery.MustExec(DB)

	_, err = builder.Exec(DB, builder.TruncateTableQuery("testtable2"))
	if err != nil {
		t.Fatal(err)
	}
	insertquery := builder.NewInsertQuery("testtable1")
	fields := querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = builder.NewInsertQuery("testtable1")
	fields = querybuilder.NewFields()
	fields.Set("id", "testid2").Set("body", "testbody2")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = builder.NewInsertQuery("testtable2")
	insertquery.Insert.Add("id", nil).Add("body2", nil)
	selectquery := builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.AddRaw("raw").Add("testtable1.body")
	selectquery.Where.Condition = builder.Equal("testtable1.id", "testid")
	insertquery.Insert.WithSelect(selectquery)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	var body string
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable2")
	selectquery.Select.Add(("testtable2.body2"))
	selectquery.Where.Condition = builder.Equal("testtable2.id", "raw")
	row := selectquery.QueryRow(DB)
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody" {
		t.Fatal(body)
	}
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.Add("testtable1.body")
	selectquery.Where.Condition = builder.Equal("testtable1.id", "testid")
	insertquery = builder.NewInsertQuery("testtable2")
	insertquery.Insert.Add("id", "subquery").AddSelect("body2", selectquery)
	_, err = insertquery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable2")
	selectquery.Select.Add("testtable2.body2")
	selectquery.Where.Condition = builder.Equal("testtable2.id", "subquery")
	row = selectquery.QueryRow(DB)
	body = ""
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody" {
		t.Fatal(body)
	}
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable1")
	selectquery.Select.Add("body")
	selectquery.Where.Condition = builder.Equal("testtable1.id", "testid2")
	updatequery := builder.NewUpdateQuery("testtable2")
	updatequery.Where.Condition = builder.Equal("id", "subquery")
	updatequery.Update.AddSelect("body2", selectquery)
	_, err = updatequery.Query().Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
	selectquery = builder.NewSelectQuery()
	selectquery.From.Add("testtable2")
	selectquery.Select.Add("testtable2.body2")
	selectquery.Where.Condition = builder.Equal("testtable2.id", "subquery")
	row = selectquery.QueryRow(DB)
	body = ""
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody2" {
		t.Fatal(body)
	}
	selectquery = builder.NewSelectQuery()
	selectquery.From.AddAlias("t1", "testtable1")
	selectquery.Select.Add("t1.body")
	selectquery.Where.Condition = builder.New("t1.body=t2.body2")
	selectwithsubqueryquery := builder.NewSelectQuery()
	selectwithsubqueryquery.From.AddAlias("t2", "testtable2")
	selectwithsubqueryquery.Select.AddSelect(selectquery)
	selectwithsubqueryquery.Where.Condition = builder.Equal("t2.id", "subquery")
	row = selectwithsubqueryquery.QueryRow(DB)
	body = ""
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody2" {
		t.Fatal(body)
	}

}
func init() {

}
