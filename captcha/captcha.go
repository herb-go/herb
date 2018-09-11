package captcha

import (
	"encoding/json"
	"net/http"
)

type driver interface {
	Type() string
	Config(scene string, w http.ResponseWriter, r *http.Request) (json.RawMessage, error)
	Reset(scene string, w http.ResponseWriter, r *http.Request) (json.RawMessage, error)
	Verify(r *http.Request, token string) (bool, error)
}

type CommonOutput struct {
	Type    string
	Enabled bool
	Config  json.RawMessage
}

func New() *Captcha {
	return &Captcha{
		DisabledScenes: map[string]bool{},
	}
}
func defaultEnabledChecker(captcha *Captcha, scene string, w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, nil
}

type Captcha struct {
	driver         driver
	Disabled       bool
	DisabledScenes map[string]bool
	EnabledChecker func(captcha *Captcha, scene string, w http.ResponseWriter, r *http.Request) (bool, error)
}

func (c *Captcha) EnabledCheck(scene string, w http.ResponseWriter, r *http.Request) (bool, error) {
	if c.Disabled {
		return false, nil
	}
	return c.EnabledChecker(c, scene, w, r)
}

func (c *Captcha) ConfigJSONAction(scene string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enabled, err := c.EnabledCheck(scene, w, r)
		if err != nil {
			panic(err)
		}
		output := CommonOutput{
			Type:    c.driver.Type(),
			Enabled: enabled,
		}
		if enabled {
			config, err := c.driver.Config(scene, w, r)
			if err != nil {
				panic(err)
			}
			output.Config = config
		}
		outputJSON, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(outputJSON)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Captcha) ResetJSONAction(scene string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enabled, err := c.EnabledCheck(scene, w, r)
		if err != nil {
			panic(err)
		}
		output := CommonOutput{
			Type:    c.driver.Type(),
			Enabled: enabled,
		}
		if enabled {
			config, err := c.driver.Reset(scene, w, r)
			if err != nil {
				panic(err)
			}
			output.Config = config
		}
		outputJSON, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(outputJSON)
		if err != nil {
			panic(err)
		}
	}
}
func (c *Captcha) Verify(r *http.Request, scene string, token string) (bool, error) {
	return c.driver.Verify(r, token)
}
