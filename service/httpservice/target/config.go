package target

type ClientConfig struct {
	ClientDriver string
	Client       func(v interface{}) error `config:", lazyload"`
}

type Server struct {
	URLTarget
	ClientConfig
}
