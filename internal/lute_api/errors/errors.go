package errors

import "fmt"

type IllegalFileName struct {
	Filename string
}

func (e *IllegalFileName) Error() string {
	return fmt.Sprintf("Illegal file name: '%s'", e.Filename)
}
