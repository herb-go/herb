package identifier

type Config struct {
	Driver string
	Config func(v interface{}) error `config:", lazyload"`
}
