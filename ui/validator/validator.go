package validator

import (
	"errors"

	"github.com/herb-go/herb/ui"
)

//ErrNoValidateMethod error rasied when model validate method not overrided.
var ErrNoValidateMethod = errors.New("error no validate method for model")

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
// type ValidatedResult struct {
// 	Validated bool
// 	Validator Fields
// }

//Validator model struct.
type Validator struct {
	ui.Language
	ui.LabelsComponent
	errors []*FieldError
}

//MustValidate return model validate result.
//Panic if any error raised.
func MustValidate(m Fields) bool {
	err := m.Validate()
	if err != nil {
		panic(err)
	}
	return !m.HasError()
}

//ComponentID get ui  component id
func (v *Validator) ComponentID() string {
	return ""
}

func (v *Validator) getTextf(field, msg string) string {
	return ui.Replace(msg, map[string]string{
		"label": v.GetFieldLabel(field),
		"field": field,
	})
}

//AddPlainError add plain error with given field and msg.
//Msg will not be translated.
func (v *Validator) AddPlainError(field string, msg string) {
	f := FieldError{
		Field: field,
		Label: v.GetFieldLabel(field),
		Msg:   msg,
	}
	v.errors = append(v.Errors(), &f)
}

//GetFieldLabel get label by given label name.
//Return field name itself if not found in field labels of model.
func (v *Validator) GetFieldLabel(field string) string {
	l := v.GetLabel(field)
	if l == "" {
		return field
	}
	return l
}

//AddError add error by given field and plain msg.
//Msg will be translated.
func (v *Validator) AddError(field string, msg string) {
	v.AddPlainError(field, msg)
}

//AddErrorf add error by given field and formatted msg.
//Msg will be translated.
func (v *Validator) AddErrorf(field string, msg string) {
	v.AddPlainError(field, v.getTextf(field, msg))
}

//ValidateField validated field then add error with given field name and plain msg if not validated.
func (v *Validator) ValidateField(validated bool, field string, msg string) {
	if !validated {
		v.AddError(field, msg)
	}

}

//ValidateFieldf validated field then add error with given field name and formatted msg if not validated.
func (v *Validator) ValidateFieldf(validated bool, field string, msg string) {
	if !validated {
		v.AddErrorf(field, msg)
	}
}

//ValidateFieldLabelf validated field then add error with given field name and  string interfcae msg if not validated.
func (v *Validator) ValidateFieldLabelf(validated bool, field string, msg ui.Label) {
	if !validated {
		v.AddErrorf(field, msg.Label())
	}

}

//Errors return error list of model
func (v *Validator) Errors() []*FieldError {
	if v.errors == nil {
		return []*FieldError{}
	}
	return v.errors
}

//HasError return if model has any error.
func (v *Validator) HasError() bool {
	return len(v.Errors()) != 0
}

//Validate method used to validate model.
//Fail validation will add error to model.
//Return any error if rasied.
//You must override this method for your own model,otherwise ErrNoValidateMethod will be raised.
func (v *Validator) Validate() error {
	return ErrNoValidateMethod
}

//Fields interface for model that can be validated.
type Fields interface {
	ui.TranslationLanguage
	ui.Component
	ui.ComponentLabels
	//HasError return if model has any error.
	HasError() bool
	//Errors return error list of model
	Errors() []*FieldError
	//AddError add error by given field and plain msg.
	AddError(field string, msg string)
	//AddErrorf add error by given field and formatted msg.
	AddErrorf(field string, msg string)
	//Validate method used to validate model.
	//Fail validation will add error to model.
	//Return any error if rasied.
	Validate() error
	//GetFieldLabel get label by given label name.
	GetFieldLabel(field string) string
}
