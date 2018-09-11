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
