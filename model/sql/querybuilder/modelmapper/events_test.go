package modelmapper

import "testing"

func TestEvent(t *testing.T) {
	type Model struct {
		CommonQueryEvents
		value interface{}
	}
	model := &Model{}
	if model.BeforeUpdate() != nil {
		t.Fatal(model)
	}
	if model.AfterUpdate() != nil {
		t.Fatal(model)
	}
	if model.BeforeInsert() != nil {
		t.Fatal(model)
	}
	if model.AfterInsert() != nil {
		t.Fatal(model)
	}
	if model.AfterDelete() != nil {
		t.Fatal(model)
	}
	if model.AfterFind() != nil {
		t.Fatal(model)
	}
}
