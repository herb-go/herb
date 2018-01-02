package user

import "errors"

var ErrAccountBindExists = errors.New("account binded exists")
var ErrAccountUnbindNotExists = errors.New("account unbinded not exists")
