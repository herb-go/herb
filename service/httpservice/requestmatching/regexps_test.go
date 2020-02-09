package requestmatching

import (
	"net/http"
	"testing"
)

func TestRegexps(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.html", nil)
	re := NewRegExps()
	if !MustMatch(r, re) {
		t.Fatal(re)
	}
	re.Add("notexists")
	if MustMatch(r, re) {
		t.Fatal(re)
	}
	re.Add(".tml")
	if !MustMatch(r, re) {
		t.Fatal(re)
	}
}
