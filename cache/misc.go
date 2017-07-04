package cache

import (
	"bytes"
	"crypto/rand"
)

var TokenMask = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomBytes(length int) ([]byte, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	return token, err
}
func NewRandomBytes(length int, origin []byte) ([]byte, error) {
	for {
		token, err := RandomBytes(length)
		if err != nil {
			return token, err
		}
		if bytes.Compare(token, origin) != 0 {
			return token, nil
		}
	}
}

func RandMaskedBytes(mask []byte, length int) ([]byte, error) {
	token := make([]byte, length)
	masked := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return masked, err
	}
	l := len(mask)
	for k, v := range token {
		index := int(v) % l
		masked[k] = mask[index]
	}
	return masked, nil
}

func NewRandMaskedBytes(mask []byte, length int, origin []byte) ([]byte, error) {
	for {
		token, err := RandMaskedBytes(mask, length)
		if err != nil {
			return token, err
		}
		if bytes.Compare(token, origin) != 0 {
			return token, nil
		}
	}
}
