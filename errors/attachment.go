package errors

import (
	"errors"
)

var ErrAttachmentIsNotAFile = errors.New("file is a dir")
