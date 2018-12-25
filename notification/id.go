package notification

import uuid "github.com/satori/go.uuid"

type IDGenerator func() (string, error)

var DefaultIDGenerator = func() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), err
}
