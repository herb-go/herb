package httpinfo

type Formatter interface {
	Format([]byte) ([]byte, bool, error)
}
