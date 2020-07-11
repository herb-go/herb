package httpinfo

import "testing"

func TestValidator(t *testing.T) {
	ok, err := ValidatorAlways.Validate(nil, nil)
	if ok != true || err != nil {
		t.Fatal(ok, err)
	}
	ok, err = ValidatorNever.Validate(nil, nil)
	if ok != false || err != nil {
		t.Fatal(ok, err)
	}
}
