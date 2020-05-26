package middlewarefactory

import "errors"

var ErrFactoryNotRegistered = errors.New("factory not registered")
var ErrConditionFactoryNotRegistered = errors.New("condition factory not registered")
