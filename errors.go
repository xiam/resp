package resp

import (
	"errors"
)

var (
	ErrInvalidResponse  = errors.New(`resp: Invalid response.`)
	ErrInvalidDelimiter = errors.New(`resp: Failed to get limits.`)
)
