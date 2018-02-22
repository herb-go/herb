package render

import "errors"
import "fmt"

//ErrViewNotExist error which raised when view not exist.
var ErrViewNotExist = errors.New("ErrorViewNotExist")

//ViewError view error struct.
type ViewError struct {
	ViewName string
	err      error
}

//Error view error message.
func (ve *ViewError) Error() string {
	return fmt.Sprintf("View %s error: %s", ve.ViewName, ve.err.Error())
}

//NewViewError create view error by view name and raw error.
func NewViewError(ViewName string, err error) *ViewError {
	return &ViewError{
		ViewName: ViewName,
		err:      err,
	}
}

//ErrRegisterFuncNotSupported raised when register func is not supported by engine.
var ErrRegisterFuncNotSupported = errors.New("render:error register func not supported")
