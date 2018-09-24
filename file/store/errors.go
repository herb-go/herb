package store

import (
	"errors"
	"fmt"
)

//ErrorType error type
type ErrorType string

//ErrorTypeNotExists error type file not  exists.
const ErrorTypeNotExists = "file-not-exists"

//Error file store error.
type Error struct {
	File string
	Err  error
	Type ErrorType
}

//NewError create new error with given file,type and error.
func NewError(file string, errorType ErrorType, err error) *Error {
	return &Error{
		File: file,
		Err:  err,
		Type: errorType,
	}
}

//Error return error msg.
func (e *Error) Error() string {
	return fmt.Sprintf("file error(%s):%s", e.File, e.Err.Error())
}

//NewNotExistsError create new file not exists error.
func NewNotExistsError(file string) *Error {
	return NewError(file, ErrorTypeNotExists, errors.New("file \""+file+"\" does not exist"))
}
