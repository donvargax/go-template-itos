package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// fault is a failure of a sealed set of the test's own, which the fake
// classifier below claims: the UI imports no slice, so it is held to its
// contract with a set that is none.
type fault struct{ what string }

func (f *fault) Error() string { return "fault " + f.what }

// classifyFaults is a slice's classifier as cmd hands one over: it claims a
// fault and gives it a code, a rule, a fix and a subject.
func classifyFaults(code int) Classifier {
	return func(err error) ([]*Problem, bool) {
		var f *fault
		if !errors.As(err, &f) {
			return nil, false
		}
		return []*Problem{NewProblem(code, "fake.fault", "the fault "+f.what, "Mend it.").With("what", f.what)}, true
	}
}

// run is err reported by a UI of the classifiers, as --json or as lines: the
// exit code, what was printed on stdout and on stderr.
func run(err error, asJSON bool, classifiers ...Classifier) (code int, stdout, stderr string) {
	var out, errs strings.Builder
	ui := &UI{Stdout: &out, Stderr: &errs, Classifiers: classifiers}
	return ui.Fail(err, asJSON), out.String(), errs.String()
}

// A classifier's problem is the contract a script reads: its exit code and
// its object, whole, with its subject's keys beside rule, message, fix and
// exit. Under --json nothing is written on stderr.
func TestAClassifiedFailureIsAnObjectOfItsRuleMessageFixExitAndSubject(t *testing.T) {
	code, stdout, stderr := run(&fault{"a"}, true, classifyFaults(CodeEnvironment))
	if code != CodeEnvironment {
		t.Errorf("exit = %d, want %d", code, CodeEnvironment)
	}
	want := `{"schema":1,"ok":false,"problems":[{"exit":82,"fix":"Mend it.","message":"the fault a","rule":"fake.fault","what":"a"}]}` + "\n"
	if stdout != want {
		t.Errorf("--json is %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want nothing under --json", stderr)
	}
}

// Without --json each problem is a line for people, its message and then its
// fix, and nothing is printed on stdout.
func TestAPlainFailureKeepsItsFixOnStderr(t *testing.T) {
	code, stdout, stderr := run(&fault{"a"}, false, classifyFaults(CodeEnvironment))
	if code != CodeEnvironment || stdout != "" {
		t.Errorf("exit = %d and stdout = %q", code, stdout)
	}
	if want := "go-template-itos: the fault a. Mend it.\n"; stderr != want {
		t.Errorf("stderr = %q, want %q", stderr, want)
	}
}

// A problem without a fix is its message alone, never a dangling sentence.
func TestAProblemWithNoFixIsItsMessageAlone(t *testing.T) {
	if got, want := NewProblem(CodeNo, "r", "no", "").plain(), "no"; got != want {
		t.Errorf("plain = %q, want %q", got, want)
	}
}

// An error no classifier claims is a defect of ours: 100, and never some
// other code a child process might have returned.
func TestAnErrorNoClassifierClaimsIsInternal(t *testing.T) {
	cause := fmt.Errorf("cannot write the report: %w", errors.New("no space"))
	code, stdout, stderr := run(cause, false, classifyFaults(CodeEnvironment))
	if code != CodeInternal {
		t.Fatalf("exit = %d, want %d", code, CodeInternal)
	}
	if stdout != "" {
		t.Errorf("standard output = %q, want none without --json", stdout)
	}
	if !strings.HasPrefix(stderr, "go-template-itos: cannot write the report: no space. This is a bug") {
		t.Errorf("stderr = %q", stderr)
	}
}

// A problem is itself, however an error wraps it, and no classifier is asked.
func TestAProblemIsItsOwnReport(t *testing.T) {
	p := NewProblem(CodeUsage, "usage", "bad", "Mend it.")
	if code, _, _ := run(fmt.Errorf("wrapped: %w", p), false); code != CodeUsage {
		t.Errorf("exit = %d, want %d", code, CodeUsage)
	}
}

