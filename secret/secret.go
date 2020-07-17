package secret

type Blob []byte

type Secret interface {
	SecretBlob() (Blob, error)
}

type Key []byte

func (k Key) SecretData() (Blob, error) {
	return Blob(k), nil
}

type ID string
