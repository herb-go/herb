package hashed

import (
	"bytes"
	"testing"
)

func TestHash(t *testing.T) {
	var status *Status
	h := New()
	if !h.isEmpty() {
		t.Fatal(h)
	}
	d1 := NewData("key1", 1235, []byte("key1"))
	status = h.set(d1, 1234)
	if status.Changed != true ||
		status.Delta != 4 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1235 ||
		status.Size != 4 {
		t.Fatal(status)
	}
	status = h.set(d1, 1234)
	if status.Changed != true ||
		status.Delta != 0 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1235 ||
		status.Size != 4 {
		t.Fatal(status)
	}
	d2 := NewData("key2", 1236, []byte("keyd2"))
	status = h.set(d2, 1234)
	if status.Changed != true ||
		status.Delta != 5 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1236 ||
		status.Size != 9 {
		t.Fatal(status)
	}
	d3 := NewData("key3", 1000, []byte("keyd2"))
	status = h.set(d3, 1234)
	if status.Changed != true ||
		status.Delta != 0 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1236 ||
		status.Size != 9 {
		t.Fatal(status)
	}
	v1 := h.get(d1.Key, 1234)
	if bytes.Compare(v1.Data, d1.Data) != 0 {
		t.Fatal(v1)
	}
	v1e := h.get(d1.Key, 1250)
	if v1e != nil {
		t.Fatal(v1e)
	}
	vne := h.get("notexist", 1234)
	if vne != nil {
		t.Fatal(v1e)
	}
	d4 := NewData("key4", 1235, []byte("keyd4"))
	status = h.update(d4, 1234)
	if status.Changed != false ||
		status.Delta != 0 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1236 ||
		status.Size != 9 {
		t.Fatal(status)
	}
	d5 := NewData("key1", 1250, []byte("1"))
	status = h.update(d5, 1234)
	if status.Changed != true ||
		status.Delta != -3 ||
		status.FirstExpired != 1236 ||
		status.LastExpired != 1250 ||
		status.Size != 6 {
		t.Fatal(status)
	}
	status = h.update(d5, 1240)
	if status.Changed != true ||
		status.Delta != -5 ||
		status.FirstExpired != 1250 ||
		status.LastExpired != 1250 ||
		status.Size != 1 {
		t.Fatal(status)
	}
	status = h.delete(d4.Key, 1240)
	if status.Changed != false ||
		status.Delta != 0 ||
		status.FirstExpired != 1250 ||
		status.LastExpired != 1250 ||
		status.Size != 1 {
		t.Fatal(status)
	}
	status = h.delete(d5.Key, 1240)
	if status.Changed != true ||
		status.Delta != -1 ||
		status.FirstExpired != 0 ||
		status.LastExpired != 0 ||
		status.Size != 0 {
		t.Fatal(status)
	}
	d1 = NewData("key1", 1235, []byte("key1"))
	h.set(d1, 1234)
	h.set(d2, 1234)
	v1 = h.get(d1.Key, 1234)
	if v1 == nil || v1.Expired != 1235 {
		t.Fatal(h)
	}
	status = h.expired(d1.Key, 1250, 1234)
	if status.Changed != true ||
		status.Delta != 0 ||
		status.FirstExpired != 1236 ||
		status.LastExpired != 1250 ||
		status.Size != 9 {
		t.Fatal(status)
	}
	v1 = h.get(d1.Key, 1234)
	if v1.Expired != 1250 {
		t.Fatal(h)
	}
	status = h.expired(d1.Key, 1250, 1240)
	if status.Changed != true ||
		status.Delta != -5 ||
		status.FirstExpired != 1250 ||
		status.LastExpired != 1250 ||
		status.Size != 4 {
		t.Fatal(status)
	}
	h.set(d2, 1234)
	status = h.selfcheck(1240)
	if status.Changed != true ||
		status.Delta != -5 ||
		status.FirstExpired != 1250 ||
		status.LastExpired != 1250 ||
		status.Size != 4 {
		t.Fatal(status)
	}
}