// The first classifier that claims an error gives its problems; a later one
// is not asked.
func TestTheFirstClassifierThatClaimsAnErrorGivesItsProblems(t *testing.T) {
	code, _, _ := run(&fault{"a"}, false, func(error) ([]*Problem, bool) { return nil, false },
		classifyFaults(CodeTemporary), classifyFaults(CodeEnvironment))
	if code != CodeTemporary {
		t.Errorf("exit = %d, want %d", code, CodeTemporary)
	}
}

// Two failures joined are each a problem of their own, in one object and a
// line each, and the code is the first of the precedence they have: one
// joined error never hides another.
func TestJoinedFailuresAreEachAProblemUnderTheFirstCodeOfThePrecedence(t *testing.T) {
	joined := errors.Join(&fault{"a"}, errors.Join(&fault{"b"}, errors.New("a bug")))
	code, stdout, stderr := run(joined, true, classifyFaults(CodeEnvironment))
	if code != CodeInternal {
		t.Fatalf("exit = %d, want %d: a bug before a missing requirement", code, CodeInternal)
	}
	if got, want := strings.Count(stdout, `"rule":`), 3; got != want {
		t.Errorf("--json holds %d problems, want %d\n%s", got, want, stdout)
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want nothing under --json", stderr)
	}
	_, _, stderr = run(joined, false, classifyFaults(CodeEnvironment))
	if got, want := strings.Count(stderr, "go-template-itos: "), 3; got != want {
		t.Errorf("stderr holds %d lines, want %d\n%s", got, want, stderr)
	}
}

// Precedence is the order of itos-cc's: the first code of it a run has, not
// the highest, and a code the list lacks, a defect, is 100.
func TestExitCodeIsTheFirstOfThePrecedence(t *testing.T) {
	of := func(codes ...int) int {
		var problems []*Problem
		for _, c := range codes {
			problems = append(problems, NewProblem(c, "r", "m", "f"))
		}
		return ExitCode(problems)
	}
	if got := of(); got != 0 {
		t.Errorf("no problem is exit %d, want 0", got)
	}
	if got := of(7); got != CodeInternal {
		t.Errorf("a code the list lacks is exit %d, want %d", got, CodeInternal)
	}
	order := []int{CodeInterrupted, CodeUsage, CodeUnknownCommand, CodeInternal, CodeEnvironment, CodeTemporary, CodeNo}
	for i, first := range order {
		for _, later := range order[i:] {
			if got := of(later, first); got != first {
				t.Errorf("%d with %d is exit %d, want %d", later, first, got, first)
			}
		}
	}
}

// The codes a script reads are fixed, and 2, 3 and 70 are none of them.
func TestTheExitCodesAreDecision13s(t *testing.T) {
	for code, want := range map[int]int{CodeNo: 1, CodeTemporary: 75, CodeUsage: 80, CodeUnknownCommand: 81,
		CodeEnvironment: 82, CodeInternal: 100, CodeInterrupted: 130} {
		if code != want {
			t.Errorf("code %d, want %d", code, want)
		}
	}
	if len(Precedence) != 7 {
		t.Errorf("Precedence holds %d codes, want 7", len(Precedence))
	}
}

// A subject never hides the problem's own keys: rule, message, fix and exit
// are written over it, and the subject's key is read back as set.
func TestASubjectNeverHidesTheProblemsOwnKeys(t *testing.T) {
	p := NewProblem(CodeUsage, "usage", "bad", "Mend it.").With("exit", 1).With("rule", "x").With("flag", "--bogus")
	raw, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["exit"] != float64(CodeUsage) || got["rule"] != "usage" || got["flag"] != "--bogus" {
		t.Errorf("object = %s", raw)
	}
	if p.Subject("flag") != "--bogus" || p.Subject("absent") != nil || p.Error() != "bad" {
		t.Errorf("subject, error = %v, %v, %q", p.Subject("flag"), p.Subject("absent"), p.Error())
	}
}

