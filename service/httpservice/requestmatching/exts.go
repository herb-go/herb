package requestmatching

import (
	"net/http"
	"path/filepath"
)

type Exts map[string]bool

func (e Exts) MatchRequest(r *http.Request) (bool, error) {
	if len(e) == 0 {
		return true, nil
	}
	return e[filepath.Ext(r.URL.Path)], nil
}
