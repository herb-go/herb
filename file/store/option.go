package store

//Option store option interface.
type Option interface {
	ApplyTo(*Store) error
}

// OptionConfig option config in map format.
type OptionConfig struct {
	Driver string
	Config func(interface{}) error `config:", lazyload"`
}

//ApplyTo apply option to file store.
func (o *OptionConfig) ApplyTo(store *Store) error {
	driver, err := NewDriver(o.Driver, o.Config)
	if err != nil {
		return err
	}
	store.Driver = driver
	return nil
}

//NewOptionConfig create new option config.
func NewOptionConfig() *OptionConfig {
	return &OptionConfig{}
}
