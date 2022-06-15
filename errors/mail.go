package errors

import (
	"errors"
)

var (
	ErrEmailProviderIsDisabled = errors.New("email provider is disabled")
	ErrEmptyAddress            = errors.New("mail address is empty")
	ErrEmptyData               = errors.New("empty data")
)
