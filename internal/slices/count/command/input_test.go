package command

import (
	"errors"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/cli"
)

// An argument that names a file is taken as given, - for standard input and a
// name that starts with a dash among them: a lossless check never changes it.
func TestAnArgumentThatNamesAFileIsTakenAsGiven(t *testing.T) {
	for _, arg := range []string{"notes.txt", "-", "-draft.txt", " spaced name ", "dir/a b.txt", "../up"} {
		input, err := NewInput(arg)
		if err != nil || input.Path() != arg {
			t.Errorf("NewInput(%q) = %q, %v", arg, input.Path(), err)
		}
	}
}

// An argument that names no file is refused as given, as the problem a script
// reads: usage, exit 80, its rule and the argument as its subject. Nothing is
// repaired, so no refusal holds a value.
func TestAnArgumentThatNamesNoFileIsRefusedNamingIt(t *testing.T) {
	for _, c := range []struct{ arg, message string }{
		{"", "the file argument is empty"},
		{"a\x00b", `the file argument "a\x00b" holds a NUL byte`},
	} {
		input, err := NewInput(c.arg)
		var p *cli.Problem
		if !errors.As(err, &p) {
			t.Fatalf("NewInput(%q) error = %v, want a problem", c.arg, err)
		}
		if p.Code != cli.CodeUsage || p.Rule != "count.input-invalid" || p.Message != c.message ||
			p.Fix == "" || p.Subject("file") != c.arg || input.Path() != "" {
			t.Errorf("NewInput(%q) = %+v and %q", c.arg, p, input.Path())
		}
	}
}
