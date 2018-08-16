package user

import "errors"

//ErrAccountBindingExists error rasied when account exists when binding
var ErrAccountBindingExists = errors.New("account binding exists")

//ErrAccountUnbindingNotExists error rasied when account does not exist when unbinding
var ErrAccountUnbindingNotExists = errors.New("account unbinding dose not exist")
