// Package disk adapts read-only operating-system file access to the counting
// domain's Files port.
package disk

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/itos-corp/go-template-itos/internal/counting/port"
)

// Files reads input files from the local filesystem.
type Files struct{}

var _ port.Files = Files{}

// Read reads path without modifying or executing its contents and translates
// operating-system errors into the port's typed failures. It does not act on
// ctx: what a cancelled read returns is not yet settled.
func (Files) Read(ctx context.Context, path string) ([]byte, port.ReadFailure) {
	data, err := os.ReadFile(path) //nolint:gosec // caller selects the input path through the Files port
	if err != nil {
		return nil, classifyReadError(path, err)
	}
	return data, nil
}

func classifyReadError(path string, err error) port.ReadFailure {
	if errors.Is(err, fs.ErrNotExist) {
		return &port.Missing{Path: path}
	}
	return &port.Unreadable{Path: path}
}
