package forwarded

import "net/http"

const StatusDisabled = 0

func defaultTokenFailedAction(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

type Middleware struct {
	Status               int
	ForwardedForHeader   string
	ForwardedHostHeader  string
	ForwardedProtoHeader string
	ForwardedTokenHeader string
	ForwardedTokenValue  string
	tokenFailedAction    http.HandlerFunc
}

func (m *Middleware) SetTokenFailedAction(action func(w http.ResponseWriter, r *http.Request)) {
	m.tokenFailedAction = action
}
func (m *Middleware) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if m.Status == StatusDisabled {
		next(w, r)
		return
	}
	if m.ForwardedTokenHeader != "" {
		if m.ForwardedTokenValue == "" || r.Header.Get(m.ForwardedTokenHeader) != m.ForwardedTokenValue {
			action := m.tokenFailedAction
			if action == nil {
				action = defaultTokenFailedAction
			}
			action(w, r)
			return
		}
	}
	var headers = r.Header
	if m.ForwardedForHeader != "" {
		forwardedFor := headers.Get(m.ForwardedForHeader)
		if forwardedFor != "" {
			r.RemoteAddr = forwardedFor + ":-1"
		}
	}
	if m.ForwardedProtoHeader != "" {
		forwardedProto := headers.Get(m.ForwardedProtoHeader)
		if forwardedProto != "" {
			r.URL.Scheme = forwardedProto
		}
	}
	if m.ForwardedHostHeader != "" {
		forwardedHost := headers.Get(m.ForwardedHostHeader)
		if forwardedHost != "" {
			r.Host = forwardedHost
		}
	}
	next(w, r)
}
