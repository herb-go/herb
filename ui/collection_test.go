package ui

import "testing"

func TestCollection(t *testing.T) {
	var result string
	message := NewMessages()
	message.Set("testlabel", "translated testlabel")
	message.Set("testlabel2", "translated testlabel2")
	messagemap := map[string]string{
		"test":  "testlabel",
		"test2": "testlabel2",
		"test3": "testlabel3",
	}
	c := NewCollection(message, messagemap)
	result = c.GetLabel("test")
	if result != "translated testlabel" {
		t.Fatal(result)
	}
	result = c.GetLabel("test2")
	if result != "translated testlabel2" {
		t.Fatal(result)
	}
	result = c.GetLabel("test3")
	if result != "testlabel3" {
		t.Fatal(result)
	}
	result = c.GetLabel("test4")
	if result != "test4" {
		t.Fatal(result)
	}
	c = NewCollection(nil, messagemap)
	result = c.GetLabel("test")
	if result != "testlabel" {
		t.Fatal(result)
	}
}
