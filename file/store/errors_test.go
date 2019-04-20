package store

import (
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	err := NewNotExistsError("test.dat")
	if err == nil {
		t.Fatal(err)
	}
	if err.File != "test.dat" {
		t.Fatal(err)
	}
	if err.Type != ErrorTypeNotExists {
		t.Fatal(err)
	}
	if !strings.Contains(err.Error(), "test.dat") {
		t.Fatal(err)
	}
}
