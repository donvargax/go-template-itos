// Package count owns the private example's file-counting behavior.
package count

import "github.com/itos-corp/go-template-itos/internal/count/port"

// Result is the number of logical lines and whitespace-delimited words.
type Result struct {
	Lines int
	Words int
}

// Count counts logical lines and words in data.
func Count(data []byte) Result { return Result{} }

// File counts the contents of path through the Files port.
func File(files port.Files, path string) (Result, error) { return Result{}, nil }
