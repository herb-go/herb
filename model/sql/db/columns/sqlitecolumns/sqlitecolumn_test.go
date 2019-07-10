package sqlitecolumns

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/db/columns"
)

var target = []*columns.Column{
	&columns.Column{Field: "id", ColumnType: "int", AutoValue: true, PrimayKey: true, NotNull: true},
	&columns.Column{Field: "f_nullable", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: false},
	&columns.Column{Field: "f_bool", ColumnType: "byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_smallint", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_mediumint", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_int", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_integer", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_tinyint", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_int2", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_int8", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_biginteger", ColumnType: "int64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_float", ColumnType: "float32", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_double", ColumnType: "float64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_doubleprecision", ColumnType: "float64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_real", ColumnType: "float64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_datetime", ColumnType: "time.Time", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_char", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_varchar", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_character", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_nvarchar", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_nchar", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_text", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_blob", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
}

func TestColumns(t *testing.T) {
	config := db.NewConfig()
	err := json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	db := db.New()
	err = db.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	sql, err := ioutil.ReadFile("./sql/sqlite.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		t.Fatal(err)
	}
	defer db.DB().Close()
	c, err := columns.Driver("sqlite3")
	if err != nil {
		t.Fatal(err)
	}
	err = c.Load(db, "columns")
	if err != nil {
		t.Fatal(err)
	}
	columns, err := c.Columns()
	if err != nil {
		t.Fatal(err)
	}

	if len(columns) != len(target) {
		t.Fatal(columns)
	}

	for k := range columns {
		if columns[k].Field != target[k].Field ||
			columns[k].ColumnType != target[k].ColumnType ||
			columns[k].AutoValue != target[k].AutoValue ||
			columns[k].PrimayKey != target[k].PrimayKey ||
			columns[k].NotNull != target[k].NotNull {
			t.Fatal(columns[k])
		}
	}
}

func TestConvertType(t *testing.T) {
	_, err := ConvertType("unknowtype")
	if err == nil || !strings.Contains(err.Error(), "unknowtype") {
		t.Fatal(err)
	}
}
