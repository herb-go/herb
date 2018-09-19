package captcha

import "github.com/herb-go/herb/cache"

type Config struct {
	Enabled        bool
	Driver         string
	DisabledScenes map[string]bool
	Config         cache.ConfigMap
}

func (c *Config) ApplyTo(captcha *Captcha) error {
	d, err := NewDriver(c.Driver, &c.Config, "")
	if err != nil {
		return err
	}
	captcha.driver = d
	captcha.Enabled = c.Enabled
	captcha.DisabledScenes = c.DisabledScenes
	return nil
}
