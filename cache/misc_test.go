package cache

import "testing"
import "bytes"

func TestRandom(t *testing.T) {
	var testLength = 10
	var testMaxLength = 1000
	var testBytes = []byte("s")
	b, err := RandomBytes(testLength)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != testLength {
		t.Errorf("RandomBytes length error %s", err)
	}
	for i := 0; i < 1000; i++ {
		b, err := NewRandomBytes(1, testBytes)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(testBytes, b) == 0 {
			t.Errorf("NewRandomBytes got same bytes error %s", string(b))
		}
	}
	b, err = RandMaskedBytes(TokenMask, testMaxLength)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != testMaxLength {
		t.Errorf("RandMaskedBytes length error %s", err)
	}
	for _, v := range b {
		if !bytes.Contains(TokenMask, []byte{v}) {
			t.Errorf("RandMaskedBytes TokenMask error %s", err)
		}
	}
	for i := 0; i < 1000; i++ {
		b, err := NewRandMaskedBytes(TokenMask, 1, testBytes)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(testBytes, b) == 0 {
			t.Errorf("NewRandomBytes got same bytes error %s", string(b))
		}
		if !bytes.Contains(TokenMask, b) {
			t.Errorf("RandMaskedBytes TokenMask error %s", err)
		}
	}
}
func TestMarshalMsgpack(t *testing.T) {
	var testData = "1234567890"
	var result string
	bytes, err := MarshalMsgpack(testData)
	if err != nil {
		t.Fatal(err)
	}
	err = UnmarshalMsgpack(bytes, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != testData {
		t.Errorf("UnmarshalMsgpack error %s", result)
	}
}
