package command

import (
	"errors"
	"fmt"

	"github.com/itos-corp/go-template-itos/internal/cli"
	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
)

// Classify is the count slice's classifier: it maps the counting port's
// sealed read failures to problems, each with its code, rule and fix, and
// says whether err was one of them. Any other error is not this slice's, so
// the UI tries the next classifier and, with none left, reports it as the
// bug it is.
//
// Each kind below is one the count command's scenarios read: a missing file
// and an unreadable one, a directory among the latter, both the
// environment's, and a read an interrupt cancelled, exit 130 (decision 8).
// A kind of the set that no case names is the 100 of a bug.
func Classify(err error) ([]*cli.Problem, bool) {
	var failure port.ReadFailure
	if !errors.As(err, &failure) {
		return nil, false
	}
	switch e := failure.(type) {
	case *port.Missing:
		return []*cli.Problem{cli.NewProblem(cli.CodeEnvironment, "count.input-missing",
			fmt.Sprintf("the file %s is not there", e.Path),
			"Check the path, or give - to count standard input.").With("file", e.Path)}, true
	case *port.Unreadable:
		return []*cli.Problem{cli.NewProblem(cli.CodeEnvironment, "count.input-unreadable",
			fmt.Sprintf("the file %s cannot be read as text", e.Path),
			"Give a file this program can read, not a folder.").With("file", e.Path)}, true
	case *port.Cancelled:
		return []*cli.Problem{cli.Interrupted()}, true
	}
	return []*cli.Problem{cli.Internal(failure)}, true
}
