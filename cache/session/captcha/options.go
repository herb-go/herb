package captcha

//Config captcha config struct.
type Config struct {
	Enabled        bool
	Driver         string
	DisabledScenes map[string]bool
	AddrWhiteList  []string
	Config         func(interface{}) error `config:", lazyload"`
}

//ApplyTo apply config to captcha.
func (c *Config) ApplyTo(captcha *Captcha) error {
	d, err := NewDriver(c.Driver, c.Config)
	if err != nil {
		return err
	}
	captcha.driver = d
	captcha.Enabled = c.Enabled
	captcha.DisabledScenes = c.DisabledScenes
	captcha.AddrWhiteList = c.AddrWhiteList
	return nil
}
