package model

import (
	"testing"
)

func TestMessage(t *testing.T) {
	var message1 = NewMessages().SetMessage("test1", "messagetest1")
	var message2 = &Messages{
		"test2": "messagetest2",
	}
	var mc = NewMessagesChain(message1)
	m, ok := message1.LoadMessage("test1")
	if ok == false {
		t.Error(ok)
	}
	if m != "messagetest1" {
		t.Error(m)
	}
	m = message1.GetMessage("test1")
	if m != "messagetest1" {
		t.Error(m)
	}
	m, ok = message1.LoadMessage("test2")
	if ok == true {
		t.Error(ok)
	}
	if m != "test2" {
		t.Error(m)
	}
	m = message1.GetMessage("test2")
	if m != "test2" {
		t.Error(m)
	}

	m, ok = mc.LoadMessage("test1")
	if ok == false {
		t.Error(ok)
	}
	if m != "messagetest1" {
		t.Error(m)
	}
	m = mc.GetMessage("test1")
	if m != "messagetest1" {
		t.Error(m)
	}
	m, ok = mc.LoadMessage("test2")
	if ok == true {
		t.Error(ok)
	}
	if m != "test2" {
		t.Error(m)
	}
	m = mc.GetMessage("test2")
	if m != "test2" {
		t.Error(m)
	}
	mc.Use(message2)
	m, ok = mc.LoadMessage("test2")
	if ok == false {
		t.Error(ok)
	}
	if m != "messagetest2" {
		t.Error(m)
	}
	m = mc.GetMessage("test2")
	if m != "messagetest2" {
		t.Error(m)
	}
}
