package guarder

import "testing"

func TestCredentialDriver(t *testing.T) {
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
		RegisterCredential("test", nil)
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
		registerTokenMapFactory()
		registerTokenMapFactory()
	}()
	_, err := NewCredentialDriver("notexistsdriver", nil, "")
	if err == nil {
		t.Fatal(err)
	}
}
