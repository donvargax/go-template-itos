// The steps, shared by the two drivers (decision 14). A Driver is where a
// scenario's folders and files live and how go-template-itos runs: the
// binary driver (features_test.go) builds the program and starts it as a
// process, and the in-process driver (cmd/go-template-itos's
// scenarios_test.go) calls cmd's run against the ports' fakes. The steps
// judge go-template-itos by its command line, its exit code and its output
// alone, and import nothing of its code, so one set of scenarios holds under
// either. Harvested from itos's (github.com/donvargax/itos, features/
// steps_test.go) and the owner's itos-template's, its domain's steps left
// behind, and shaped as itos-cc's two drivers are.
package features

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/cucumber/godog"
)

// stampedVersion is the version the binary under test is stamped with, as a
// release's build stamps the release's, so a scenario can tell the stamp
// from the version a build without one says. The scenarios that read it
// carry @binary-only: a run in process has no build to stamp.
const stampedVersion = "1.2.3-features"

// stampedCommit is the commit the binary under test is stamped with, as a
// release's build stamps the commit it was built from.
const stampedCommit = "0123456789abcdef0123456789abcdef01234567"

// Run is one go-template-itos run a step asks for: its words, the
// variables added to the clean environment, and the text it reads from its
// standard input, "" for none.
type Run struct {
	Args  []string
	Env   []Variable
	Stdin string
}

// Variable is one NAME=value a run's environment is given.
type Variable struct{ Name, Value string }

// Result is what a finished run gave.
type Result struct {
	Exit           int
	Stdout, Stderr string
}

// Driver is where a scenario's fixtures live and how go-template-itos runs.
// Setup starts a scenario with nothing in it and Teardown removes what it
// made. Write and Mkdir make a file and a folder at a path relative to the
// scenario's folder, as Run's words name them.
type Driver interface {
	Setup() error
	Teardown()
	Write(name string, data []byte) error
	Mkdir(name string) error
	Run(Run) (Result, error)
}

// Waiter is what a Driver holds if it can run go-template-itos with its
// standard input left open and interrupt it (decision 14: a signal is the
// binary driver's). Start returns without waiting for the run, Interrupt
// ends it once it waits and gives its result, and Stop ends one still under
// way, so no scenario leaves a process behind.
type Waiter interface {
	Start(args []string) error
	Interrupt() (Result, error)
	Stop()
}

// world is one scenario's state.
type world struct {
	d            Driver
	result       Result
	stdin        string // the text the run reads from its standard input, "" for none
	outsideShard bool   // the scenario is another share's (-shard): skipped, nothing set up
}

