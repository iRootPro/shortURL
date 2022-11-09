package error

import "fmt"

type NotUniqueRecordError struct {
	URL string
}

func (e *NotUniqueRecordError) Error() string {
	return fmt.Sprintf("%s already exist in db", e.URL)
}
