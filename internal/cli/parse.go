package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/donvargax/kong"
)

// A command line kong could not parse is a problem of its own rule, in the
// person's words, told apart by the types of the kong fork's failures and of
// Flags' refusals, never by their text (decision 12).

// ParseProblem is err, a failure to parse a command line, as the problem
// docs/CLI.md names it, each kind told by its type:
//
//   - flags.unknown: kong's UnknownFlagError, naming the flag as written and
//     the one meant when it can guess it.
//   - flags.switch-value, flags.value-missing and flags.repeated: Flags' own
//     refusals, and kong's MissingValueError of a flag.
//   - flags.value-invalid: a flag's value kong's mapper or its enum refused.
//   - command.unknown, exit 81: kong's UnexpectedArgumentError where a
//     command was expected, naming the one meant when it can guess it.
//   - args.unexpected, args.missing and command.missing: an argument no
//     command takes, one a command needs, and a command a group needs.
//
// Any other failure of the command line is a usage problem in kong's words;
// an error that is no parse failure, such as a model kong refuses, is
// internal, a bug.
func ParseProblem(err error) *Problem {
	var (
		parse          *kong.ParseError
		unknown        *kong.UnknownFlagError
		repeated       *RepeatedFlagError
		switchValue    *SwitchValueError
		valueMissing   *ValueMissingError
		missing        *kong.MissingValueError
		invalid        *kong.InvalidValueError
		enum           *kong.EnumError
		unexpected     *kong.UnexpectedArgumentError
		positionals    *kong.MissingPositionalsError
		missingCommand *kong.MissingCommandError
	)
	if !errors.As(err, &parse) || parse.Context == nil {
		return Internal(err)
	}
	node := parse.Context.Selected()
	if node == nil {
		node = parse.Context.Model.Node
	}
	switch {
	case errors.As(err, &unknown):
		return unknownFlag(node, unknown.Flag)
	case errors.As(err, &repeated):
		return NewProblem(CodeUsage, "flags.repeated", fmt.Sprintf("--%s is given more than once", repeated.Flag.Name),
			"Give it once.").With("flag", "--"+repeated.Flag.Name)
	case errors.As(err, &switchValue):
		name := switchValue.Flag.Name
		return NewProblem(CodeUsage, "flags.switch-value", fmt.Sprintf("--%s takes no value, and was given %q", name, switchValue.Value),
			fmt.Sprintf("Write --%s alone.", name)).With("flag", "--"+name)
	case errors.As(err, &valueMissing):
		return needsAValue(valueMissing.Flag)
	case errors.As(err, &missing) && missing.Value != nil && missing.Value.Flag != nil:
		return needsAValue(missing.Value.Flag)
	case errors.As(err, &invalid) && invalid.Value.Flag != nil:
		return valueInvalid(node, invalid.Value, invalid.Err.Error(), "")
	case errors.As(err, &enum) && enum.Value.Flag != nil:
		return valueInvalid(node, enum.Value, enum.Error(), enum.Got)
	case errors.As(err, &unexpected) && unexpected.ExpectedCommand:
		return UnknownCommand(node, unexpected.Arg)
	case errors.As(err, &unexpected) && len(node.Positional) == 0:
		name := commandPath(node)
		return NewProblem(CodeUsage, "args.unexpected", fmt.Sprintf("%s takes no arguments, and was given %q", name, unexpected.Arg),
			fmt.Sprintf("Run '%s %s' alone.", programOf(node), name)).With("argument", unexpected.Arg)
	case errors.As(err, &positionals):
		var names []string
		for _, p := range positionals.Positionals {
			names = append(names, "<"+p.Name+">")
		}
		return NewProblem(CodeUsage, "args.missing", fmt.Sprintf("%s needs %s", commandPath(node), strings.Join(names, " ")),
			fmt.Sprintf("Run '%s %s --help' for what it takes.", programOf(node), commandPath(node)))
	case errors.As(err, &missingCommand):
		return NewProblem(CodeUsage, "command.missing", fmt.Sprintf("%s needs a subcommand", commandPath(node)),
			fmt.Sprintf("Name one of %s; run '%s --help' for them.", andList(commandNames(node)), strings.TrimSpace(programOf(node)+" "+commandPath(node))))
	}
	return NewProblem(CodeUsage, "usage", strings.TrimSpace(commandPath(node)+": "+err.Error()),
		fmt.Sprintf("Run '%s --help' for the commands and the flags each takes.", programOf(node)))
}

