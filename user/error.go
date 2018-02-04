package user

import "errors"

//ErrAccountBindExists error rasied when account exists when binding
var ErrAccountBindExists = errors.New("account binded exists")

//ErrAccountUnbindNotExists error rasied when account does not exist when unbinding
var ErrAccountUnbindNotExists = errors.New("account unbinded dose not exist")
