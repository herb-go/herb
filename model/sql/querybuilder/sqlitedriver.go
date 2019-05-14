package querybuilder

type SqliteBuilderDriver struct {
	EmptyBuilderDriver
}

var SqliteDriver = &SqliteBuilderDriver{}

func init() {
	RegisterDriver("sqlite3", SqliteDriver)
}
