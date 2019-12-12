package channel

import (
	"net/http"
)

type Handler struct {
	Stoped  bool
	handler http.Handler
}
