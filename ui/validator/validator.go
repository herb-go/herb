package validator

import (
	"errors"
	"fmt"

	"github.com/herb-go/herb/ui"
)

//ErrNoValidateMethod error rasied when model validate method not overrided.
var ErrNoValidateMethod = errors.New("error no validate method for model")

//ErrValidatorNotInited error raised when model id is empty.
var ErrValidatorNotInited = errors.New("model id is empty.You must use inited model . ")

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
	Validator Fields
}

//Validator model struct.
type Validator struct {
	modelID string
	errors  []*FieldError
	labels  ui.Labels
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

//ValidatorID get model id
func (v *Validator) ValidatorID() string {
	return v.modelID
}

//SetValidatorID set id to model.
func (v *Validator) SetValidatorID(id string) {
	v.modelID = id
}

func (v *Validator) getTextf(field, msg string) string {
	return fmt.Sprintf(msg+"%[3]s", v.GetFieldLabel(field), field, "")
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

//SetFieldLabels set field labels to model
func (v *Validator) SetFieldLabels(labels map[string]string) {
	v.SetFieldLabelsCollection(ui.MapLabels(labels))
}

//SetFieldLabelsCollection set field labels collection to model
func (v *Validator) SetFieldLabelsCollection(labels ui.Labels) {
	v.labels = labels
}

//GetFieldLabel get label by given label name.
//Return field name itself if not found in field labels of model.
func (v *Validator) GetFieldLabel(field string) string {
	if v.labels == nil {
		return field
	}
	return v.labels.GetLabel(field)
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
func (v *Validator) ValidateField(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		v.AddError(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Validator: v,
	}
}

//ValidateFieldf validated field then add error with given field name and formatted msg if not validated.
func (v *Validator) ValidateFieldf(validated bool, field string, msg string) *ValidatedResult {
	if !validated {
		v.AddErrorf(field, msg)
	}
	return &ValidatedResult{
		Validated: validated,
		Validator: v,
	}
}

//ValidateFieldfString validated field then add error with given field name and  string interfcae msg if not validated.
func (v *Validator) ValidateFieldfString(validated bool, field string, msg ui.Label) *ValidatedResult {
	if !validated {
		v.AddErrorf(field, msg.Label())
	}
	return &ValidatedResult{
		Validated: validated,
		Validator: v,
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
	//ValidatorID return model id.
	ValidatorID() string
	//SetValidatorID set model id.
	SetValidatorID(string)
	//SetFieldLabelsCollection set field labels collection to model
	SetFieldLabelsCollection(labels ui.Labels)
	//GetFieldLabel get label by given label name.
	GetFieldLabel(field string) string
}
