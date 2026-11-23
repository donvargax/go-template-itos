// The tests are in the slice's own package: a slice imports no other slice,
// and depguard cannot tell a slice's own import from another's.
package command

import (
	"strings"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/cli"
	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
	"github.com/itos-corp/go-template-itos/internal/engine/counting/port/porttest"
)

// A run's ends, so a test reads what the person or the script would.
type ends struct {
	stdout, stderr strings.Builder
}

// uiOf is a UI writing to ends.
func uiOf(e *ends) *cli.UI {
	return &cli.UI{In: strings.NewReader(""), Stdout: &e.stdout, Stderr: &e.stderr, Classifiers: []cli.Classifier{Classify}}
}

// count runs the command over the fake's files named in data. Standard input
// is the fake's - key: the command hands the argument to the port as given,
// and turning - into a reader is the stdin adapter's, held by its own tests
// and the scenarios.
func count(t *testing.T, e *ends, file string, data map[string][]byte, asJSON bool) int {
	t.Helper()
	return CLI{File: file, JSON: asJSON}.Run(t.Context(), uiOf(e), Handler{Files: porttest.Files{Data: data}})
}

// A person reads the counts and the path as given on stdout, and is told to
// use --json on stderr: the two are the plain output's whole contract here.
func TestPlainOutputIsTheCountsAndTheHint(t *testing.T) {
	var e ends
	if code := count(t, &e, "notes.txt", map[string][]byte{"notes.txt": []byte("one two three\nfour five")}, false); code != 0 {
		t.Fatalf("Run exit = %d, want 0", code)
	}
	if got, want := e.stdout.String(), "notes.txt: lines 2, words 5\n"; got != want {
		t.Errorf("standard output = %q, want %q", got, want)
	}
	if got, want := e.stderr.String(), "go-template-itos: for scripts, use --json\n"; got != want {
		t.Errorf("error output = %q, want %q", got, want)
	}
}

// The counts are never the input's own text: a file holding a secret must
// not reach the output.
func TestPlainOutputNeverPrintsWhatItRead(t *testing.T) {
	var e ends
	if code := count(t, &e, "secret.txt", map[string][]byte{"secret.txt": []byte("hunter2")}, false); code != 0 {
		t.Fatalf("Run exit = %d, want 0", code)
	}
	if strings.Contains(e.stdout.String(), "hunter2") {
		t.Errorf("standard output prints the input: %q", e.stdout.String())
	}
}

// One object on stdout, one line, and nothing on stderr: --json prints the
// hint nowhere, so a script's output is only what it reads.
func TestJSONOutputIsOneObjectAndNoHint(t *testing.T) {
	var e ends
	if code := count(t, &e, "notes.txt", map[string][]byte{"notes.txt": []byte("one two three\nfour five")}, true); code != 0 {
		t.Fatalf("Run exit = %d, want 0", code)
	}
	want := `{"schema":1,"ok":true,"file":"notes.txt","lines":2,"words":5}` + "\n"
	if got := e.stdout.String(); got != want {
		t.Errorf("standard output = %q, want %q", got, want)
	}
	if e.stderr.Len() != 0 {
		t.Errorf("error output = %q, want none", e.stderr.String())
	}
}

// Standard input is named - in the object too, so a script needs no second
// rule to tell where the counts came from.
func TestJSONOutputNamesStandardInputAsTheArgument(t *testing.T) {
	var e ends
	if code := count(t, &e, "-", map[string][]byte{"-": []byte("one two three")}, true); code != 0 {
		t.Fatalf("Run exit = %d, want 0", code)
	}
	want := `{"schema":1,"ok":true,"file":"-","lines":1,"words":3}` + "\n"
	if got := e.stdout.String(); got != want {
		t.Errorf("standard output = %q, want %q", got, want)
	}
}

// A failure is the domain's, given its code by the UI: a missing file is
// the environment's, exit 82, and the object says so with the rule and what
// to do about it.
func TestAMissingFileIsTheEnvironments(t *testing.T) {
	var e ends
	if code := count(t, &e, "absent.txt", nil, true); code != 82 {
		t.Fatalf("Run exit = %d, want 82", code)
	}
	want := `{"schema":1,"ok":false,"problems":[{"exit":82,` +
		`"file":"absent.txt",` +
		`"fix":"Check the path, or give - to count standard input.",` +
		`"message":"the file absent.txt is not there","rule":"count.input-missing"}]}` + "\n"
	if got := e.stdout.String(); got != want {
		t.Errorf("standard output = %q, want %q", got, want)
	}
	if e.stderr.Len() != 0 {
		t.Errorf("error output = %q, want none under --json", e.stderr.String())
	}
}

// Without --json a failure is a line for a person, and the code is still the
// failure's.
func TestAFailureWithoutJSONIsALineForAPerson(t *testing.T) {
	var e ends
	if code := count(t, &e, "absent.txt", nil, false); code != 82 {
		t.Fatalf("Run exit = %d, want 82", code)
	}
	if e.stdout.Len() != 0 {
		t.Errorf("standard output = %q, want none", e.stdout.String())
	}
	if got, want := e.stderr.String(), "go-template-itos: the file absent.txt is not there. Check the path, or give - to count standard input.\n"; got != want {
		t.Errorf("error output = %q, want %q", got, want)
	}
}

// An unreadable input, a directory among them, is the environment's too, and
// names its own rule rather than the missing one's.
func TestAnUnreadableInputIsTheEnvironments(t *testing.T) {
	var e ends
	files := porttest.Files{Failures: map[string]port.ReadFailure{"adir": porttest.Unreadable("adir")}}
	if code := (CLI{File: "adir", JSON: true}).Run(t.Context(), uiOf(&e), Handler{Files: files}); code != 82 {
		t.Fatalf("Run exit = %d, want 82", code)
	}
	if !strings.Contains(e.stdout.String(), `"rule":"count.input-unreadable"`) {
		t.Errorf("standard output = %q, want the unreadable rule", e.stdout.String())
	}
}

// An argument that names no file is refused where it enters, before any read:
// a usage failure, exit 80, its object naming the argument, and nothing on
// stderr under --json.
func TestAnEmptyArgumentIsRefusedBeforeAnyRead(t *testing.T) {
	var e ends
	if code := count(t, &e, "", map[string][]byte{"": []byte("one")}, true); code != 80 {
		t.Fatalf("Run exit = %d, want 80", code)
	}
	if !strings.Contains(e.stdout.String(), `"rule":"count.input-invalid"`) || !strings.Contains(e.stdout.String(), `"file":""`) {
		t.Errorf("standard output = %q, want the invalid-input problem naming the argument", e.stdout.String())
	}
	if e.stderr.Len() != 0 {
		t.Errorf("error output = %q, want none under --json", e.stderr.String())
	}
}
