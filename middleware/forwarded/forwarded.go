package forwarded

import "net/http"

func defaultTokenFailedAction(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

//Middleware main middleware struct.
type Middleware struct {
	//Enabled if this middleware is enabled.
	Enabled bool
	//ForwardedForHeader request header name which stores real ip.
	//If set to empty string,this feature will be disabeld.
	ForwardedForHeader string
	//ForwardedHostHeader request header name which stores real host.
	//If set to empty string,this feature will be disabeld.
	ForwardedHostHeader string
	//ForwardedProtoHeader request header name which stores real proto.
	//If set to empty string,this feature will be disabeld.
	ForwardedProtoHeader string
	//ForwardedTokenHeader request header name which stores token.
	//If set to empty string,this feature will be disabeld.
	ForwardedTokenHeader string
	//ForwardedTokenValue value which request header must equal.
	ForwardedTokenValue string
	tokenFailedAction   http.HandlerFunc
}

//SetTokenFailedAction set action which will execute when token verification fail
func (m *Middleware) SetTokenFailedAction(action func(w http.ResponseWriter, r *http.Request)) {
	m.tokenFailedAction = action
}

//ServeMiddleware return middleware.
func (m *Middleware) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !m.Enabled {
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

func (m *Middleware) Warnings() []string {
	if m.Enabled && (m.ForwardedForHeader != "" ||
		m.ForwardedHostHeader != "" ||
		m.ForwardedProtoHeader != "") &&
		m.ForwardedTokenHeader == "" {
		return []string{"Forwarded middleware is running without available ForwardedTokenHeader Value."}
	}
	return nil
}
