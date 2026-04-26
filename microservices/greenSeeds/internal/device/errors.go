package device

import "errors"

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrDeviceBusy      = errors.New("device busy: active session in progress")
)
