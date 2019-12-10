package store

import "testing"

func TestRegister(t *testing.T) {
	fs := Factories()
	if len(fs) != 1 {
		t.Fatal(fs)
	}
	UnregisterAll()
	fs = Factories()
	if len(fs) != 0 {
		t.Fatal(fs)
	}
	RegisterAssets()
}

func TestEmptyDriver(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
	}()
	Register("test", nil)
}

func TestDupDriver(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
	}()
	Register("assets", func(func(interface{}) error) (Driver, error) {
		return nil, nil
	})
}

func TestNotExistDriver(t *testing.T) {
	_, err := NewDriver("notexist", nil)
	if err == nil {
		t.Fatal("err")
	}
}
