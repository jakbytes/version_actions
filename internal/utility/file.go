package utility

import (
	"errors"
	"os"
)

// Open opens a file using os.Open, passes it to a handler function, and ensures that it is closed after the handler is done.
func Open(filePath string, handler func(*os.File) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	return errors.Join(handler(file), file.Close())
}

func OpenFile(name string, flag int, perm os.FileMode, handler func(*os.File) error) error {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}
	return errors.Join(handler(file), file.Close())
}

// Create creates a file using os.Create, passes it to a handler function, and ensures that it is closed after the handler is done.
func Create(filePath string, handler func(*os.File) error) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	return errors.Join(handler(file), file.Close())
}
