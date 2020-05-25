package middlewarefactory

import "errors"

var ErrFactoryNotRegistered = errors.New("factory not registered")
var ErrConditionFactoryNotRegistered = errors.New("condition factory not registered")
var ErrFactoryRegistered = errors.New("factory registered")
var ErrConditionFactoryRegistered = errors.New("condition factory registered")
