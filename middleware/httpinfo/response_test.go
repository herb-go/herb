package httpinfo

import "testing"

func TestResponse(t *testing.T) {
	resp := NewResponse()
	if resp.StatusCode != 200 || resp.Written {
		t.Fatal(resp)
	}
}
