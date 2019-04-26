package querybuilder

import "testing"

func TestBuilder(t *testing.T) {
	Debug = true
	defer func() { Debug = false }()
	RegisterBuilder("test", DefaultBuilder)
	b := NewBuilder()
	b.Driver = "test"
}
