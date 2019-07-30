package translate

import "testing"

func TestTranslations(t *testing.T) {
	defer func() {
		DefaultTranslations = NewTranslations()
		Lang = ""
	}()
	DefaultTranslations = NewTranslations()
	m := NewMessages()
	m.Set("test", "test message")
	DefaultTranslations.SetMessages("testlang", "testmodule", m)

	result := Get("notexistmodule", "notexistkey")
	if result != "notexistkey" {
		t.Fatal(result)
	}
	Lang = "notexistlang"
	result = Get("notexistmodule", "notexistkey")
	if result != "notexistkey" {
		t.Fatal(result)
	}
	Lang = "testlang"
	result = Get("notexistmodule", "notexistkey")
	if result != "notexistkey" {
		t.Fatal(result)
	}
	result = Get("", "notexistkey")
	if result != "notexistkey" {
		t.Fatal(result)
	}
	result = Get("testmodule", "notexistkey")
	if result != "notexistkey" {
		t.Fatal(result)
	}
	result = Get("testmodule", "test")
	if result != "test message" {
		t.Fatal(result)
	}
	m = GetMessages("", "testmodule")
	if m == nil {
		t.Fatal(m)
	}
	c := m.Collection(nil)
	if c == nil {
		t.Fatal(c)
	}
}
