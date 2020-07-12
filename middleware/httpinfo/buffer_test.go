package httpinfo

import (
	"net/http/httptest"
	"testing"
)

func TestBuffer(t *testing.T) {
	resp := NewResponse()
	c := NewBufferController(nil, resp)
	success := resp.UpdateController(c)
	if success != true {
		t.Fatal(success)
	}
	resp = NewResponse()
	writer := httptest.NewRecorder()
	c = NewBufferController(nil, resp)
	resp.WrapWriter(writer).Write([]byte{1})
	success = resp.UpdateController(c)
	if success != false {
		t.Fatal(success)
	}
}
