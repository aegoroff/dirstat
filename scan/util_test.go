package scan

import (
	"errors"
	"testing"
)

type errCloser struct{}

func (e *errCloser) Close() error {
	return errors.New("new error")
}

func Test_Close_ThatReturnsError(t *testing.T) {
	// Arrange
	ec := &errCloser{}

	// Act
	Close(ec)

	// Assert
}

func Test_Close_Nil(t *testing.T) {
	// Arrange

	// Act
	Close(nil)

	// Assert
}
