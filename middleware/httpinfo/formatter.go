package httpinfo

type Formatter interface {
	Format([]byte) ([]byte, bool, error)
}

type FormatterFunc func([]byte) ([]byte, bool, error)

func (f FormatterFunc) Format(data []byte) ([]byte, bool, error) {
	return f(data)
}
