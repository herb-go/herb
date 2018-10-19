package session

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestAES(t *testing.T) {
	var byteData = []byte("testdata12345678")
	var iv = []byte("1234567890abcdef")
	var key = []byte("key")
	encryptedBytes, err := AESEncrypt(byteData, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	resultBytes, err := AESDecrypt(encryptedBytes, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBytes, byteData) != 0 {
		t.Errorf("AES decrypt error")
	}
	encryptedBase64, err := AESEncryptBase64(byteData, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	_, err = base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		t.Fatal(err)
	}
	resultBase64, err := AESDecryptBase64(encryptedBase64, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBase64, byteData) != 0 {
		t.Errorf("AES decrypt error")
	}

}

func TestAESNonce(t *testing.T) {
	var byteData = []byte("testdata12345678abc")
	var key = []byte("key")
	var repeatLimit = 20
	encryptedBytes, err := AESNonceEncrypt(byteData, key)
	if err != nil {
		t.Fatal(err)
	}
	resultBytes, err := AESNonceDecrypt(encryptedBytes, key)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBytes, byteData) != 0 {
		t.Errorf("AES decrypt error")
	}
	encryptedBase64, err := AESNonceEncryptBase64(byteData, key)
	if err != nil {
		t.Fatal(err)
	}
	_, err = base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		t.Fatal(err)
	}
	resultBase64, err := AESNonceDecryptBase64(encryptedBase64, key)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBase64, byteData) != 0 {
		t.Errorf("AES decrypt error")
	}
	repeatBytes := make([][]byte, repeatLimit)
	for i := 0; i < repeatLimit; i++ {
		repeatBytes[i], err = AESNonceEncrypt(byteData, key)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, v1 := range repeatBytes {
		sameCount := 0
		for _, v2 := range repeatBytes {
			if bytes.Compare(v1, v2) == 0 {
				sameCount++
			}
		}
		if sameCount == repeatLimit {
			t.Errorf("Repeat nonce test fail")
		}
	}
	repeatBytesBase64 := make([]string, repeatLimit)
	for i := 0; i < repeatLimit; i++ {
		repeatBytesBase64[i], err = AESNonceEncryptBase64(byteData, key)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, v1 := range repeatBytesBase64 {
		sameCount := 0
		for _, v2 := range repeatBytesBase64 {
			if v1 == v2 {
				sameCount++
			}
		}
		if sameCount == repeatLimit {
			t.Errorf("Repeat nonce test fail")
		}
	}
}
