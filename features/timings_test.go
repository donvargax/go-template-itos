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
	defer func(f string) {
		*timingsFile = f
		timings.scenarios = nil
		timings.kinds = map[string]*kindTime{}
	}(*timingsFile)
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
