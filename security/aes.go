package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var filledByte = []byte{0}

func formatKey(key []byte, size int) []byte {
	var data = make([]byte, size)
	copy(data, key)
	return data
}
func AESEncrypt(unencrypted []byte, key []byte) (encrypted []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	cryptKey := formatKey(key, aes.BlockSize)
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return
	}
	data := PKCS7Padding(unencrypted, aes.BlockSize)
	crypter := cipher.NewCBCEncrypter(block, cryptKey)
	encrypted = make([]byte, len(data))
	crypter.CryptBlocks(encrypted, data)
	return
}
func AESEncryptBase64(unencrypted []byte, key []byte) (encrypted string, err error) {
	d, err := AESEncrypt(unencrypted, key)
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(d), nil
}
func AESDecrypt(encrypted []byte, key []byte) (dencrypted []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	cryptKey := formatKey(key, aes.BlockSize)
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return
	}
	crypter := cipher.NewCBCDecrypter(block, cryptKey)
	data := make([]byte, len(encrypted))
	crypter.CryptBlocks(data, encrypted)
	dencrypted = PKCS7Unpadding(data)
	return
}

func AESDecryptBase64(encrypted string, key []byte) (dencrypted []byte, err error) {
	d, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return
	}
	return AESDecrypt(d, key)
}
