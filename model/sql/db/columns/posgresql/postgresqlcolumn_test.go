package posgresqlcolumn

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/db/columns"
)

var target = []*columns.Column{
	&columns.Column{Field: "id", ColumnType: "int", AutoValue: true, PrimayKey: true, NotNull: true},
	&columns.Column{Field: "f_nullable", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: false},
	&columns.Column{Field: "f_bigint", ColumnType: "int64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_bigserial", ColumnType: "int64", AutoValue: true, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_boolean", ColumnType: "bool", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_bytea", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_char", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_varchar", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_float8", ColumnType: "float64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_integer", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_json", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_jsonb", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_real", ColumnType: "float32", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_text", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_timestamptz", ColumnType: "time.Time", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_timestamp", ColumnType: "time.Time", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_smallint", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_smallserial", ColumnType: "int", AutoValue: true, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_uuid", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_xml", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
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
	defer db.DB().Close()
	c, err := columns.Driver("postgres")
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
