package target

type ClientConfig struct {
	ClientDriver string
	Client       func(v interface{}) error `config:", lazyload"`
}

type Server struct {
	URLTarget
	ClientConfig
}

func (s *Server) CreatePlan() (Plan, error) {
	var err error
	p := NewPlan()
	p.Target = &s.URLTarget
	p.Doer, err = NewDoer(s.ClientDriver, s.Client)
	if err != nil {
		return nil, err
	}
	return p, nil
}