// unknownFlag is the problem of word, a flag node's command does not take,
// as it was written: the flag meant, when one is near enough to guess.
func unknownFlag(node *kong.Node, word string) *Problem {
	name, _, _ := strings.Cut(strings.TrimLeft(word, "-"), "=")
	var names []string
	for _, f := range flagsOf(node) {
		names = append(names, "--"+f.Name)
	}
	command := commandPath(node)
	if command == "" {
		command = programOf(node)
	}
	fix := fmt.Sprintf("Run '%s --help' for its flags.", strings.TrimSpace(programOf(node)+" "+commandPath(node)))
	fix = DidYouMean("--"+name, names, func(flag string) string { return flag }) + fix
	return NewProblem(CodeUsage, "flags.unknown", fmt.Sprintf("%s takes no flag %s", command, word), fix).With("flag", word)
}

// needsAValue is the problem of flag, which takes a value, given none.
func needsAValue(flag *kong.Flag) *Problem {
	return NewProblem(CodeUsage, "flags.value-missing", fmt.Sprintf("--%s needs a value", flag.Name),
		fmt.Sprintf("Give it one, such as --%s %s.", flag.Name, flag.FormatPlaceHolder())).With("flag", "--"+flag.Name)
}

// valueInvalid is the problem of a value a flag cannot take, in kong's
// words; given is the value when it is known.
func valueInvalid(node *kong.Node, value *kong.Value, message, given string) *Problem {
	p := NewProblem(CodeUsage, "flags.value-invalid", message,
		fmt.Sprintf("Run '%s --help' for its flags.", strings.TrimSpace(programOf(node)+" "+commandPath(node)))).With("flag", "--"+value.Name)
	if given != "" {
		p = p.With("value", given)
	}
	return p
}

// UnknownCommand is the problem of word, no command of node's: exit 81, the
// one meant named when one is near enough to guess.
func UnknownCommand(node *kong.Node, word string) *Problem {
	full := strings.TrimSpace(commandPath(node) + " " + word)
	fix := fmt.Sprintf("Run '%s --help' for the commands.", strings.TrimSpace(programOf(node)+" "+commandPath(node)))
	if node.Parent != nil {
		fix = fmt.Sprintf("The subcommands of %s are %s. Run '%s %s --help' for them.", commandPath(node), andList(commandNames(node)), programOf(node), commandPath(node))
	}
	fix = DidYouMean(word, commandNames(node), func(command string) string {
		return "'" + strings.TrimSpace(commandPath(node)+" "+command) + "'"
	}) + fix
	return NewProblem(CodeUnknownCommand, "command.unknown", fmt.Sprintf("there is no command %q", full), fix).With("command", full)
}

// flagsOf are the flags node's command takes, its parents' too, as a
// person sees them: none hidden.
func flagsOf(node *kong.Node) []*kong.Flag {
	var flags []*kong.Flag
	for n := node; n != nil; n = n.Parent {
		for _, f := range n.Flags {
			if !f.Hidden {
				flags = append(flags, f)
			}
		}
	}
	return flags
}

// commandPath is node's command as a person writes it after the program:
// count, or a group's subcommand after it.
func commandPath(node *kong.Node) string {
	var words []string
	for n := node; n != nil && n.Type == kong.CommandNode; n = n.Parent {
		words = append([]string{n.Name}, words...)
	}
	return strings.Join(words, " ")
}

// programOf is the name of the program node is a command of.
func programOf(node *kong.Node) string {
	for node.Parent != nil {
		node = node.Parent
	}
	return node.Name
}

// commandNames are the commands of node a person may name.
func commandNames(node *kong.Node) []string {
	var names []string
	for _, child := range node.Children {
		if child.Type == kong.CommandNode && !child.Hidden {
			names = append(names, child.Name)
		}
	}
	return names
}

// andList joins words as "a", "a and b", or "a, b, and c".
func andList(words []string) string {
	switch len(words) {
	case 0:
		return ""
	case 1:
		return words[0]
	case 2:
		return words[0] + " and " + words[1]
	}
	return strings.Join(words[:len(words)-1], ", ") + ", and " + words[len(words)-1]
}
