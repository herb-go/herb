package requestinforeader

type Config struct {
	Name   string
	Type   string
	Config func(v interface{}) error
}

func (c *Config) Register() error {
	f, err := GetFactory(c.Type)
	if err != nil {
		return err
	}
	reader, err := f.CreateReader(c.Config)
	if err != nil {
		return err
	}
	Register(c.Name, reader)
	return nil
}