// Initialize registers the steps and the timing hooks, over d.
func Initialize(sc *godog.ScenarioContext, d Driver) {
	w := &world{d: d}
	clock := timeScenario(sc)
	defer clock.end(sc)
	sc.Before(func(ctx context.Context, s *godog.Scenario) (context.Context, error) {
		w.result, w.stdin = Result{}, ""
		if w.outsideShard = !inShard(s); w.outsideShard {
			return ctx, godog.ErrSkip
		}
		return ctx, timed("(the scenario's set-up)", d.Setup)
	})
	// The scenario's error is godog's already: returned here, it would be
	// reported twice, as the hook's and as the step's.
	sc.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if w.outsideShard {
			return ctx, nil
		}
		_ = timed("(the scenario's clean-up)", func() error {
			if waiter, ok := d.(Waiter); ok {
				waiter.Stop()
			}
			d.Teardown()
			return nil
		})
		return ctx, nil
	})

	// Split before anything else is expanded, so a word with a space in it
	// stays one argument. Quotes also preserve empty words for completion
	// requests.
	sc.Step(`^go-template-itos runs with "([^"]*)"$`, func(args string) error {
		return w.runWith(args, "")
	})
	sc.Step(`^go-template-itos runs with "([^"]*)" and the environment "([^"]*)"$`, func(args, vars string) error {
		return w.runWith(args, vars)
	})
	// A run that waits on its input until something ends it, and the
	// interrupt that does, which only a driver that is a Waiter can give.
	sc.Step(`^go-template-itos runs with "([^"]*)" and its standard input left open$`, w.startWithInputOpen)
	sc.Step(`^it is interrupted while it waits$`, w.interruptWhileWaiting)

	sc.Step(`^it exits with code (\d+)$`, w.exitsWith)
	sc.Step(`^its standard output lists "([^"]*)"$`, func(text string) error {
		if !slices.Contains(outputLines(w.result.Stdout), text) {
			return fmt.Errorf("standard output does not list %q\n%s", text, w.report())
		}
		return nil
	})
	sc.Step(`^its standard output does not list "([^"]*)"$`, func(text string) error {
		if slices.Contains(outputLines(w.result.Stdout), text) {
			return fmt.Errorf("standard output lists %q\n%s", text, w.report())
		}
		return nil
	})
	sc.Step(`^the last line of its standard output is "([^"]*)"$`, func(text string) error {
		lines := outputLines(w.result.Stdout)
		last := strings.TrimSuffix(lines[len(lines)-1], "\r")
		if last != text {
			return fmt.Errorf("the last line of standard output is %q, not %q\n%s", last, text, w.report())
		}
		return nil
	})
	sc.Step(`^its standard output says "([^"]*)"$`, func(text string) error {
		if !strings.Contains(w.result.Stdout, text) {
			return fmt.Errorf("standard output does not say %q\n%s", text, w.report())
		}
		return nil
	})
	sc.Step(`^its standard output does not say "([^"]*)"$`, func(text string) error {
		if strings.Contains(w.result.Stdout, text) {
			return fmt.Errorf("standard output says %q\n%s", text, w.report())
		}
		return nil
	})
	sc.Step(`^its standard output is empty$`, func() error {
		if w.result.Stdout != "" {
			return fmt.Errorf("standard output is not empty\n%s", w.report())
		}
		return nil
	})
	sc.Step(`^its error output is empty$`, func() error {
		if w.result.Stderr != "" {
			return fmt.Errorf("error output is not empty\n%s", w.report())
		}
		return nil
	})
	sc.Step(`^its error output says "([^"]*)"$`, func(text string) error {
		if !strings.Contains(w.result.Stderr, text) {
			return fmt.Errorf("error output does not say %q\n%s", text, w.report())
		}
		return nil
	})
	sc.Step(`^its error output does not say "([^"]*)"$`, func(text string) error {
		if strings.Contains(w.result.Stderr, text) {
			return fmt.Errorf("error output says %q\n%s", text, w.report())
		}
		return nil
	})
	sc.Step(`^the first line of its standard output is "([^"]*)" and the stamped version$`, func(name string) error {
		return w.firstLineIs(name + " " + stampedVersion)
	})
	sc.Step(`^the second line of its standard output is "([^"]*)" and the stamped commit$`, func(word string) error {
		return w.secondLineIs(word + " " + stampedCommit)
	})

	// The --json object, as a script reads it.
	sc.Step(`^its JSON output gives the subject "([^"]*)" as "([^"]*)"$`, w.jsonGivesSubject)
	sc.Step(`^its JSON output names the problem "([^"]*)"$`, w.jsonNamesProblem)
	sc.Step(`^its JSON output gives file "([^"]*)", (\d+) lines and (\d+) words$`, w.jsonGivesCounts)
	sc.Step(`^its JSON output gives each problem a rule, a message and a fix$`, w.jsonProblemsComplete)
	sc.Step(`^its JSON output gives the stamped version and commit$`, w.jsonGivesBuild)
	sc.Step(`^its standard output is one line$`, w.stdoutIsOneLine)

	// The fixtures a scenario asks for. godog takes one parameter per
	// capturing group and none variadic, so the second text is its own
	// group, empty where the scenario names one.
	sc.Step(`^the file "([^"]*)" contains "([^"]*)"(?:, "([^"]*)")?$`, func(name, first, second string) error {
		texts := []string{first}
		if second != "" {
			texts = append(texts, second)
		}
		return w.fileContains(name, texts...)
	})
	sc.Step(`^the file "([^"]*)" holds a line of (\d+) characters$`, w.fileHoldsLine)
	sc.Step(`^standard input holds "([^"]*)"$`, func(text string) error {
		w.stdin = text
		return nil
	})
	sc.Step(`^an empty folder "([^"]*)"$`, w.emptyFolder)
}

