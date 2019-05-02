package querybuilder_test

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/herb-go/herb/model/sql/db"
)

func TestMysql(t *testing.T) {
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = DB.Exec("truncate table testtable1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = DB.Exec("truncate table testtable2")
	if err != nil {
		t.Fatal(err)
	}
}
func init() {

}
