package middlewarefactory_test

import (
	"testing"
	"time"

	"github.com/herb-go/herb/middleware/middlewarefactory"
)

var emptyTimeCondition = middlewarefactory.TimeCondition{}

func newTimeCondition(start, end int64) *middlewarefactory.TimeCondition {
	return &middlewarefactory.TimeCondition{
		Start: start,
		End:   end,
	}
}
func TestTimeCondition(t *testing.T) {
	f := middlewarefactory.NewTimeConditionFactory()
	c, err := f(mustNewLoader(emptyTimeCondition))
	if err != nil {
		t.Fatal(err)
	}
	result, err := c.CheckRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatal(c)
	}
	now := time.Now()
	successStart := now.Add(-time.Hour).Unix()
	failStart := now.Add(time.Hour).Unix()
	successEnd := now.Add(time.Hour).Unix()
	failEnd := now.Add(-time.Hour).Unix()

	c, err = f(mustNewLoader(newTimeCondition(failStart, failEnd)))
	if err != nil {
		t.Fatal(err)
	}
	result, err = c.CheckRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Fatal(c)
	}
	c, err = f(mustNewLoader(newTimeCondition(successStart, failEnd)))
	if err != nil {
		t.Fatal(err)
	}
	result, err = c.CheckRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Fatal(c)
	}
	c, err = f(mustNewLoader(newTimeCondition(failStart, successEnd)))
	if err != nil {
		t.Fatal(err)
	}
	result, err = c.CheckRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Fatal(c)
	}
	c, err = f(mustNewLoader(newTimeCondition(successStart, successEnd)))
	if err != nil {
		t.Fatal(err)
	}
	result, err = c.CheckRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatal(c)
	}
}
