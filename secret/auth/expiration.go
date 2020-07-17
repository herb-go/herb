package auth

import (
	"time"
)

type Expiration time.Time

func (e *Expiration) NotExpiratd() bool {
	if e == ExpirationNever {
		return true
	}
	return time.Time(*e).After(time.Now())
}

var ExpirationNever *Expiration = nil
