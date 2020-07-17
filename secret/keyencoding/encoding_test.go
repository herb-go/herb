package keyencoding

import (
	"bytes"
	"testing"
)

func TestEncoding(t *testing.T) {
	secret := []byte{1, 2, 3, 4, 5}
	encoded, err := Base64Encoding.Encode(secret)
	if err != nil {
		panic(err)
	}
	if encoded == string(secret) {
		t.Fatal(encoded)
	}
	decoded, err := Base64Encoding.Decode(encoded)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(decoded, secret) != 0 {
		t.Fatal(decoded)
	}
	encoded, err = NopEncoding.Encode(secret)
	if err != nil {
		panic(err)
	}
	if encoded != string(secret) {
		t.Fatal(encoded)
	}
	decoded, err = NopEncoding.Decode(string(secret))
	if err != nil {
		panic(err)
	}
	if bytes.Compare(decoded, secret) != 0 {
		t.Fatal(decoded)
	}
}
