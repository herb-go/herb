package user

import "errors"

//ErrAccountBindingExists error rasied when account exists when binding
var ErrAccountBindingExists = errors.New("account binding exists")

//ErrAccountUnbindingNotExists error rasied when account does not exist when unbinding
var ErrAccountUnbindingNotExists = errors.New("account unbinding does not exist")

//ErrUserNotExists error raseid when user does not exist.
var ErrUserNotExists = errors.New("user does not exist")

//ErrUserExists error raseid when user does  exist.
var ErrUserExists = errors.New("user exists")
