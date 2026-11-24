package features

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTimingsReportTheScenariosAndTheKindsOfStep(t *testing.T) {
	file := filepath.Join(t.TempDir(), "timings.txt")
	// A run given -timings (the nightly's) keeps the scenarios' timings in the same state, which End
	// writes once the suite is done: this test starts from none and leaves them as it found them.
	savedFile, savedStart, savedScenarios, savedKinds := *timingsFile, timings.start, timings.scenarios, timings.kinds
	timings.scenarios, timings.kinds = nil, map[string]*kindTime{}
	defer func() {
		*timingsFile, timings.start, timings.scenarios, timings.kinds = savedFile, savedStart, savedScenarios, savedKinds
	}()
	*timingsFile = file
	startTimings()
	timeKind(`it exits with code "…"`, 2*time.Second)
	timeKind(`it exits with code "…"`, time.Second)
	if err := timed("(the scenario's set-up)", func() error { return nil }); err != nil {
		t.Fatal(err)
	}
	timings.scenarios = append(timings.scenarios, scenarioTime{"slow one", "a.feature", 3 * time.Second})
	if err := writeTimings(); err != nil {
		t.Fatal(err)
	}
	report, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	text := string(report)
	for _, want := range []string{"1 scenarios in", "The slowest scenarios:", "3.00s  a.feature: slow one", `   2  1.500s`, `it exits with code "…"`, "(the scenario's set-up)"} {
		if !strings.Contains(text, want) {
			t.Errorf("the report lacks %q:\n%s", want, text)
		}
	}
}

func TestTimingsKeepNothingWithoutTheFlag(t *testing.T) {
	defer func(f string) { *timingsFile = f }(*timingsFile)
	*timingsFile = ""
	timeKind("a kind", time.Second)
	if _, kept := timings.kinds["a kind"]; kept {
		t.Error("a kind was timed with no -timings")
	}
	if err := writeTimings(); err != nil {
		t.Errorf("writeTimings with no file = %v", err)
	}
}
