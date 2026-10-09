// Package porttest provides in-memory fakes for counting-domain tests.
package porttest

import (
	"errors"

	"github.com/itos-corp/go-template-itos/internal/count/port"
)

// Files holds test inputs and configured read failures in memory.
type Files struct {
	Data     map[string][]byte
	Failures map[string]error
}

var _ port.Files = Files{}

// Read returns a copy of the configured file, or its configured failure.
func (f Files) Read(path string) ([]byte, error) {
	if err, ok := f.Failures[path]; ok {
		return nil, err
	}
	data, ok := f.Data[path]
	if !ok {
		return nil, &port.Missing{Path: path}
	}
	return append([]byte(nil), data...), nil
}

// Unreadable returns a configured unreadable-file failure.
func Unreadable(path string) error { return &port.Unreadable{Path: path} }

// Failure is a convenience for a fixed in-memory read failure.
func Failure(message string) error { return errors.New(message) }
