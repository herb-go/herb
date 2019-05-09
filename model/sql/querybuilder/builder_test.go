package querybuilder

import "testing"

func TestBuilder(t *testing.T) {
	Debug = true
	defer func() { Debug = false }()
	RegisterDriver("test", DefaultDriver)
	b := NewBuilder()
	b.Driver = "test"
}
