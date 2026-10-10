# CLI design guidelines

These are the rules for the `go-template-itos` command line. Each rule names its source. They are
adapted from itos-template's `docs/CLI.md`, which took them from
[itos](https://github.com/donvargax/itos), with the owner's permission (`CONTRIBUTORS.md`). Where the
CLI does not follow a rule yet, the rule says so and names the item that will change that.

Use these rules when you add or change a command, a flag, an exit code or an output. A decision
in `docs/decisions/` can change a rule; change this document in the same commit. Rules for
configuration files are in `docs/CONFIG.md`.

## Sources

| Key      | Source                                                                                                                                                         |
| -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CLIG     | Command Line Interface Guidelines, <https://clig.dev>. The anchor after the key names the section, for example CLIG `#help`.                                   |
| GNU-CLI  | GNU Coding Standards, "Standards for Command Line Interfaces", <https://www.gnu.org/prep/standards/html_node/Command_002dLine-Interfaces.html>                 |
| GNU-VER  | GNU Coding Standards, "--version", <https://www.gnu.org/prep/standards/html_node/_002d_002dversion.html>                                                       |
| GNU-HELP | GNU Coding Standards, "--help", <https://www.gnu.org/prep/standards/html_node/_002d_002dhelp.html>                                                             |
| GNU-ERR  | GNU Coding Standards, "Formatting Error Messages", <https://www.gnu.org/prep/standards/html_node/Errors.html>                                                  |
| POSIX    | POSIX.1-2024 Base Definitions, section 12.2, "Utility Syntax Guidelines", <https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap12.html#tag_12_02> |
| COBRA    | Cobra README, "Concepts", <https://github.com/spf13/cobra/blob/main/README.md#concepts>                                                                        |
| ITOS     | itos's rule that only machine output is the contract, <https://github.com/donvargax/itos/tree/main/docs/decisions>                                             |

## The contract with scripts

Scripts can rely on two things only (ITOS):

- The exit code.
- The `--json` output, less every key named `message` or `fix`. Those keys hold the same
  sentences as the plain output.

All plain output is for people and can change in any release. A script reads `--json` and the
exit code, never the plain output and never a `message`. A change to the contract is a breaking
change, and a breaking change makes a major release. Once there is a release, CI is meant to run
the last release's scenarios against the new binary. Not yet: there is no release, and
`cli-releases` brings one.

### Exit codes

| Code | Meaning                                                                                                     |
| ---- | ----------------------------------------------------------------------------------------------------------- |
| 0    | Success.                                                                                                    |
| 1    | A check said no. No command returns it yet.                                                                 |
| 2    | A usage error: the command line is wrong.                                                                   |
| 3    | The environment is missing something: the input is missing or cannot be read.                               |
| 70   | An internal error that no code classified. Report it.                                                       |
| 75   | A temporary failure, which may pass when run again with no change. No command returns it yet.               |

These codes follow grep and diff (0 yes, 1 no, 2 trouble) and the BSD `sysexits.h` values
`EX_SOFTWARE` for 70 and `EX_TEMPFAIL` for 75 (CLIG `#the-basics`). A script tells failures with
the same code apart by their rule IDs, such as `COUNT_INPUT_MISSING` and `COUNT_INPUT_UNREADABLE`.

## Rules

### Command names and grammar

1. Keep the program name short and lowercase. (POSIX guidelines 1 and 2; CLIG `#naming`.)
2. Write a subcommand name in lowercase, with dashes between words. (CLIG `#naming`.)
3. Name a group of commands with a noun, and an action in a group with a verb in the imperative.
   (CLIG `#subcommands`; COBRA `#concepts`.)
4. Use the singular for a group name.
5. Do not give two commands similar names or overlapping meanings. (CLIG `#subcommands`.)
6. Do not name a command with an everyday verb that a request could use for a different command.
   An agent picks a command by its name before it reads the help.
7. Do not add an implicit default subcommand. With no arguments, the program prints its help on
   stdout and exits 0. (CLIG `#future-proofing`.)
8. Do not accept an abbreviation of a subcommand. Add an alias only by naming it explicitly.
   (CLIG `#future-proofing`.)

### Help and version

9. Show help for `--help`, for `<command> --help` and for `-h` in any position. Write help to
   stdout and exit 0. (CLIG `#help`; GNU-HELP.) Not yet: there is no `help <command>`, which
   exits 2 as an unknown command.
10. In the help of each command, give the shape of its `--json` output and its exit codes.
    `count --help` does both.
11. Support `--version` and `version`, printing the same. The first line is
    `go-template-itos <version>`. A second line, `commit <full commit id>`, follows when the build
    knows the commit it was built from. (GNU-CLI; GNU-VER; CLIG `#arguments-and-flags`.)
12. For an unknown command, exit 2. If you can guess the command the person meant, name it.
    (CLIG `#help`.)
13. For a group with no subcommand, name the subcommands the group takes.
14. End the help with an example or two and the address for issue reports. (CLIG `#help`;
    GNU-HELP.) Not yet: the help has neither.

### Flags and arguments

kong declares every command's flags in one place, which gives rules 19 to 21 and 24 by
construction; check each when a flag is added. kong refuses an unknown flag by itself.
`internal/cli`'s `Flags` hold kong to the rest of rule 20 in this CLI's words. Every flag is a
switch (a bool, with its `--no-` pair) or takes one value (a string). A cumulative flag (a slice)
may repeat, but no command has one yet.

15. Give each flag a long form. Give a one-letter form only to the most common flags. (CLIG
    `#arguments-and-flags`; GNU-CLI.)
16. Use the standard name when one exists: `--json`, `-q`/`--quiet`, `-h`/`--help`, `--force`,
    `--version`. (CLIG `#arguments-and-flags`.)
17. Let a flag mean the same thing in every command. (CLIG `#subcommands`.)
18. Use a flag to change an action, never to select a different action. (COBRA `#concepts`.)
19. Accept `--flag=value` and `--flag value`. (CLIG `#arguments-and-flags`; GNU-CLI.)
20. Refuse with exit 2 an unknown flag, a flag with no value, a switch given a value, and a flag
    that is not cumulative given twice. A switch and its `--no-` pair count as one flag, so
    `--json --no-json` is refused. An environment variable is not the flag
    (decision 1). (CLIG `#robustness-guidelines`.)
21. Give an option-argument to its option only. Never read it as a global flag. (POSIX guidelines
    6 and 14.)
22. Do not make an option-argument optional. (POSIX guideline 7.) A choice that can be declined is
    a switch with its `--no-` pair (kong's negatable flags), never `--flag [yes|no]`.
23. Let `--` end the options, and let `-` mean stdin or stdout. (POSIX guidelines 10 and 13; CLIG
    `#arguments-and-flags`.)
24. Accept flags in any position. (CLIG `#arguments-and-flags`.)
25. Check each argument before you use it, and refuse a bad one with exit 2. (CLIG
    `#robustness-guidelines`.)

### Output

26. Write the main output to stdout. Write logs, progress and errors to stderr. (CLIG
    `#the-basics`.)
27. Print JSON only with `--json`. `--json` prints one object with `"schema": 1`, and a later
    release only adds keys to it. `count --json` prints `ok`, `file` (the argument as given, `-`
    for stdin), `lines` and `words`. (CLIG `#output`; ITOS.)
28. When the plain output looks like data, write a line on stderr telling the reader to use
    `--json` in scripts.
29. With `--json`, print an object for every failure too, a usage error included: `"ok": false`
    and a `problems` list, each with its `rule`, `message` and `fix`. (CLIG `#output`; ITOS.)
30. Do not use colour, and pass `NO_COLOR` on to any program the CLI runs. (CLIG `#output`,
    `#environment-variables`.)
31. Log through `log/slog` to stderr, warnings and above only. Logs are text on a terminal and
    JSON lines otherwise, so a log never mixes into the main output. (CLIG `#output`.)

### Errors and exit codes

32. Derive the exit code from the kind of error, in the UI alone. The domain's and adapters'
    errors are sealed sets, and each kind gets its code in one switch with no default. Only an
    error no switch classified, a bug, exits 70. Read the kind of a failure of a program you run
    from what it says; never pass its own exit code through. (CLIG `#the-basics`.)
33. Start each error line with `go-template-itos:`. Write it for people: say what happened and
    what to do next. Do not show a raw command line as the message. (GNU-ERR; CLIG `#errors`.)
34. Keep the help and the code in agreement on each exit code; a scenario checks each code the
    help names.

### Environment variables

35. Start each environment variable the CLI reads with `GO_TEMPLATE_ITOS_`, in uppercase with
    underscores. Give one to a flag only where it is worth setting once for a shell, through
    kong's `env` tag. `GO_TEMPLATE_ITOS_JSON` is the one there is. (CLIG
    `#environment-variables`.)
36. Read settings in this order: flag, then environment variable, then configuration file. A flag
    overrides its variable: `GO_TEMPLATE_ITOS_JSON=1` with `--no-json` prints the plain output.
    There is no configuration file yet. (CLIG `#configuration`.)
37. Do not let a run in CI depend on the network for an update check, and give an opt-out for any
    such check. (CLIG `#future-proofing`, "Don't create a time bomb".)

### Prompts

38. Ask a question only when stdin and stdout are terminals, and give a flag for each question, so
    a script never needs a terminal. (CLIG `#interactivity`.) No command asks anything yet.

### Changing the interface

39. Renaming a command or a flag, or changing an exit code or a key of the `--json` output, is a
    breaking change. Release the breaking changes that are ready together, in one major release.
    A commit that changes the contract says so in an `Upgrading` footer.
40. Do not keep code for compatibility with an old interface: an old name exits 2, naming the new
    one. Only the previous-release check judges compatibility, against the contract above.

### Shell completion

41. `completion <shell>` prints a short completion script for `bash`, `zsh`, `fish` or
    `powershell`, and exits 2 for any other shell. The script asks the same executable for
    candidates through the hidden `__complete` command, so it carries no second copy of the
    command model. `__complete` takes the command line's words that follow it. The last word is
    the one being completed, and an empty last word means a new word. It prints one candidate per line,
    then exactly one final instruction: `:files` leaves file-name completion to the shell, and
    `:none` asks it to offer nothing more. Candidates come only from kong's command and flag
    model. The printed comments say where to install each script.
