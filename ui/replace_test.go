package ui

import "testing"

func TestReplace(t *testing.T) {
	var result string
	m := map[string]string{
		"test":  "mapped-test-",
		"test2": "mapped-test2-",
		"{":     "mapeed-{-",
	}
	result = Replace("abcde", m)
	if result != "abcde" {
		t.Fatal(result)
	}
	result = Replace("{{test}}2", m)
	if result != "mapped-test-2" {
		t.Fatal(result)
	}
	result = Replace("\\{{test\\}}2", m)
	if result != "{{test}}2" {
		t.Fatal(result)
	}
	result = Replace("\\\\{{test}}2", m)
	if result != "\\mapped-test-2" {
		t.Fatal(result)
	}
	result = Replace("{{test}}{{notexist}}{{test2}}", m)
	if result != "mapped-test-{{notexist}}mapped-test2-" {
		t.Fatal(result)
	}
}
