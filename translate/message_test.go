package translate

import "testing"

func TestMessage(t *testing.T) {
	var result string
	defer func() {
		DefaultTranslations = NewTranslations()
		Lang = ""
	}()
	Lang = "testlang"
	DefaultTranslations = NewTranslations()
	m := NewMessages()
	m.Set("test", "translated test")
	DefaultTranslations.SetMessages("testlang", "testmodule", m)
	message := NewMessage("testmodule", "test")
	result = message.Translate("")
	if result != "translated test" {
		t.Fatal(result)
	}
	result = message.TranslateWith(DefaultTranslations, "")
	if result != "test" {
		t.Fatal(result)
	}
	result = message.Translated("").String()
	if result != "translated test" {
		t.Fatal(result)
	}
}

func TestTemplateMessage(t *testing.T) {
	var result string
	defer func() {
		DefaultTranslations = NewTranslations()
		Lang = ""
	}()
	Lang = "testlang"
	DefaultTranslations = NewTranslations()
	m := NewMessages()
	m.Set("{{test}}", "translated {{test}}")
	DefaultTranslations.SetMessages("testlang", "testmodule", m)
	message := NewTemplateMessage("testmodule", "{{test}}", map[string]string{"test": "testtemplate"})
	result = message.Translate("")
	if result != "translated testtemplate" {
		t.Fatal(result)
	}
	result = message.TranslateWith(DefaultTranslations, "")
	if result != "testtemplate" {
		t.Fatal(result)
	}
	result = message.Translated("").String()
	if result != "translated testtemplate" {
		t.Fatal(result)
	}
}
