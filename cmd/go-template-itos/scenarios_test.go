package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"

	"github.com/itos-corp/go-template-itos/features"
	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
	"github.com/itos-corp/go-template-itos/internal/engine/counting/port/porttest"
)

// The feature files run in process (decision 14): the steps of package
// features, through a driver that keeps every file in memory and runs
// go-template-itos through run with buffers and porttest's fake Files. This
// is the default and builds no binary. The binary driver (features.TestFeatures)
// runs the same scenarios against the built program where CI opts in; a
// scenario only it can hold carries @binary-only and is left out here, never
// copied.
//
//	go test ./cmd/go-template-itos -count=1 -run '^TestScenarios$'
//	go test ./cmd/go-template-itos -count=1 -run '^TestScenarios$' -scenarios=<regexp>
//
// -shard=i/n and -timings=<file> work here as they do there.

func TestScenarios(t *testing.T) {
	filter, err := features.TagFilter(features.Selection(), filepath.Join("..", "..", "features"), features.BinaryOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err := features.Begin(); err != nil {
		t.Fatal(err)
	}
	suite := godog.TestSuite{
		Name:                "go-template-itos in process",
		ScenarioInitializer: func(sc *godog.ScenarioContext) { features.Initialize(sc, &memory{}) },
		Options: &godog.Options{
			Format:      "pretty",
			Paths:       []string{filepath.Join("..", "..", "features")},
			Tags:        filter,
			Strict:      true,
			Concurrency: 1,
			TestingT:    t,
		},
	}
	failed := suite.Run() != 0
	if err := features.End(); err != nil {
		t.Error(err)
	}
	if failed {
		t.Fatal("a scenario failed")
	}
}

// memory is the in-process driver: one scenario's files and folders, held in
// memory by the path a scenario names, which the fake Files answers for.
// Nothing it does reaches the disk.
type memory struct {
	files   map[string][]byte
	folders map[string]bool
}

var _ features.Driver = (*memory)(nil)

func (m *memory) Setup() error {
	m.files, m.folders = map[string][]byte{}, map[string]bool{}
	return nil
}

func (*memory) Teardown() {}

func (m *memory) Write(name string, data []byte) error {
	m.files[name] = data
	return nil
}

// Mkdir makes a folder, which cannot be read as a file: the fake is told so,
// as the disk's read of a directory fails.
func (m *memory) Mkdir(name string) error {
	m.folders[name] = true
	return nil
}

// Run runs go-template-itos through run, in process, over the fakes
// inProcessAdapters binds, in the clean environment a scenario expects: no
// GO_TEMPLATE_ITOS_ variable of the caller's, and the run's own set.
func (m *memory) Run(r features.Run) (features.Result, error) {
	restore := cleanEnv(r.Env)
	defer restore()
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), r.Args, strings.NewReader(r.Stdin), &stdout, &stderr, m.inProcessAdapters())
	return features.Result{Exit: code, Stdout: stdout.String(), Stderr: stderr.String()}, nil
}

// inProcessAdapters is diskAdapters' twin: every port bound to porttest's
// fake over the memory. A new port is bound here by hand, and
// TestEveryPortIsBoundOnDiskAndInProcess fails until it is.
func (m *memory) inProcessAdapters() adapters {
	failures := map[string]port.ReadFailure{}
	for name := range m.folders {
		failures[name] = porttest.Unreadable(name)
	}
	return adapters{Files: porttest.Files{Data: m.files, Failures: failures}}
}

// cleanEnv unsets every GO_TEMPLATE_ITOS_ variable and sets the run's, and
// returns what puts the caller's back.
func cleanEnv(vars []features.Variable) func() {
	saved := map[string]*string{}
	save := func(name string) {
		if _, done := saved[name]; done {
			return
		}
		if value, ok := os.LookupEnv(name); ok {
			saved[name] = &value
		} else {
			saved[name] = nil
		}
	}
	for _, kv := range os.Environ() {
		if name, _, _ := strings.Cut(kv, "="); strings.HasPrefix(name, "GO_TEMPLATE_ITOS_") {
			save(name)
			_ = os.Unsetenv(name)
		}
	}
	for _, v := range vars {
		save(v.Name)
		_ = os.Setenv(v.Name, v.Value)
	}
	return func() {
		for name, value := range saved {
			if value == nil {
				_ = os.Unsetenv(name)
			} else {
				_ = os.Setenv(name, *value)
			}
		}
	}
}
