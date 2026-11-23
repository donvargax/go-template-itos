// Package cli is the UI every command shares: where a command reads and
// writes, the problems it reports, the lines on stderr, the --json object,
// and the one place a run's problems become its exit code.
//
// It imports nothing of the project's. A slice's errors are a sealed set of
// plain types carrying no exit code, rule ID or wording; cmd hands the UI
// one Classifier per slice, each turning that slice's failures into problems
// in one switch naming every kind. An error no classifier claims is the 100
// of an unclassified bug, never a child process's status passed through.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
)

// UI is a command's ends: what it reads and where it writes, and the
// classifiers that turn its slices' failures into problems. Nothing here
// asks a question, so there is no terminal to know of.
type UI struct {
	In     io.Reader
	Stdout io.Writer
	Stderr io.Writer
	// Classifiers are cmd's, one per slice, tried in turn.
	Classifiers []Classifier
}

// Classifier turns err, a failure of one slice's sealed set, into its
// problems, each with its code, and says whether err was of that set. A
// classifier names every kind of its set in one switch, so gochecksumtype
// refuses it when a kind is added and left out.
type Classifier func(err error) ([]*Problem, bool)

// Fail reports err, a command failing, and returns its exit code: with
// withJSON the failure's object on stdout and nothing on stderr, otherwise
// each problem a line on stderr.
func (u *UI) Fail(err error, withJSON bool) int {
	return u.Emit(withJSON, u.Problems(err))
}

// Usage reports err, a usage error of a command that is not kong's, as Fail
// does: exit 80, the command named first as it was written, so the sentence
// says what was read.
func (u *UI) Usage(command string, err error, withJSON bool) int {
	message := err.Error()
	if command != "" {
		message = command + ": " + message
	}
	return u.Emit(withJSON, []*Problem{NewProblem(CodeUsage, "usage", message,
		"Run 'go-template-itos --help' for the commands and the flags each takes.")})
}

// Emit prints a run's outcome, which is failing, and returns its exit code:
// with asJSON the object {"schema": 1, "ok": false, "problems": […]} on one
// line of stdout and nothing on stderr; otherwise each problem on stderr, as
// "go-template-itos: <message>. <fix>".
func (u *UI) Emit(asJSON bool, problems []*Problem) int {
	if asJSON {
		u.JSON(struct {
			Schema   int        `json:"schema"`
			OK       bool       `json:"ok"`
			Problems []*Problem `json:"problems"`
		}{1, false, problems})
	} else {
		for _, p := range problems {
			u.Line(p.plain())
		}
	}
	return ExitCode(problems)
}

// ExitCode is the exit code of a run with problems: 0 with none, else the
// first code of Precedence one of them has. A problem of a code the list
// lacks is a defect of ours, so with no listed code it is 100.
func ExitCode(problems []*Problem) int {
	for _, code := range Precedence {
		if slices.ContainsFunc(problems, func(p *Problem) bool { return p.Code == code }) {
			return code
		}
	}
	if len(problems) > 0 {
		return CodeInternal
	}
	return 0
}

// Problems is err as the problems it reports. A problem is itself; a
// failure a classifier claims is that classifier's problems; anything else
// is one internal problem, err whole. Errors joined (errors.Join) are each
// their own problems, so one never hides another.
func (u *UI) Problems(err error) []*Problem {
	if err == nil {
		return nil
	}
	var problems []*Problem
	for _, leaf := range leaves(err) {
		problems = append(problems, u.classify(leaf)...)
	}
	return problems
}

// classify is err's problems, from the first that claims it: err itself
// when it is a problem, or a classifier's, else err as an internal one.
func (u *UI) classify(err error) []*Problem {
	var p *Problem
	if errors.As(err, &p) {
		return []*Problem{p}
	}
	for _, c := range u.Classifiers {
		if problems, ok := c(err); ok {
			return problems
		}
	}
	return []*Problem{Internal(err)}
}

// leaves are the errors err joins, each in turn, or err.
func leaves(err error) []error {
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return []error{err}
	}
	var all []error
	for _, e := range joined.Unwrap() {
		all = append(all, leaves(e)...)
	}
	return all
}

// Internal is err as an internal problem, a defect of ours: exit 100.
func Internal(err error) *Problem {
	return NewProblem(CodeInternal, "internal", err.Error(), "This is a bug in go-template-itos; please report it.")
}

// Interrupted is the problem of a command an interrupt cancelled: exit 130,
// what shells report for Ctrl-C (128 and SIGINT's 2), on every system.
func Interrupted() *Problem {
	return NewProblem(CodeInterrupted, "interrupted", "interrupted before it finished", "Run it again and let it finish.")
}

// JSON writes v on stdout as one line of JSON, as --json prints it; a
// failure to is a line on stderr.
func (u *UI) JSON(v any) {
	enc := json.NewEncoder(u.Stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		u.Line(err.Error())
	}
}

// Line writes message on stderr as a line for people, go-template-itos:
// first. A stderr that cannot be written leaves nowhere else to say so.
func (u *UI) Line(message string) {
	_, _ = fmt.Fprintf(u.Stderr, "go-template-itos: %s\n", message)
}
