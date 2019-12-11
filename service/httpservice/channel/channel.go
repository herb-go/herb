package channel

import (
	"net/http"
	"sync"
)

var channels = sync.Map{}

type Channel struct {
	handler http.Handler
}
