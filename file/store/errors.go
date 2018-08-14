package store

import (
	"errors"
	"fmt"
)

type ErrorType string

const ErrorTypeNotExists = "file-not-exists"

type Error struct {
	File string
	Err  error
	Type ErrorType
}

func NewError(file string, errorType ErrorType, err error) *Error {
	return &Error{
		File: file,
		Err:  err,
		Type: errorType,
	}
}
func (e *Error) Error() string {
	return fmt.Sprintf("file error(%s):%s", e.File, e.Err.Error())
}

func NewNotExistsError(file string) *Error {
	return NewError(file, ErrorTypeNotExists, errors.New("file \""+file+"\" does not exist"))
}
