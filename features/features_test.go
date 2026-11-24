// The features: every *.feature file in this folder, run by godog through go
// test against the go-template-itos binary built once a run from this tree,
// each scenario a subtest of TestFeatures. Harvested from itos's
// (github.com/donvargax/itos, features/features_test.go) and the owner's
// itos-template's, without either one's generator behaviour.
//
// This is the binary driver (decision 14). It runs only where
// GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1 is set, which CI's workflows set;
// anywhere else TestFeatures logs that it ran nothing, builds nothing and
// passes, since the scenarios run in process against the ports' fakes
// (cmd/go-template-itos's TestScenarios). A scenario only this driver can
// hold carries @binary-only. Every command below assumes the variable set.
//
//	go test ./features -count=1                        every live scenario
//	go test ./features -count=1 -scenarios=<regexp>    the live scenarios with a
//	                                                   tag the expression matches
//	go test ./features -count=1 -shard=1/3             the first of three shares
//	go test ./features -count=1 -timings=<file>        any run, writing where it
//	                                                   spent its time to <file>
//	GOCOVERDIR=<dir> go test ./features -count=1       either, go-template-itos
//	                                                   built with -cover, what
//	                                                   the scenarios ran in <dir>
//
// godog's own tag filter takes exact tags joined by commas; itos's run
// templates join the IDs they select with |, as a regular expression
// (itos.yaml's tests.scenario.run), so -scenarios takes that expression and
// is turned here into godog's filter. A scenario tagged @wip never runs.
//
// The scenarios run go-template-itos as another process, where Go records
// coverage only for a binary built with -cover and run with GOCOVERDIR naming
// an existing folder (go.dev/doc/build-cover); go test's own -cover sees none
// of it. So with GOCOVERDIR set, the binary is built with -cover
// -covermode=atomic -coverpkg=./..., every package of the module counted, and
// what the scenarios run lands in that folder, where itos-cc merges it with go
// test's. The build keys on GOCOVERDIR alone: a -cover binary run without it
// writes "warning: GOCOVERDIR not set, no coverage data emitted" to its
// stderr, which breaks every scenario that reads stderr, and only
// GOCOVERDIR's being set guarantees every command a scenario starts has it.
// w.env() passes it on unchanged, every environment a step runs a command in
// is w.env() changed elsewhere (noGit, knowingNoOne, withGitConfig). Never
// ITOS_CC_TEST_COVERDIR, which itos-cc sets alone for a listed test's first
// coverage run and w.env() drops with every ITOS_ variable. With GOCOVERDIR
// unset the build is the plain one.
package features

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
)

// binaryDriverVariable opts into the binary driver.
const binaryDriverVariable = "GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER"

// versionSymbol and commitSymbol are the variables the build stamps
// (-ldflags -X).
const (
	versionSymbol = "github.com/itos-corp/go-template-itos/internal/version.stamp"
	commitSymbol  = "github.com/itos-corp/go-template-itos/internal/version.commit"
)

func TestFeatures(t *testing.T) {
	if os.Getenv(binaryDriverVariable) != "1" {
		t.Logf("ran nothing: the binary driver runs only with %s=1; the scenarios run in process under go test ./cmd/go-template-itos", binaryDriverVariable)
		return
	}
	filter, err := TagFilter(Selection(), ".")
	if err != nil {
		t.Fatal(err)
	}
	if err := Begin(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(removeCallerPath)
	root, err := moduleRoot()
	if err != nil {
		t.Fatal(err)
	}
	bin, err := build(root, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	suite := godog.TestSuite{
		Name:                "go-template-itos",
		ScenarioInitializer: func(sc *godog.ScenarioContext) { Initialize(sc, &binary{root: root, bin: bin}) },
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			Tags:     filter,
			Strict:   true,
			TestingT: t,
		},
	}
	failed := suite.Run() != 0
	if err := End(); err != nil {
		t.Error(err)
	}
	if failed {
		t.Fatal("a scenario failed")
	}
}

// build builds cmd/go-template-itos from the module at root into dir,
// stamped with stampedVersion and stampedCommit, once for every scenario of
// the run: go-template-itos, or go-template-itos.exe on windows, where a
// program is found by its extension. No VCS stamping, so the build runs no
// git of the caller's. With GOCOVERDIR set it is a -cover build (see the
// package's comment).
func build(root, dir string) (string, error) {
	name := "go-template-itos"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(dir, name)
	args := []string{"build", "-trimpath", "-buildvcs=false",
		"-ldflags", "-X " + versionSymbol + "=" + stampedVersion + " -X " + commitSymbol + "=" + stampedCommit}
	if os.Getenv("GOCOVERDIR") != "" {
		args = append(args, "-cover", "-covermode=atomic", "-coverpkg=./...")
	}
	cmd := exec.Command("go", append(args, "-o", bin, "./cmd/go-template-itos")...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("building go-template-itos: %v\n%s", err, out)
	}
	return bin, nil
}
