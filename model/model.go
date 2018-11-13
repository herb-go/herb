package model

import (
	"errors"
	"fmt"
)

//ErrNoValidateMethod error rasied when model validate method not overrided.
var ErrNoValidateMethod = errors.New("error no validate method for model")

//ErrModelNotInited error raised when model id is empty.
var ErrModelNotInited = errors.New("model id is empty.You must use inited model . ")

//FieldError field error info struct
type FieldError struct {
	//Field field name.
	Field string
	//Label field label.
	//If field not found in model's labels,field label is same as field name.
	Label string
	//Msg error message
	Msg string
}

//ValidatedResult result of  validation.
type ValidatedResult struct {
	Validated bool
	Model     Validator
}

//Model model struct.
type Model struct {
	modelID     string
	errors      []FieldError
	messages    MessagesCollection
	fieldLabels map[string]string
}

//MustValidate return model validate result.
//Panic if any error raised.
func MustValidate(m Validator) bool {
	err := m.Validate()
	if err != nil {
		panic(err)
	}
	return !m.HasError()
}

//ModelID get model id
func (model *Model) ModelID() string {
	return model.modelID
}

//SetModelID set id to model.
func (model *Model) SetModelID(id string) {
	model.modelID = id
}

//SetMessages set translate messages to model.
//Messages will be used to translate msg field in AddPlainError,AddError,AddErrorf ,AddErrorf and ValidateFieldf method.
func (model *Model) SetMessages(m Messages) {
	model.messages = m
}
func (model *Model) getMessageText(msg string) string {
	var msgtext string
	if model.messages != nil {
		msgtext = model.messages.GetMessage(msg)
	} else {
		msgtext = DefaultMessages.GetMessage(msg)
	}
	return msgtext
}
func (model *Model) getMessageTextf(field, msg string) string {
	msg = model.getMessageText(msg)
	return fmt.Sprintf(msg+"%[3]s", model.GetFieldLabel(field), field, "")
}

//AddPlainError add plain error with given field and msg.
//Msg will not be translated.
func (model *Model) AddPlainError(field string, msg string) {
	f := FieldError{
		Field: field,
		Label: model.GetFieldLabel(field),
		Msg:   msg,
	}
	model.errors = append(model.Errors(), f)
}

//SetFieldLabels set field labels to model
func (model *Model) SetFieldLabels(labels map[string]string) {
	model.fieldLabels = labels
}

//GetFieldLabel get label by given label name.
//Return field name itself if not found in field labels of model.
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

//AddError add error by given field and plain msg.
//Msg will be translated.
func (model *Model) AddError(field string, msg string) {
	model.AddPlainError(field, model.getMessageText(msg))
}

//AddErrorf add error by given field and formatted msg.
//Msg will be translated.
func (model *Model) AddErrorf(field string, msg string) {
	model.AddPlainError(field, model.getMessageTextf(field, msg))
}

//ValidateField validated field then add error with given field name and plain msg if not validated.
func (model *Model) ValidateField(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		model.AddError(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Model:     model,
	}
}

//ValidateFieldf validated field then add error with given field name and formatted msg if not validated.
func (model *Model) ValidateFieldf(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		model.AddErrorf(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Model:     model,
	}
}

//Errors return error list of model
func (model *Model) Errors() []FieldError {
	if model.errors == nil {
		return []FieldError{}
	}
	return model.errors
}

//HasError return if model has any error.
func (model *Model) HasError() bool {
	return len(model.Errors()) != 0
}

//Validate method used to validate model.
//Fail validation will add error to model.
//Return any error if rasied.
//You must override this method for your own model,otherwise ErrNoValidateMethod will be raised.
func (model *Model) Validate() error {
	return ErrNoValidateMethod
}

//Validator interface for model that can be validated.
type Validator interface {
	//HasError return if model has any error.
	HasError() bool
	//Errors return error list of model
	Errors() []FieldError
	//AddError add error by given field and plain msg.
	AddError(field string, msg string)
	//AddErrorf add error by given field and formatted msg.
	AddErrorf(field string, msg string)
	//Validate method used to validate model.
	//Fail validation will add error to model.
	//Return any error if rasied.
	Validate() error
	//ModelID return model id.
	ModelID() string
	//SetModelID set model id.
	SetModelID(string)
}
