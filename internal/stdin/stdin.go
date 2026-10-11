// Package stdin is the Files adapter for standard input: the reader behind
// the argument -, and the Files it was given for any other path.
//
// The domain knows no paths of the command line and no standard input: it is
// handed bytes through the Files port, and this adapter is the one place the
// - of a command line becomes the reader behind it. It imports nothing of
// ours but the port: cmd gives it the disk adapter as its fallback, so no
// adapter imports another.
package stdin

import (
	"context"
	"io"

	"github.com/itos-corp/go-template-itos/internal/counting/port"
)

// Dash is the argument that names standard input.
const Dash = "-"

// Files reads all of In where the path is Dash, and any other path through
// Fallback.
type Files struct {
	Fallback port.Files
	In       io.Reader
}

var _ port.Files = Files{}

// Read reads path, and refuses any path but Dash through In, so a caller
// that names a file reads it through the fallback and one that names - reads
// a reader. A reader that fails is the environment's, as a file that cannot
// be read is. It hands ctx to the fallback and does not act on it: what a
// cancelled read returns is not yet settled.
func (f Files) Read(ctx context.Context, path string) ([]byte, port.ReadFailure) {
	if path != Dash {
		return f.Fallback.Read(ctx, path)
	}
	data, err := io.ReadAll(f.In)
	if err != nil {
		return nil, &port.Unreadable{Path: path}
	}
	return data, nil
}
