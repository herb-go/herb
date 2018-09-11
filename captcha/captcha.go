package captcha

import (
	"encoding/json"
	"net/http"
)

type driver interface {
	Type() string
	Config(w http.ResponseWriter, r *http.Request) (json.RawMessage, error)
	Reset(w http.ResponseWriter, r *http.Request) (json.RawMessage, error)
	Verify(r *http.Request, token string) (bool, error)
}

type CommonOutput struct {
	Type    string
	Enabled bool
	Config  json.RawMessage
}

func New() *Captcha {
	return &Captcha{}
}
func defaultEnabledChecker(captcha *Captcha, w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, nil
}

type Captcha struct {
	driver         driver
	Disabled       bool
	EnabledChecker func(captcha *Captcha, w http.ResponseWriter, r *http.Request) (bool, error)
}

func (c *Captcha) EnabledCheck(w http.ResponseWriter, r *http.Request) (bool, error) {
	if c.Disabled {
		return false, nil
	}
	return c.EnabledChecker(c, w, r)
}

func (c *Captcha) ConfigAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enabled, err := c.EnabledCheck(w, r)
	if err != nil {
		panic(err)
	}
	output := CommonOutput{
		Type:    c.driver.Type(),
		Enabled: enabled,
	}
	if enabled {
		config, err := c.driver.Config(w, r)
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

func (c *Captcha) ResetAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enabled, err := c.EnabledCheck(w, r)
	if err != nil {
		panic(err)
	}
	output := CommonOutput{
		Type:    c.driver.Type(),
		Enabled: enabled,
	}
	if enabled {
		config, err := c.driver.Reset(w, r)
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

func (c *Captcha) Verify(r *http.Request, token string) (bool, error) {
	return c.driver.Verify(r, token)
}
