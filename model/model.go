package model

import (
	"errors"
	"fmt"
)

var ErrNoValidateModel = errors.New("error no validate method for model.")

type FieldError struct {
	Field string
	Label string
	Msg   string
}
type ValidatedResult struct {
	Validated bool
	Model     Validator
}
type Model struct {
	errors      []FieldError
	messages    ModelMessages
	fieldLabels map[string]string
}

func MustValidate(m Validator) bool {
	err := m.Validate()
	if err != nil {
		panic(err)
	}
	return !m.HasError()
}

func (model *Model) SetMessages(m Messages) {
	model.messages = m
}
func (model *Model) getMessageText(msg string) string {
	var msgtext string
	if model.messages != nil {
		msgtext = model.messages.GetMessage(msg)
	} else if DefaultMessages != nil {
		msgtext = DefaultMessages.GetMessage(msg)
	} else {
		msgtext = msg
	}
	return msgtext
}
func (model *Model) getMessageTextf(field, msg string) string {
	msg = model.getMessageText(msg)
	return fmt.Sprintf(msg+"%[3]s", model.GetFieldLabel(field), field, "")
}
func (model *Model) AddPlainError(field string, msg string) {
	f := FieldError{
		Field: field,
		Label: model.GetFieldLabel(field),
		Msg:   msg,
	}
	model.errors = append(model.Errors(), f)
}
func (model *Model) SetFieldLabels(labels map[string]string) {
	model.fieldLabels = labels
}
func (model *Model) GetFieldLabel(field string) string {
	if model.fieldLabels == nil {
		return DefaultMessages.GetMessage(field)
	}
	label, ok := model.fieldLabels[field]
	if ok == false {
		return field
	}
	return label
}
func (model *Model) AddError(field string, msg string) {
	model.AddPlainError(field, model.getMessageText(msg))
}
func (model *Model) AddErrorf(field string, msg string) {
	model.AddPlainError(field, model.getMessageTextf(field, msg))
}
func (model *Model) ValidateField(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		model.AddError(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Model:     model,
	}
}
func (model *Model) ValidateFieldf(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		model.AddErrorf(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Model:     model,
	}
}
func (model *Model) Errors() []FieldError {
	if model.errors == nil {
		return []FieldError{}
	} else {
		return model.errors
	}
}

func (model *Model) HasError() bool {
	return len(model.Errors()) != 0
}
func (model *Model) Validate() error {
	return ErrNoValidateModel
}

type Validator interface {
	HasError() bool
	Errors() []FieldError
	AddError(field string, msg string)
	AddErrorf(field string, msg string)
	Validate() error
}
