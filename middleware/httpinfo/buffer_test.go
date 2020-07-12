package httpinfo

import (
	"net/http/httptest"
	"testing"
)

func TestBuffer(t *testing.T) {
	resp := NewResponse()
	success := resp.BuildBuffer(nil, nil)
	if success != true {
		t.Fatal(success)
	}
	success = resp.BuildBuffer(nil, nil)
	if success != false {
		t.Fatal(success)
	}
	resp = NewResponse()
	writer := httptest.NewRecorder()
	resp.WrapWriter(writer).Write([]byte{1})
	success = resp.BuildBuffer(nil, nil)
	if success != false {
		t.Fatal(success)
	}
}
