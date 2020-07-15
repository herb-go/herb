package token

import (
	"encoding/base64"
)

type Secret []byte

type Encoding struct {
	Encode func(Secret) (string, error)
	Decode func(string) (Secret, error)
}

var StringEncoding = &Encoding{
	Encode: func(s Secret) (string, error) {
		return string(s), nil
	},
	Decode: func(s string) (Secret, error) {
		return Secret(s), nil
	},
}

var Base64Encoding = &Encoding{
	Encode: func(s Secret) (string, error) {
		return base64.StdEncoding.EncodeToString(s), nil
	},
	Decode: func(s string) (Secret, error) {
		return base64.StdEncoding.DecodeString(s)
	},
}
