package middlewarefactory_test

import (
	"net/http"
	"testing"

	"github.com/herb-go/herb/middleware/middlewarefactory"
)

type fixedCondition struct {
	Value bool
}

func (c *fixedCondition) CheckRequest(*http.Request) (bool, error) {
	return c.Value, nil
}

var fixedSuccessCondition = &fixedCondition{
	Value: true,
}

var fixedFailCondition = &fixedCondition{
	Value: false,
}

func mustCheck(result bool, err error) bool {
	if err != nil {
		panic(err)
	}
	return result
}
func mustCheckCondition(c middlewarefactory.Condition) bool {
	return mustCheck(c.CheckRequest(nil))
}
func TestCondition(t *testing.T) {
	if !mustCheck(middlewarefactory.Not(nil, fixedFailCondition)) {
		t.Fatal()
	}
	if mustCheck(middlewarefactory.Not(nil, fixedSuccessCondition)) {
		t.Fatal()
	}
	if !mustCheck(middlewarefactory.And(nil)) {
		t.Fatal()
	}
	if mustCheck(middlewarefactory.And(nil, fixedFailCondition, fixedSuccessCondition)) {
		t.Fatal()
	}
	if !mustCheck(middlewarefactory.And(nil, fixedSuccessCondition, fixedSuccessCondition)) {
		t.Fatal()
	}
	if mustCheck(middlewarefactory.Or(nil)) {
		t.Fatal()
	}
	if !mustCheck(middlewarefactory.Or(nil, fixedFailCondition, fixedSuccessCondition)) {
		t.Fatal()
	}
	if !mustCheck(middlewarefactory.Or(nil, fixedSuccessCondition, fixedSuccessCondition)) {
		t.Fatal()
	}
	if mustCheck(middlewarefactory.Or(nil, fixedFailCondition, fixedFailCondition)) {
		t.Fatal()
	}
}

func TestPlainCondition(t *testing.T) {
	var c *middlewarefactory.PlainCondition
	c = middlewarefactory.NewPlainCondition()
	if !mustCheckCondition(c) {
		t.Fatal()
	}
	c.Not = true
	if mustCheckCondition(c) {
		t.Fatal()
	}
	c = middlewarefactory.NewPlainCondition()
	c.Disabled = true
	if mustCheckCondition(c) {
		t.Fatal()
	}
	c = middlewarefactory.NewPlainCondition()
	c.Condition = fixedFailCondition
	if mustCheckCondition(c) {
		t.Fatal()
	}
	c = middlewarefactory.NewPlainCondition()
	c.Conditions = append(c.Conditions, fixedSuccessCondition, fixedSuccessCondition)
	if !mustCheckCondition(c) {
		t.Fatal()
	}
	c = middlewarefactory.NewPlainCondition()
	c.Conditions = append(c.Conditions, fixedSuccessCondition, fixedFailCondition)
	if mustCheckCondition(c) {
		t.Fatal()
	}
	c = middlewarefactory.NewPlainCondition()
	c.Or = true
	c.Conditions = append(c.Conditions, fixedSuccessCondition, fixedSuccessCondition)
	if !mustCheckCondition(c) {
		t.Fatal()
	}
	c = middlewarefactory.NewPlainCondition()
	c.Or = true
	c.Conditions = append(c.Conditions, fixedSuccessCondition, fixedFailCondition)
	if !mustCheckCondition(c) {
		t.Fatal()
	}
}
