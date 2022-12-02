package error

import (
	"errors"
	"fmt"
)

type NotUniqueRecordError struct {
	URL string
}

func (e *NotUniqueRecordError) Error() string {
	return fmt.Sprintf("%s already exist in db", e.URL)
}

var (
	ErrDeleteLink   = errors.New("delete link")
	ErrLinkNotFound = errors.New("link not found")
)
