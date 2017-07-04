package security

import "bytes"

/**
 *   Reference http://blog.studygolang.com/167.html
 */
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	d := make([]byte, padding+len(data))
	copy(d, data)
	copy(d[len(data):], padtext)
	return d

}

/**
 *  Reference http://blog.studygolang.com/167.html
 */
func PKCS7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	d := make([]byte, length-unpadding)
	copy(d, data)
	return d
}
