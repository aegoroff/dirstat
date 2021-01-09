package out

import (
	"bytes"
)

type memoryEnvironment struct {
	*stringEnvironment
	buffer *bufferClosable
}

type bufferClosable struct {
	*bytes.Buffer
}

func (*bufferClosable) Close() error {
	return nil
}

func (m *memoryEnvironment) String() string {
	return m.buffer.String()
}

// NewMemoryEnvironment creates new memory PrintEnvironment implementation
func NewMemoryEnvironment() StringEnvironment {
	buffer := &bufferClosable{
		bytes.NewBufferString(""),
	}
	return &memoryEnvironment{
		stringEnvironment: newStringEnvironment(buffer),
		buffer:            buffer,
	}
}
