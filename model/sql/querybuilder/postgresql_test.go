package querybuilder_test

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/querybuilder"
	_ "github.com/lib/pq"
)

func TestPostgresql(t *testing.T) {
	type Result struct {
		ID   string
		Body string
	}
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(PostgreConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	table1 := querybuilder.NewTable(DB.Table("testtable1"))

	truncatequery := table1.QueryBuilder().New("truncate table testtable1")
	truncatequery.MustExec(table1)

	_, err = table1.QueryBuilder().Exec(table1, table1.QueryBuilder().New("truncate table testtable2"))
	if err != nil {
		t.Fatal(err)
	}

	// builder := table1.QueryBuilder()
	fields := querybuilder.NewFields()
	var count int
	fields.Set(table1.QueryBuilder().CountField(), &count)
	countquery := table1.BuildCount()
	r := countquery.QueryRow(table1)
	err = countquery.Result().BindFields(fields).ScanFrom(r)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal(err)
	}
	insertquery := table1.NewInsert()
	fields = querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	result := &Result{}
	fields = querybuilder.NewFields()
	fields.Set("id", &result.ID).
		Set("body", &result.Body)
	selectquery := table1.NewSelect()
	selectquery.Select.AddFields(fields)
	row := selectquery.QueryRow(table1)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid" && result.Body != "testbody" {
		t.Fatal(result)
	}
	insertquery = table1.NewInsert()
	fields = querybuilder.NewFields()
	fields.Set("id", "testid2").Set("body", "testbody2")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	countquery = table1.BuildCount()
	c, err := table1.Count(countquery)
	if err != nil {
		t.Fatal(err)
	}
	if c != 2 {
		t.Fatal(c)
	}
	fields.Set("id", nil).
		Set("body", nil)
	selectquery = table1.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.OrderBy.Add("id", false)
	rows, err := selectquery.QueryRows(table1)
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
	selectquery = table1.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.OrderBy.Add("id", false)
	selectquery.Limit.SetLimit(1)
	selectquery.Limit.SetOffset(1)
	rows, err = selectquery.QueryRows(table1)
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
	updatequery := table1.NewUpdate()
	updatequery.Update.AddFields(fields)
	updatequery.Where.Condition = table1.QueryBuilder().Equal("id", "testid2")
	_, err = updatequery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	result = &Result{}
	fields.Set("id", &result.ID).Set("body", &result.Body)
	selectquery = table1.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.Where.Condition = table1.QueryBuilder().Equal("id", "testid2")
	row = selectquery.QueryRow(table1)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid2" || result.Body != "testbody2updated" {
		t.Fatal(result)
	}

	deletequery := table1.NewDelete()
	deletequery.Where.Condition = table1.QueryBuilder().Equal(table1.FieldAlias("id"), "testid2")
	_, err = deletequery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}

	fields.Set("id", nil).
		Set("body", nil)
	selectquery = table1.NewSelect()
	selectquery.Select.AddFields(fields)
	rows, err = selectquery.QueryRows(table1)
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
func TestPostgresqlJoin(t *testing.T) {
	type Result struct {
		ID    string
		Body  string
		Body2 string
	}
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(PostgreConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	table1 := querybuilder.NewTable(DB.Table("testtable1"))
	table1.SetAlias("t1")
	builder := table1.QueryBuilder()

	truncatequery := table1.QueryBuilder().New("truncate table testtable1")
	truncatequery.MustExec(table1)
	table2 := querybuilder.NewTable(DB.Table("testtable2"))

	_, err = DB.Exec("truncate table testtable2")
	if err != nil {
		t.Fatal(err)
	}
	insertquery := table1.NewInsert()
	fields := querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = table2.NewInsert()
	fields = querybuilder.NewFields()

	fields.Set("id", "testid").Set("body2", "testbody2")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	result := &Result{}
	fields = querybuilder.NewFields()
	fields.Set("t1.id", &result.ID).Set("t1.body", &result.Body).Set("t2.body2", &result.Body2)
	selectquery := table1.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.Join.LeftJoin().Alias("t2", table2.TableName()).On(builder.New("t1.id = t2.id"))
	row := selectquery.QueryRow(table1)
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
	selectquery = table1.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.Join.InnerJoin().Alias("t2", table2.TableName()).Using("id")
	row = selectquery.QueryRow(table1)
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
	selectquery = table1.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.Join.RightJoin().Alias("t2", table2.TableName()).Using("id")
	row = selectquery.QueryRow(table1)
	err = selectquery.Result().BindFields(fields).ScanFrom(row)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testid" || result.Body != "testbody" || result.Body2 != "testbody2" {
		t.Fatal(*result)
	}
}

