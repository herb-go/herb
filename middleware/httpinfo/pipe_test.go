package httpinfo

import (
	"net/http/httptest"
	"testing"
)

func TestPipe(t *testing.T) {
	resp := NewResponse()
	p := NewBufferPipe(nil, resp)
	success := resp.UpdatePipe(p)
	if success != true {
		t.Fatal(success)
	}
	resp = NewResponse()
	writer := httptest.NewRecorder()
	p = NewBufferPipe(nil, resp)
	resp.WrapWriter(writer).Write([]byte{1})
	success = resp.UpdatePipe(p)
	if success != false {
		t.Fatal(success)
	}
}
