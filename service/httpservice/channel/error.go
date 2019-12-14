package channel

import "errors"

var ErrChannelUsed = errors.New("channel used")

var ErrChannelNotRegistered = errors.New("channel not registered")

var ErrChannelStarted = errors.New("channel started")

var ErrChannelStopped = errors.New("channel stopped")
