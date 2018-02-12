package render

import "errors"
import "fmt"

var ErrorViewNotExist = errors.New("ErrorViewNotExist")

type ViewError struct {
	ViewName string
	err      error
}

func (ve *ViewError) Error() string {
	return fmt.Sprintf("View %s error: %s", ve.ViewName, ve.err.Error())
}

func NewViewError(ViewName string, err error) *ViewError {
	return &ViewError{
		ViewName: ViewName,
		err:      err,
	}
}
