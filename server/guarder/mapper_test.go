package guarder

import "testing"

func TestMapperDriver(t *testing.T) {
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
		RegisterMapper("test", nil)
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
		registerBasicAuthFactory()
		registerBasicAuthFactory()
	}()
	_, err := NewMapperDriver("notexistsdriver", nil, "")
	if err == nil {
		t.Fatal(err)
	}
}
