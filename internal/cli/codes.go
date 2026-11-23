package cli

// The exit codes, as docs/CLI.md's table gives them: itos-cc's, with 75 kept
// for a failure that may pass when run again (decision 13). 2, 3 and 70 are
// left alone, as shells, Go and sysexits already use them.
const (
	CodeNo             = 1   // a check said no
	CodeTemporary      = 75  // may pass when run again unchanged
	CodeUsage          = 80  // a usage or config error
	CodeUnknownCommand = 81  // a command or subcommand that is none
	CodeEnvironment    = 82  // a requirement not met: the environment lacks something
	CodeInternal       = 100 // an error nothing classified
	CodeInterrupted    = 130 // an interrupt stopped it (decision 8)
)

// Precedence picks the exit code of a run with problems of several codes:
// the first of these it has, not the highest, the one that says the most
// about why nothing else could be trusted (decision 13). The interrupt is
// first, as nothing else in the run can be trusted after it.
var Precedence = []int{CodeInterrupted, CodeUsage, CodeUnknownCommand, CodeInternal, CodeEnvironment, CodeTemporary, CodeNo}
