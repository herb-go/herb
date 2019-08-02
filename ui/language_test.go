package ui

import "testing"

func TestTranslatable(t *testing.T) {
	l := Language{}
	if l.Lang() != "" {
		t.Fatal(l)
	}
	l.SetLang("zh-cn")
	if l.Lang() != "zh-cn" {
		t.Fatal(l)
	}
}
