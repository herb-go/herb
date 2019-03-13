package msgpackmarshaler

import (
	"testing"

	"github.com/herb-go/herb/cache"
)

func TestMsgpack(t *testing.T) {
	var testdata = map[string]string{
		"test": "test",
	}
	marshaler, err := cache.NewMarshaler("msgpack")
	if err != nil {
		t.Fatal(err)
	}
	var v = map[string]string{}
	bs, err := marshaler.Marshal(testdata)
	if err != nil {
		t.Fatal(err)
	}
	err = marshaler.Unmarshal(bs, &v)
	if err != nil {
		t.Fatal(err)
	}
	if v["test"] != testdata["test"] {
		t.Fatal(v)
	}
}
