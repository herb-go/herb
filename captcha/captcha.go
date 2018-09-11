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
	Type   string
	Config json.RawMessage
}

func New() *Captcha {
	return &Captcha{}
}

type Captcha struct {
	driver driver
}

func (c *Captcha) ActionConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config, err := c.driver.Config(w, r)
	if err != nil {
		panic(err)
	}
	output, err := json.Marshal(CommonOutput{
		Type:   c.driver.Type(),
		Config: config,
	})
	if err != nil {
		panic(err)
	}
	_, err = w.Write(output)
	if err != nil {
		panic(err)
	}
}

func (c *Captcha) ActionReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config, err := c.driver.Reset(w, r)
	if err != nil {
		panic(err)
	}
	output, err := json.Marshal(CommonOutput{
		Type:   c.driver.Type(),
		Config: config,
	})
	if err != nil {
		panic(err)
	}
	_, err = w.Write(output)
	if err != nil {
		panic(err)
	}
}

func (c *Captcha) Verify(r *http.Request, token string) (bool, error) {
	return c.driver.Verify(r, token)
}
