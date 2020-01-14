package requestmatching

import "net/http"

type Methods map[string]bool

func (m Methods) MatchRequest(r *http.Request) (bool, error) {
	if len(m) == 0 {
		return true, nil
	}
	return m[r.Method], nil
}
