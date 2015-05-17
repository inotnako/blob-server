package storage

import (
	"io"
)

type IdRequestError interface {
	NotFound() bool
	IllFormed() bool
	Error() string
}

type Storage interface {
	Post(reader io.Reader) (string, error)
	Get(id string, writer io.Writer) IdRequestError
	GetList() ([]string, error)
	Delete(id string) IdRequestError
}