// runWith runs go-template-itos with args, split as a scenario wrote them,
// and with each NAME=value of vars, a semicolon-separated list, set in its
// environment.
func (w *world) runWith(args, vars string) error {
	fields, err := commandWords(args)
	if err != nil {
		return err
	}
	run := Run{Args: fields, Stdin: w.stdin}
	for _, one := range strings.Split(vars, ";") {
		if name, value, ok := strings.Cut(one, "="); ok {
			run.Env = append(run.Env, Variable{name, value})
		}
	}
	w.result, err = w.d.Run(run)
	return err
}

// startWithInputOpen starts go-template-itos with args and its standard
// input held open, and returns without waiting for it.
func (w *world) startWithInputOpen(args string) error {
	waiter, ok := w.d.(Waiter)
	if !ok {
		return errors.New("this driver cannot hold a run's standard input open: the scenario is @binary-only")
	}
	fields, err := commandWords(args)
	if err != nil {
		return err
	}
	return waiter.Start(fields)
}

// interruptWhileWaiting interrupts the run startWithInputOpen started once
// it reads its input, and takes what it wrote and its exit code for the
// Then steps.
func (w *world) interruptWhileWaiting() error {
	waiter, ok := w.d.(Waiter)
	if !ok {
		return errors.New("this driver cannot interrupt a run: the scenario is @binary-only")
	}
	var err error
	w.result, err = waiter.Interrupt()
	return err
}

