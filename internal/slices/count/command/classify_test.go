package command

import (
	"errors"
	"fmt"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/cli"
	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
)

// Every kind the counting port's sealed set has is read by a scenario by its
// exit code and its --json rule. Both are the contract a script relies on,
// and they are checked here whole, with the subject a script reads beside
// them: a kind whose code or rule changed would break every script that
// reads it, and no other check here would say so.
func TestEachReadFailureKindHasItsCodeRuleSentenceFixAndSubject(t *testing.T) {
	for _, c := range []struct {
		failure port.ReadFailure
		code    int
		rule    string
		message string
		fix     string
		subject any
	}{
		{&port.Missing{Path: "absent.txt"}, cli.CodeEnvironment, "count.input-missing",
			"the file absent.txt is not there", "Check the path, or give - to count standard input.", "absent.txt"},
		{&port.Unreadable{Path: "adirectory"}, cli.CodeEnvironment, "count.input-unreadable",
			"the file adirectory cannot be read as text", "Give a file this program can read, not a folder.", "adirectory"},
		{&port.Cancelled{Path: "-"}, cli.CodeInterrupted, "interrupted",
			"interrupted before it finished", "Run it again and let it finish.", nil},
	} {
		problems, ok := Classify(fmt.Errorf("wrapped: %w", c.failure))
		if !ok || len(problems) != 1 {
			t.Fatalf("%T: claimed %v, problems %v", c.failure, ok, problems)
		}
		p := problems[0]
		if p.Code != c.code || p.Rule != c.rule || p.Message != c.message || p.Fix != c.fix || p.Subject("file") != c.subject {
			t.Errorf("%T: %+v with file %v", c.failure, p, p.Subject("file"))
		}
	}
}

// A kind of the set no case names is a bug, the 100 of an unclassified one,
// and an error of none of the set is not this slice's to claim.
func TestAKindNoCaseNamesIsInternalAndAForeignErrorIsNotClaimed(t *testing.T) {
	problems, ok := Classify(unnamed{})
	if !ok || len(problems) != 1 || problems[0].Code != cli.CodeInternal || problems[0].Rule != "internal" {
		t.Errorf("an unnamed kind: %v, %v", problems, ok)
	}
	if problems, ok := Classify(errors.New("no space")); ok || problems != nil {
		t.Errorf("a foreign error: %v, %v", problems, ok)
	}
}

// unnamed is a failure of the port's set that Classify names no case for.
type unnamed struct{ port.ReadFailure }

func (unnamed) Error() string { return "unnamed" }
