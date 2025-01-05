package errors

import "fmt"

type IllegalFileName struct {
	Filename string
}

type NoUploadRequest struct {
	FileId string
}

func (e *NoUploadRequest) Error() string {
	return fmt.Sprintf("No file upload request was made for: '%s'", e.FileId)
}

func (e *IllegalFileName) Error() string {
	return fmt.Sprintf("Illegal file name: '%s'", e.Filename)
}