// commandWords splits a scenario's command line without involving a shell.
// Quotes group words and preserve an explicitly empty argument, as the
// completion protocol needs for a new token. Backslashes remain literal, so
// Windows paths survive the harness.
func commandWords(line string) ([]string, error) {
	var words []string
	var word strings.Builder
	var quote rune
	started := false
	for _, r := range line {
		if quote != 0 {
			if r == quote {
				quote = 0
				continue
			}
			word.WriteRune(r)
			started = true
			continue
		}
		switch r {
		case '\'', '"':
			quote = r
			started = true
		case ' ', '\t', '\n', '\r':
			if started {
				words = append(words, word.String())
				word.Reset()
				started = false
			}
		default:
			word.WriteRune(r)
			started = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unclosed quote in scenario command %q", line)
	}
	if started {
		words = append(words, word.String())
	}
	return words, nil
}

// Fixture steps.

// fileContains writes name in the scenario's folder, one line per text the
// scenario named, each ending in the newline a text file's last line has.
func (w *world) fileContains(name string, texts ...string) error {
	if len(texts) == 0 || texts[0] == "" {
		return fmt.Errorf("the file %q is given no text to contain", name)
	}
	return w.d.Write(name, []byte(strings.Join(texts, "\n")+"\n"))
}

// fileHoldsLine writes name holding one line of count characters, long
// enough to beat a token limit a reader might have, and no other line.
func (w *world) fileHoldsLine(name string, count int) error {
	return w.d.Write(name, []byte(strings.Repeat("x", count)+"\n"))
}

// emptyFolder makes a folder in the scenario's folder.
func (w *world) emptyFolder(name string) error {
	return w.d.Mkdir(name)
}

// Then steps.

func (w *world) report() string {
	return fmt.Sprintf("exit %d\n--- stdout\n%s--- stderr\n%s", w.result.Exit, w.result.Stdout, w.result.Stderr)
}

func outputLines(output string) []string {
	return strings.Split(strings.TrimSuffix(output, "\n"), "\n")
}

func (w *world) exitsWith(code int) error {
	if w.result.Exit != code {
		return fmt.Errorf("go-template-itos exited %d, not %d\n%s", w.result.Exit, code, w.report())
	}
	return nil
}

// secondLineIs is whether standard output's second line is text, a missing
// line read as an empty one and a line ending in \r\n as one in \n.
func (w *world) secondLineIs(text string) error {
	_, rest, _ := strings.Cut(w.result.Stdout, "\n")
	second, _, _ := strings.Cut(rest, "\n")
	if second = strings.TrimSuffix(second, "\r"); second != text {
		return fmt.Errorf("the second line of standard output is %q, not %q\n%s", second, text, w.report())
	}
	return nil
}

// firstLineIs is whether standard output's first line is text, a line ending
// in \r\n read as one in \n.
func (w *world) firstLineIs(text string) error {
	first, _, _ := strings.Cut(w.result.Stdout, "\n")
	if first = strings.TrimSuffix(first, "\r"); first != text {
		return fmt.Errorf("the first line of standard output is %q, not %q\n%s", first, text, w.report())
	}
	return nil
}

// The JSON steps read standard output as the one object --json prints.

// json is the object standard output holds, one line of JSON as --json
// prints it.
func (w *world) json() (map[string]any, error) {
	var object map[string]any
	if err := json.Unmarshal([]byte(w.result.Stdout), &object); err != nil {
		return nil, fmt.Errorf("standard output is not one JSON object: %v\n%s", err, w.report())
	}
	return object, nil
}

// jsonNamesProblem is whether the object's problems list one with that rule.
func (w *world) jsonNamesProblem(rule string) error {
	object, err := w.json()
	if err != nil {
		return err
	}
	problems, _ := object["problems"].([]any)
	for _, one := range problems {
		problem, _ := one.(map[string]any)
		if problem["rule"] == rule {
			return nil
		}
	}
	return fmt.Errorf("no problem of the JSON output has the rule %q\n%s", rule, w.report())
}

// jsonGivesSubject is whether the object's one problem holds the key, beside
// its rule, message and fix, as the value: what the problem is about, which a
// script reads without parsing a sentence.
func (w *world) jsonGivesSubject(key, value string) error {
	object, err := w.json()
	if err != nil {
		return err
	}
	problems, _ := object["problems"].([]any)
	for _, one := range problems {
		problem, _ := one.(map[string]any)
		if problem[key] == value {
			return nil
		}
	}
	return fmt.Errorf("no problem of the JSON output has %s %q\n%s", key, value, w.report())
}

// jsonGivesCounts is whether the success object's file is file, its lines
// the count and its words the count.
func (w *world) jsonGivesCounts(file string, lines, words int) error {
	object, err := w.json()
	if err != nil {
		return err
	}
	if object["ok"] != true {
		return fmt.Errorf("the JSON output's ok is %v, not true\n%s", object["ok"], w.report())
	}
	for _, want := range []struct {
		key  string
		got  any
		want any
	}{
		{"file", object["file"], file},
		{"lines", numberOf(object["lines"]), lines},
		{"words", numberOf(object["words"]), words},
	} {
		if want.got != want.want {
			return fmt.Errorf("the JSON output's %s is %v, not %v\n%s", want.key, want.got, want.want, w.report())
		}
	}
	return nil
}

// jsonGivesBuild is whether the success object of version --json is schema
// 1 and ok, and gives the version and the commit the harness stamped.
func (w *world) jsonGivesBuild() error {
	object, err := w.json()
	if err != nil {
		return err
	}
	for _, want := range []struct {
		key  string
		want any
	}{
		{"schema", 1.0},
		{"ok", true},
		{"version", stampedVersion},
		{"commit", stampedCommit},
	} {
		if got := object[want.key]; got != want.want {
			return fmt.Errorf("the JSON output's %s is %v, not %v\n%s", want.key, got, want.want, w.report())
		}
	}
	return nil
}

// stdoutIsOneLine is whether standard output is exactly one line, ended by
// its newline, as --json prints its object.
func (w *world) stdoutIsOneLine() error {
	if strings.Count(w.result.Stdout, "\n") != 1 || !strings.HasSuffix(w.result.Stdout, "\n") {
		return fmt.Errorf("standard output is not one line\n%s", w.report())
	}
	return nil
}

// jsonProblemsComplete is whether the object's problems all carry a rule, a
// message and a fix, and there is one.
func (w *world) jsonProblemsComplete() error {
	object, err := w.json()
	if err != nil {
		return err
	}
	problems, _ := object["problems"].([]any)
	if len(problems) == 0 {
		return fmt.Errorf("the JSON output lists no problem\n%s", w.report())
	}
	for i, one := range problems {
		problem, _ := one.(map[string]any)
		for _, key := range []string{"rule", "message", "fix"} {
			if text, _ := problem[key].(string); text == "" {
				return fmt.Errorf("problem %d of the JSON output has no %s\n%s", i+1, key, w.report())
			}
		}
	}
	return nil
}

// numberOf is the number a decoded JSON value holds, a float64 as encoding/
// json decodes every number into.
func numberOf(value any) int {
	number, _ := value.(float64)
	return int(number)
}
