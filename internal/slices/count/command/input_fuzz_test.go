package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/cli"
)

// FuzzNewInput holds the checked input to what no argument may break, whatever the command line
// hands it: an argument naming a file comes back as given and one naming none is refused, as the
// problem a script reads, naming the argument, with no value left to trust. Go runs the seeds with
// every go test; the nightly (tools/bin/nightly-fuzz) is where the fuzzer explores. A failing
// input it writes under testdata/fuzz is kept by committing that file, which every later run
// replays.
func FuzzNewInput(f *testing.F) {
	for _, seed := range []string{
		"", "-", "notes.txt", "-draft.txt", " spaced name ", "dir/a b.txt", "../up",
		"a\x00b", "\x00", "a\xffb", "é\u2003x", "--", "C:\\notes.txt",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, arg string) {
		input, err := NewInput(arg)
		if arg != "" && !strings.ContainsRune(arg, 0) {
			if err != nil || input.Path() != arg {
				t.Fatalf("NewInput(%q) = %q, %v; want the argument as given", arg, input.Path(), err)
			}
			return
		}
		var problem *cli.Problem
		if !errors.As(err, &problem) {
			t.Fatalf("NewInput(%q) error = %v, want a problem", arg, err)
		}
		if problem.Code != cli.CodeUsage || problem.Rule != "count.input-invalid" || problem.Subject("file") != arg || problem.Fix == "" {
			t.Fatalf("NewInput(%q) = %+v, want a usage problem count.input-invalid naming the argument", arg, problem)
		}
		if input.Path() != "" {
			t.Fatalf("NewInput(%q) holds the path %q after a refusal", arg, input.Path())
		}
	})
}
