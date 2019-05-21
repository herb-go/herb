package querybuilder

// SqliteBuilderDriver sqlite builder driver
type SqliteBuilderDriver struct {
	EmptyBuilderDriver
}

// SqliteDriver sqlite driver
var SqliteDriver = &SqliteBuilderDriver{}

func init() {
	RegisterDriver("sqlite3", SqliteDriver)
}
