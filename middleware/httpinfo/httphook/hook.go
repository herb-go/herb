package httphook

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
)

type Hook struct {
	handler         Handler
	validator       httpinfo.Validator
	bufferValidator httpinfo.Validator
}

func (hook *Hook) Clone() *Hook {
	return &Hook{
		handler:   hook.handler,
		validator: hook.validator,
	}
}
func (hook *Hook) WithHandler(h Handler) *Hook {
	newhook := hook.Clone()
	newhook.handler = h
	return newhook
}

func (hook *Hook) WithValidator(v httpinfo.Validator) *Hook {
	newhook := hook.Clone()
	newhook.validator = v
	return newhook
}

func (hook *Hook) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resp := httpinfo.NewResponse()
	writer := resp.WrapWriter(w)
	if hook.bufferValidator != nil {
		resp.BuildBuffer(r, hook.bufferValidator)
	}
	next(writer, r)
	ok, err := hook.validator.Validate(r, resp)
	if err != nil {
		panic(err)
	}
	if ok {
		hook.handler.Handle(r, resp)
	}
}

func New() *Hook {
	return roothook
}

var roothook = &Hook{
	handler:   NopHandler,
	validator: httpinfo.ValidatorNever,
}
