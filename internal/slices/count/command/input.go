package command

import (
	"fmt"
	"strings"

	"github.com/itos-corp/go-template-itos/internal/cli"
)

// Input is the file argument, checked where it enters: the one thing the
// command line says to count, a path or - for standard input. Its field is
// unexported and NewInput the one constructor that checks it, so holding an
// Input proves the argument was checked; the Query, the Handler and the port
// beneath them trust it and check it no more. The zero Input was never
// checked, and nothing here makes one.
type Input struct{ path string }

// NewInput checks arg, the file argument as the command line gave it, and
// refuses what names no file: an empty one, and one holding a NUL byte, which
// no system takes in a path. A refusal is the problem count.input-invalid,
// exit 80, naming the argument as given; nothing is repaired, so an empty
// argument never becomes standard input nor the working folder.
func NewInput(arg string) (Input, error) {
	switch {
	case arg == "":
		return Input{}, invalid(arg, "the file argument is empty")
	case strings.ContainsRune(arg, 0):
		return Input{}, invalid(arg, fmt.Sprintf("the file argument %q holds a NUL byte", arg))
	}
	return Input{path: arg}, nil
}

// invalid is the problem of arg, an argument that names no file.
func invalid(arg, message string) error {
	return cli.NewProblem(cli.CodeUsage, "count.input-invalid", message,
		"Give the path of a file to count, or - for standard input.").With("file", arg)
}

// Path is the argument as it was given, which the Files port reads and the
// output names.
func (i Input) Path() string { return i.path }
