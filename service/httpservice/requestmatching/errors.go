package requestmatching

import "errors"

//ErrHeaderNotValidated error raised if given header is not validated.
var ErrHeaderNotValidated = errors.New("requestmatching:header is not validated")
