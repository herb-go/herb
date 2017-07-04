package main

import "github.com/herb-go/herb/model"

type exampleFormModel struct {
	model.ModelErrors
	Data string
}

func (m *exampleFormModel) Validate() error {
	m.ValidateFieldf("Data", MsgFormFieldRequired, m.Data != "")
	return nil
}

func InitModel() {

}
