package utility

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	filename := "file_test.go"
	err := Open(filename, func(file *os.File) error {
		return nil
	})
	require.Nil(t, err)
}
func TestOpen_Error(t *testing.T) {
	filename := "file_test.go"
	err := Open(filename, func(file *os.File) error {
		return assert.AnError
	})
	require.NotNil(t, err)
	require.Error(t, assert.AnError, err)

	filename = "a file that is most likely not there"
	err = Open(filename, func(file *os.File) error {
		return nil
	})
	require.NotNil(t, err)
	require.Error(t, assert.AnError, err)
}

func TestCreate(t *testing.T) {
	filename := "testing123.txt"
	err := Create(filename, func(file *os.File) error {
		return nil
	})
	require.Nil(t, err)

	err = os.Remove("testing123.txt")
	require.Nil(t, err)
}

func TestCreate_Error(t *testing.T) {
	filename := "testing123.txt"
	err := Create(filename, func(file *os.File) error {
		return assert.AnError
	})
	require.NotNil(t, err)
	require.Error(t, assert.AnError, err)

	err = os.Remove("testing123.txt")
	require.Nil(t, err)

	filename = ""
	err = Create(filename, func(file *os.File) error {
		return assert.AnError
	})
	require.NotNil(t, err)
	require.Error(t, assert.AnError, err)
}
