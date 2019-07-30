package translate

import (
	"testing"
)

func TestMessages(t *testing.T) {
	var message1 = NewMessages().Set("test1", "messagetest1")

	m, ok := message1.Load("test1")
	if ok == false {
		t.Error(ok)
	}
	if m != "messagetest1" {
		t.Error(m)
	}
	m = message1.Get("test1")
	if m != "messagetest1" {
		t.Error(m)
	}
	m, ok = message1.Load("test2")
	if ok == true {
		t.Error(ok)
	}
	if m != "test2" {
		t.Error(m)
	}
	m = message1.Get("test2")
	if m != "test2" {
		t.Error(m)
	}

}
