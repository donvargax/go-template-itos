package features

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

var scenarioListOutput = flag.String("scenario-list-output", "", "write the live scenarios, one ID, file and driver set per tab-separated line, to this file")

// TestScenarioList is the entry point of tools/bin/scenario-list: it lists
// the live scenarios with no itos, and runs no step.
func TestScenarioList(t *testing.T) {
	if *scenarioListOutput == "" {
		t.Skip("invoked only by tools/bin/scenario-list")
	}
	listed, err := List(".")
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	for _, one := range listed {
		fmt.Fprintf(&out, "%s\t%s\t%s\n", one.ID, one.File, one.Drivers)
	}
	if err := os.WriteFile(*scenarioListOutput, []byte(out.String()), 0o600); err != nil {
		t.Fatal(err)
	}
}
