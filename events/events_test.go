package events

import (
	"testing"
	"time"
)

func TestEvents(t *testing.T) {
	eventtype := Type("testtype")
	testtarget := "target"
	testdata := 15
	e := NewEvent()
	if e.Type != "" || e.Data != nil || e.Target != "" {
		t.Error(e)
	}
	e.WithType(eventtype).WithData(testdata).WithTarget(testtarget)
	d, ok := e.Data.(int)
	if e.Type != eventtype || ok != true || d != testdata || e.Target != testtarget {
		t.Error(e)
	}
}
func TestDefaultEvents(t *testing.T) {
	eventtype := Type("testtype")
	eventtype2 := Type("testtype2")
	e := eventtype.NewEvent()
	e2 := eventtype2.NewEvent()
	if e.Type == eventtype2 {
		t.Error(e)
	}
	if e.Type != eventtype {
		t.Error(e)
	}
	e3 := DefaultEventService.NewEvent()
	if e3.Type != "" {
		t.Error(e)
	}
	var result1 bool
	var result2 bool

	//test register event handler again.
	EmittEventType := WrapEmit(eventtype)
	OnEventType := WrapOn(eventtype)
	OnEventType(func(e *Event) {
		result1 = true
	})
	OnEventType(func(e *Event) {
		result2 = true
	})
	EmittEventType(e)
	time.Sleep(1 * time.Millisecond)
	if result1 != true || result2 != true {
		t.Error(e)
	}
	if DefaultEventService.Emit(e2) != false {
		t.Error(e2)
	}
	result1 = false
	result2 = false
	EmittEventType(nil)
	time.Sleep(1 * time.Millisecond)
	if result1 != true || result2 != true {
		t.Error(e)
	}
	if DefaultEventService.Emit(e2) != false {
		t.Error(e2)
	}

}
