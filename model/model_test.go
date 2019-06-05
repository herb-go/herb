package model

import (
	"testing"
)

type testModel struct {
	Field1 string
	Field2 string
	Field3 string
	Model
}
type testNoValidateModel struct {
	Model
}

func newTestModel() *testModel {
	m := &testModel{}
	m.SetMessages(testMessage)
	m.SetFieldLabels(testLabels)
	return m
}
func (m *testModel) Validate() error {
	m.ValidateFieldf(m.Field1 != "", "Field1", "Field1 required.")
	m.ValidateField(m.Field2 != "", "Field2", "Field2 required.")
	m.ValidateFieldf(m.Field3 != "", "Field3", "Field3 required")

	return nil
}

var testMessage = &Messages{
	"Field1 required.": "%[1]s required.",
}

var testLabels = Messages{
	"Field1": "test field1",
}

var defaultMessage = &Messages{
	"Field2":          "default field2",
	"Field3 required": "test field3",
}

func TestModel(t *testing.T) {
	DefaultMessages = NewMessageChain()
	DefaultMessages.Use(defaultMessage)
	m := newTestModel()
	m.SetModelID("test")
	if m.ModelID() != "test" {
		t.Fatal(m.ModelID())
	}
	MustValidate(m)
	if !m.HasError() {
		t.Error(m.HasError())
	}
	modelerrors := m.Errors()
	if len(modelerrors) != 3 {
		t.Error(modelerrors)
	}
	if modelerrors[0].Field != "Field1" {
		t.Error(modelerrors[0].Field)
	}
	if modelerrors[0].Label != "test field1" {
		t.Error(modelerrors[0].Label)
	}
	if modelerrors[0].Msg != "test field1 required." {
		t.Error(modelerrors[0].Msg)
	}
	if modelerrors[1].Field != "Field2" {
		t.Error(modelerrors[1].Field)
	}
	if modelerrors[1].Label != "Field2" {
		t.Error(modelerrors[1].Label)
	}
	if modelerrors[1].Msg != "Field2 required." {
		t.Error(modelerrors[1].Msg)
	}
	if modelerrors[2].Field != "Field3" {
		t.Error(modelerrors[2].Field)
	}
	if modelerrors[2].Label != "Field3" {
		t.Error(modelerrors[2].Label)
	}
	if modelerrors[2].Msg != "Field3 required" {
		t.Error(modelerrors[2].Msg)
	}
	m = newTestModel()
	m.Field1 = "value1"
	m.Field2 = "value2"
	m.Field3 = "value3"
	MustValidate(m)
	if m.HasError() {
		t.Error(m.HasError())
	}
	if len(m.Errors()) != 0 {
		t.Error(m.Errors())
	}

}

func TestNilMessage(t *testing.T) {
	DefaultMessages = NewMessageChain()
	DefaultMessages.Use(defaultMessage)
	m := &testModel{}
	MustValidate(m)
	if !m.HasError() {
		t.Error(m.HasError())
	}
	modelerrors := m.Errors()
	if len(modelerrors) != 3 {
		t.Error(modelerrors)
	}
	if modelerrors[0].Field != "Field1" {
		t.Error(modelerrors[0].Field)
	}
	if modelerrors[0].Label != "Field1" {
		t.Error(modelerrors[0].Label)
	}
	if modelerrors[0].Msg != "Field1 required." {
		t.Error(modelerrors[0].Msg)
	}
	if modelerrors[1].Field != "Field2" {
		t.Error(modelerrors[1].Field)
	}
	if modelerrors[1].Label != "default field2" {
		t.Error(modelerrors[1].Label)
	}
	if modelerrors[1].Msg != "Field2 required." {
		t.Error(modelerrors[1].Msg)
	}
	if modelerrors[2].Field != "Field3" {
		t.Error(modelerrors[2].Field)
	}
	if modelerrors[2].Label != "Field3" {
		t.Error(modelerrors[2].Label)
	}
	if modelerrors[2].Msg != "test field3" {
		t.Error(modelerrors[2].Msg)
	}
}

func TestNoValidate(t *testing.T) {
	m := testNoValidateModel{}
	err := m.Validate()
	if err != ErrNoValidateMethod {
		t.Fatal(err)
	}
}