// A problem's text is as written: <, > and & stay as they are, so an input's
// name is copied from the object as a person wrote it.
func TestAProblemsTextIsNotHTMLEscaped(t *testing.T) {
	var stdout strings.Builder
	(&UI{Stdout: &stdout}).JSON(NewProblem(CodeUsage, "usage", "a<b&c>", "Mend it."))
	if !strings.Contains(stdout.String(), `"a<b&c>"`) {
		t.Errorf("object = %s", stdout.String())
	}
}

// A usage failure of a command that is not kong's is exit 80, names the
// command as it was written, and says what to do next.
func TestUsageNamesTheCommandAndTheWayOut(t *testing.T) {
	var stdout, stderr strings.Builder
	ui := &UI{Stdout: &stdout, Stderr: &stderr}
	if code := ui.Usage("completion", errors.New("unsupported shell"), false); code != CodeUsage {
		t.Fatalf("Usage exit = %d, want %d", code, CodeUsage)
	}
	want := "go-template-itos: completion: unsupported shell. Run 'go-template-itos --help' for the commands and the flags each takes.\n"
	if got := stderr.String(); got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
}

// A line that names no command says what was wrong and nothing more.
func TestUsageOfALineNamingNoCommandSaysOnlyWhatWasWrong(t *testing.T) {
	var stdout, stderr strings.Builder
	ui := &UI{Stdout: &stdout, Stderr: &stderr}
	ui.Usage("", errors.New("unknown flag --bogus"), false)
	if got := stderr.String(); !strings.HasPrefix(got, "go-template-itos: unknown flag --bogus. ") {
		t.Errorf("stderr = %q", got)
	}
}

// The interrupt is a problem of its own rule, with the code shells report.
func TestInterruptedIsExit130(t *testing.T) {
	p := Interrupted()
	if p.Code != CodeInterrupted || p.Rule != "interrupted" || p.Fix == "" {
		t.Errorf("Interrupted = %+v", p)
	}
}

// A writer that cannot hold what it is given is a line on stderr, not a
// silent loss: --json's contract is what a script reads.
func TestAFailureToWriteJSONIsALineOnStderr(t *testing.T) {
	var stderr strings.Builder
	ui := &UI{Stdout: failingWriter{}, Stderr: &stderr}
	ui.JSON(struct {
		Schema int `json:"schema"`
	}{1})
	if got, want := stderr.String(), "go-template-itos: the disk is full\n"; got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
}

// failingWriter is a writer that cannot be written, whatever a machine makes
// one: no disk, no permission bit.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("the disk is full") }

// No failure is no problems, and a subject that cannot be written as JSON is
// an error the caller sees, never a half object.
func TestNoFailureIsNoProblemsAndAnUnwritableSubjectIsAnError(t *testing.T) {
	if got := (&UI{}).Problems(nil); got != nil {
		t.Errorf("Problems(nil) = %v, want none", got)
	}
	if _, err := json.Marshal(NewProblem(CodeUsage, "usage", "bad", "Mend it.").With("k", make(chan int))); err == nil {
		t.Error("a subject of a channel was written as JSON")
	}
}

// With no candidate near the word, nothing is suggested.
func TestDidYouMeanOfNoNearCandidateIsEmpty(t *testing.T) {
	if got := DidYouMean("zzzzzz", []string{"count", "version"}, func(s string) string { return s }); got != "" {
		t.Errorf("DidYouMean = %q, want none", got)
	}
}

// Several near candidates are all named, in their order.
func TestDidYouMeanOfSeveralNearCandidatesNamesEachOfThem(t *testing.T) {
	got := DidYouMean("coun", []string{"count", "counts", "version"}, func(s string) string { return "<" + s + ">" })
	if want := "Did you mean one of <count>, <counts>? "; got != want {
		t.Errorf("DidYouMean = %q, want %q", got, want)
	}
}
