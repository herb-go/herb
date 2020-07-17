package keyencoding

import (
	"encoding/base64"

	"github.com/herb-go/herb/secret"
)

type Encoding struct {
	Encode func(secret.Blob) (string, error)
	Decode func(string) (secret.Blob, error)
}

var NopEncoding = &Encoding{
	Encode: func(s secret.Blob) (string, error) {
		return string(s), nil
	},
	Decode: func(s string) (secret.Blob, error) {
		return []byte(s), nil
	},
}

var Base64Encoding = &Encoding{
	Encode: func(s secret.Blob) (string, error) {
		return base64.StdEncoding.EncodeToString(s), nil
	},
	Decode: func(s string) (secret.Blob, error) {
		return base64.StdEncoding.DecodeString(s)
	},
}
