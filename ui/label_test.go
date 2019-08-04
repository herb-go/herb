package ui

import "testing"

func TestLabel(t *testing.T) {
	var ml Labels = MapLabels(map[string]string{
		"test": "test",
	})
	if ml.GetLabel("test") != "test" {
		t.Fatal(ml)
	}
	if ml.GetLabel("notexist") != "" {
		t.Fatal(ml)
	}
	var l Label = StringLabel("test")
	if l.Label() != "test" {
		t.Fatal(l)
	}
}
