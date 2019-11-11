package guarder

import "testing"

func TestIdentifierDriver(t *testing.T) {
	func() {
		defer initDrivers()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal(r)
			}
			err := r.(error)
			if err == nil {
				t.Fatal(err)
			}
		}()
		RegisterIdentifier("test", nil)
	}()
	func() {
		defer initDrivers()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal(r)
			}
			err := r.(error)
			if err == nil {
				t.Fatal(err)
			}
		}()
		registerTokenFactory()
		registerTokenFactory()
	}()
	_, err := NewIdentifierDriver("notexistsdriver", nil, "")
	if err == nil {
		t.Fatal(err)
	}
}
