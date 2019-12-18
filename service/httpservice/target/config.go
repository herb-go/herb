package target

type Config struct {
	ClientDriver string
	Client       func(v interface{}) error `config:", lazyload"`
	TargetDriver string
	Target       func(v interface{}) error `config:", lazyload"`
}
