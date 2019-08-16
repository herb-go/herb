package querybuilder

import driver "github.com/mattn/go-sqlite3"

//SqliteBuilderDriver sqlite builder driver
type SqliteBuilderDriver struct {
	EmptyBuilderDriver
}

// SqliteDriver sqlite driver
var SqliteDriver = &SqliteBuilderDriver{}

//TruncateTableCommandBuilder return truncate table query.
func (d *SqliteBuilderDriver) TruncateTableCommandBuilder(t string) string {
	return "DELETE FROM " + t
}
func init() {
	RegisterDriver("sqlite3", SqliteDriver)
}

//IsDuplicate check if error is Is duplicate error.
func (d *SqliteBuilderDriver) IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(driver.Error)
	if ok == false {
		return false
	}
	return e.Code == 19
}
