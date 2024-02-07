package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	expected := "test"
	require.Equal(t, &expected, String("test"))
}

func TestBool(t *testing.T) {
	expected := true
	require.Equal(t, &expected, Bool(true))
}

func TestInt(t *testing.T) {
	expected := 1
	require.Equal(t, &expected, Int(1))
}
