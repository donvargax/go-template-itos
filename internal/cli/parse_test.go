package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/donvargax/kong"
)

// The scenarios hold the unknown command and the unknown flag by their exit
// codes and rules. These hold each kind a command line can fail as, by its
// type, with its rule, code, subject and sentence whole.

// tree is a command line of every shape ParseProblem tells: a command that
// takes a file and flags, a leaf that takes nothing, a group that needs a
// subcommand, a hidden command, and a flag of a choice.
type tree struct {
	Count struct {
		File  string `arg:""`
		Stack string `placeholder:"STACK"`
		Kind  string `enum:"a,b" default:"a"`
		JSON  bool   `negatable:""`
		Num   int
	} `cmd:""`
	Version struct{} `cmd:""`
	Group   struct {
		Sub struct{} `cmd:""`
		Two struct{} `cmd:""`
	} `cmd:""`
	Hidden struct{} `cmd:"" hidden:""`
}

// problem is ParseProblem of the error kong gives for args.
func problem(t *testing.T, args ...string) *Problem {
	t.Helper()
	var l tree
	parser, err := kong.New(&l, append(Flags(), kong.Name("go-template-itos"), kong.Exit(func(int) { t.Fatal("kong exited") }))...)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parser.Parse(args); err == nil {
		t.Fatalf("%v parsed", args)
	}
	return ParseProblem(err)
}

func TestEachKindOfParseFailureIsAProblemOfItsOwnRule(t *testing.T) {
	for _, c := range []struct {
		args    []string
		code    int
		rule    string
		message string
		subject [2]string
	}{
		{[]string{"count", "f", "--bogus"}, CodeUsage, "flags.unknown", "count takes no flag --bogus", [2]string{"flag", "--bogus"}},
		{[]string{"--bogus", "count", "f"}, CodeUsage, "flags.unknown", "go-template-itos takes no flag --bogus", [2]string{"flag", "--bogus"}},
		{[]string{"count", "f", "--json", "--no-json"}, CodeUsage, "flags.repeated", "--json is given more than once", [2]string{"flag", "--json"}},
		{[]string{"count", "f", "--json=yes"}, CodeUsage, "flags.switch-value", `--json takes no value, and was given "yes"`, [2]string{"flag", "--json"}},
		{[]string{"count", "f", "--stack"}, CodeUsage, "flags.value-missing", "--stack needs a value", [2]string{"flag", "--stack"}},
		{[]string{"count", "f", "--stack", "--json"}, CodeUsage, "flags.value-missing", "--stack needs a value", [2]string{"flag", "--stack"}},
		{[]string{"count", "f", "--num"}, CodeUsage, "flags.value-missing", "", [2]string{"flag", "--num"}},
		{[]string{"count", "f", "--kind", "z"}, CodeUsage, "flags.value-invalid", "--kind must be one of \"a\",\"b\" but got \"z\"", [2]string{"value", "z"}},
		{[]string{"count", "f", "--num", "x"}, CodeUsage, "flags.value-invalid", "", [2]string{"flag", "--num"}},
		{[]string{"nope"}, CodeUnknownCommand, "command.unknown", `there is no command "nope"`, [2]string{"command", "nope"}},
		{[]string{"group", "nope"}, CodeUnknownCommand, "command.unknown", `there is no command "group nope"`, [2]string{"command", "group nope"}},
		{[]string{"version", "extra"}, CodeUsage, "args.unexpected", `version takes no arguments, and was given "extra"`, [2]string{"argument", "extra"}},
		{[]string{"count"}, CodeUsage, "args.missing", "count needs <file>", [2]string{}},
		{[]string{"group"}, CodeUsage, "command.missing", "group needs a subcommand", [2]string{}},
	} {
		p := problem(t, c.args...)
		if p.Code != c.code || p.Rule != c.rule {
			t.Errorf("%v: code %d and rule %q, want %d and %q (%s)", c.args, p.Code, p.Rule, c.code, c.rule, p.Message)
		}
		if c.message != "" && p.Message != c.message {
			t.Errorf("%v: message %q, want %q", c.args, p.Message, c.message)
		}
		if p.Fix == "" {
			t.Errorf("%v: no fix", c.args)
		}
		if k := c.subject[0]; k != "" && p.Subject(k) != c.subject[1] {
			t.Errorf("%v: subject %s = %v, want %q", c.args, k, p.Subject(k), c.subject[1])
		}
	}
}

// The command meant is named in the fix by the kong fork's Suggest, a
// command as a person writes it, never a hidden one.
func TestAnUnknownCommandNamesTheCommandMeant(t *testing.T) {
	p := problem(t, "coun")
	if !strings.HasPrefix(p.Fix, "Did you mean 'count'? Run 'go-template-itos --help' for the commands.") {
		t.Errorf("fix = %q", p.Fix)
	}
	p = problem(t, "group", "su")
	if want := "Did you mean 'group sub'? The subcommands of group are sub and two. Run 'go-template-itos group --help' for them."; p.Fix != want {
		t.Errorf("fix = %q, want %q", p.Fix, want)
	}
	if p = problem(t, "hidden2"); strings.Contains(p.Fix, "hidden") {
		t.Errorf("fix names a hidden command: %q", p.Fix)
	}
}

// A flag meant is named in the fix, the way the unknown command is.
func TestAnUnknownFlagNamesTheFlagMeant(t *testing.T) {
	p := problem(t, "count", "f", "--jso")
	if want := "Did you mean --json? Run 'go-template-itos count --help' for its flags."; p.Fix != want {
		t.Errorf("fix = %q, want %q", p.Fix, want)
	}
	if p = problem(t, "count", "f", "--zzzzzzzz"); p.Fix != "Run 'go-template-itos count --help' for its flags." {
		t.Errorf("fix = %q", p.Fix)
	}
}

// An error that is no parse failure is a bug, internal, and a failure of the
// command line no rule names is a usage problem in kong's words.
func TestAnErrorThatIsNoParseFailureIsInternal(t *testing.T) {
	if p := ParseProblem(errors.New("kong refuses its model")); p.Code != CodeInternal || p.Rule != "internal" {
		t.Errorf("problem = %+v", p)
	}
}

func TestAnUnnamedParseFailureIsAUsageProblemInKongsWords(t *testing.T) {
	var l struct {
		Count struct {
			A string `required:""`
		} `cmd:""`
	}
	parser, err := kong.New(&l, append(Flags(), kong.Name("go-template-itos"))...)
	if err != nil {
		t.Fatal(err)
	}
	_, err = parser.Parse([]string{"count"})
	p := ParseProblem(err)
	if p.Code != CodeUsage || p.Rule != "usage" || p.Message != "count: missing flags: --a=STRING" {
		t.Errorf("problem = %+v", p)
	}
}

func TestAndListJoinsWordsAsPeopleWriteThem(t *testing.T) {
	for in, want := range map[string]string{"": "", "a": "a", "a b": "a and b", "a b c": "a, b, and c"} {
		if got := andList(strings.Fields(in)); got != want {
			t.Errorf("andList(%q) = %q, want %q", in, got, want)
		}
	}
}