func TestPostgresqlSubquery(t *testing.T) {
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(PostgreConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	table1 := querybuilder.NewTable(DB.Table("testtable1"))
	table1.SetAlias("t1")

	// builder := table1.QueryBuilder()
	truncatequery := table1.QueryBuilder().New("truncate table testtable1")
	truncatequery.MustExec(table1)
	table2 := querybuilder.NewTable(DB.Table("testtable2"))
	table2.SetAlias("t2")

	_, err = DB.Exec("truncate table testtable2")
	if err != nil {
		t.Fatal(err)
	}
	insertquery := table1.NewInsert()
	fields := querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = table1.NewInsert()
	fields = querybuilder.NewFields()
	fields.Set("id", "testid2").Set("body", "testbody2")
	insertquery.Insert.AddFields(fields)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	insertquery = table2.NewInsert()
	insertquery.Insert.Add("id", nil).Add("body2", nil)
	selectquery := table1.NewSelect()
	selectquery.Select.AddRaw("raw").Add(table1.FieldAlias("body"))
	selectquery.Where.Condition = table1.QueryBuilder().Equal(table1.FieldAlias("id"), "testid")
	insertquery.Insert.SetSelect(selectquery)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	var body string
	selectquery = table2.NewSelect()
	selectquery.Select.Add(table2.FieldAlias("body2"))
	selectquery.Where.Condition = table2.QueryBuilder().Equal(table2.FieldAlias("id"), "raw")
	row := selectquery.QueryRow(table2)
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody" {
		t.Fatal(body)
	}
	selectquery = table1.NewSelect()
	selectquery.Select.Add(table1.FieldAlias("body"))
	selectquery.Where.Condition = table1.QueryBuilder().Equal(table1.FieldAlias("id"), "testid")
	insertquery = table2.NewInsert()
	insertquery.Insert.Add("id", "subquery").AddSelect("body2", selectquery)
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	selectquery = table2.NewSelect()
	selectquery.Select.Add(table2.FieldAlias("body2"))
	selectquery.Where.Condition = table2.QueryBuilder().Equal(table2.FieldAlias("id"), "subquery")
	row = selectquery.QueryRow(table2)
	body = ""
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody" {
		t.Fatal(body)
	}
	selectquery = table1.NewSelect()
	selectquery.Select.Add(table1.FieldAlias("body"))
	selectquery.Where.Condition = table1.QueryBuilder().Equal(table1.FieldAlias("id"), "testid2")
	updatequery := table2.NewUpdate()
	updatequery.Update.SetAlias(table2.Alias())
	updatequery.Where.Condition = table2.QueryBuilder().Equal(table2.FieldAlias("id"), "subquery")
	updatequery.Update.AddSelect("body2", selectquery)
	_, err = updatequery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}
	selectquery = table2.NewSelect()
	selectquery.Select.Add(table2.FieldAlias("body2"))
	selectquery.Where.Condition = table2.QueryBuilder().Equal(table2.FieldAlias("id"), "subquery")
	row = selectquery.QueryRow(table2)
	body = ""
	err = row.Scan(&body)
	if err != nil {
		t.Fatal(err)
	}
	if body != "testbody2" {
		t.Fatal(body)
	}
	selectquery = table1.NewSelect()
	selectquery.Select.Add(table1.FieldAlias("body"))
	selectquery.Where.Condition = table1.QueryBuilder().New("t1.body=t2.body2")
	selectwithsubqueryquery := table2.NewSelect()
	selectwithsubqueryquery.Select.AddSelect(selectquery)
	selectwithsubqueryquery.Where.Condition = table1.QueryBuilder().Equal(table2.FieldAlias("id"), "subquery")
	row = selectwithsubqueryquery.QueryRow(table2)
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
