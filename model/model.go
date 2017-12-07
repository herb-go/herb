package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FieldError struct {
	Field string
	Label string
	Msg   string
}
type ValidatedResult struct {
	Validated bool
	Model     *ModelErrors
}
type ModelErrors struct {
	errors      []FieldError
	messages    ModelMessages
	fieldLabels map[string]string
}

func MustValidate(m Model) bool {
	err := m.Validate()
	if err != nil {
		panic(err)
	}
	return !m.HasError()
}
func MustValidateJSONPost(r *http.Request, m HttpModel) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}
	err = m.InitWithRequest(r)
	if err != nil {
		panic(err)
	}
	return MustValidate(m)

}
func MustRenderErrors(w http.ResponseWriter, m Model) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)
	bytes, err := json.Marshal(m.Errors())
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}
func (model *ModelErrors) SetMessages(m *Messages) {
	model.messages = m
}
func (model *ModelErrors) getMessageText(msg string) string {
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
func (model *ModelErrors) getMessageTextf(field, msg string) string {
	msg = model.getMessageText(msg)
	return fmt.Sprintf("%[3]s"+msg, model.GetFieldLabel(field), field, "")
}
func (model *ModelErrors) AddPlainError(field string, msg string) {
	f := FieldError{
		Field: field,
		Label: model.GetFieldLabel(field),
		Msg:   msg,
	}
	model.errors = append(model.Errors(), f)
}
func (model *ModelErrors) SetFieldLabels(labels map[string]string) {
	model.fieldLabels = labels
}
func (model *ModelErrors) GetFieldLabel(field string) string {
	if model.fieldLabels == nil {
		return field
	}
	label, ok := model.fieldLabels[field]
	if ok == false {
		return field
	}
	return label
}
func (model *ModelErrors) AddError(field string, msg string) {
	model.AddPlainError(field, model.getMessageText(msg))
}
func (model *ModelErrors) AddErrorf(field string, msg string) {
	model.AddPlainError(field, model.getMessageTextf(field, msg))
}
func (model *ModelErrors) ValidateField(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		model.AddError(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Model:     model,
	}
}
func (model *ModelErrors) ValidateFieldf(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		model.AddErrorf(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Model:     model,
	}
}
func (model *ModelErrors) Errors() []FieldError {
	if model.errors == nil {
		return []FieldError{}
	} else {
		return model.errors
	}
}
func (model *ModelErrors) InitWithRequest(*http.Request) error {
	return nil
}

func (model *ModelErrors) HasError() bool {
	return len(model.Errors()) != 0
}

type Model interface {
	HasError() bool
	Errors() []FieldError
	AddError(field string, msg string)
	Validate() error
}

type HttpModel interface {
	Model
	InitWithRequest(*http.Request) error
}
