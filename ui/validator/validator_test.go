package validator

import (
	"testing"

	"github.com/herb-go/herb/ui"
)

type testModel struct {
	Field1 string
	Field2 string
	Field3 string
	Field4 string
	Validator
}
type testNoValidateModel struct {
	Validator
}

func newTestModel() *testModel {
	m := &testModel{}
	m.SetComponentLabels(ui.MapLabels(testLabels))
	return m
}

type str string

func (s str) Label() string {
	return string(s)
}
func (m *testModel) Validate() error {
	m.ValidateFieldf(m.Field1 != "", "Field1", "{{label}} required.")
	m.ValidateField(m.Field2 != "", "Field2", "Field2 required.")
	m.ValidateFieldf(m.Field3 != "", "Field3", "Field3 required")
	m.ValidateFieldLabelf(m.Field4 != "", "Field4", str("{{label}} required."))

	return nil
}

var testLabels = map[string]string{
	"Field1": "test field1",
	"Field2": "default field2",
	"Field4": "test field4",
}

func TestModel(t *testing.T) {
	m := newTestModel()
	if m.ComponentID() != "" {
		t.Fatal(m.ComponentID())
	}
	MustValidate(m)
	if !m.HasError() {
		t.Error(m.HasError())
	}
	modelerrors := m.Errors()
	if len(modelerrors) != 4 {
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
	if modelerrors[2].Msg != "Field3 required" {
		t.Error(modelerrors[2].Msg)
	}
	if modelerrors[3].Field != "Field4" {
		t.Error(modelerrors[3].Field)
	}
	if modelerrors[3].Label != "test field4" {
		t.Error(modelerrors[3].Label)
	}
	if modelerrors[3].Msg != "test field4 required." {
		t.Error(modelerrors[3].Msg)
	}
	m = newTestModel()
	m.Field1 = "value1"
	m.Field2 = "value2"
	m.Field3 = "value3"
	m.Field4 = "value4"
	MustValidate(m)
	if m.HasError() {
		t.Error(m.HasError())
	}
	if len(m.Errors()) != 0 {
		t.Error(m.Errors())
	}

}

func TestNoValidate(t *testing.T) {
	m := testNoValidateModel{}
	err := m.Validate()
	if err != ErrNoValidateMethod {
		t.Fatal(err)
	}
}
