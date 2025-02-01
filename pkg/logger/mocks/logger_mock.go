// /pkg/logger/mocks/logger_mock.go
package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Info(prefix, msg string) {
	m.Called(prefix, msg)
}

func (m *Logger) Debug(prefix, msg string) {
	m.Called(prefix, msg)
}

func (m *Logger) Fatal(prefix, msg string) {
	m.Called(prefix, msg)
}

func (m *Logger) Response(prefix, status, msg string) {
	m.Called(prefix, status, msg)
}
