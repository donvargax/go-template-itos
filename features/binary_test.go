// The binary driver: go-template-itos is a black box here. Each scenario gets
// a scratch folder in a temporary directory, runs the binary TestFeatures
// built from this tree in it, and reads its exit code and its output.
// Nothing here imports or reads go-template-itos's code, so the steps judge
// it by its command line alone. Harvested from itos's (github.com/donvargax/
// itos, features/steps_test.go) and the owner's itos-template's.
package features

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// binary is the Driver over the built program, and a Waiter: it holds a
// run's standard input open and interrupts it.
type binary struct {
	root string // the go-template-itos checkout: where go.mod is
	bin  string // the go-template-itos binary under test
	dir  string // the scratch folder

	running *running // the run started with its standard input left open, nil for none
}

var (
	_ Driver = (*binary)(nil)
	_ Waiter = (*binary)(nil)
)

// Setup makes the scenario's scratch folder.
func (b *binary) Setup() error {
	var err error
	b.dir, err = os.MkdirTemp("", "go-template-itos-features-")
	return err
}

// Teardown removes the scratch folder.
func (b *binary) Teardown() { _ = os.RemoveAll(b.dir) }

func (b *binary) Write(name string, data []byte) error {
	return os.WriteFile(filepath.Join(b.dir, name), data, 0o600)
}

func (b *binary) Mkdir(name string) error { return os.Mkdir(filepath.Join(b.dir, name), 0o755) }

// Run runs go-template-itos in the scratch folder, in the clean environment
// with the run's variables set as given.
func (b *binary) Run(r Run) (Result, error) {
	env := b.env()
	for _, v := range r.Env {
		env = setEnv(env, v.Name, v.Value)
	}
	cmd := exec.Command(b.bin, r.Args...)
	cmd.Dir = b.dir
	cmd.Env = env
	if r.Stdin != "" {
		cmd.Stdin = strings.NewReader(r.Stdin)
	}
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	result := Result{Stdout: stdout.String(), Stderr: stderr.String()}
	var exit *exec.ExitError
	switch {
	case err == nil:
	case errors.As(err, &exit):
		result.Exit = exit.ExitCode()
	default:
		return result, fmt.Errorf("running %s: %w", b.bin, err)
	}
	return result, nil
}

// The folder holding go.mod, above the working directory go test gives.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no go.mod above the working directory")
		}
		dir = parent
	}
}

// The environment every command runs in: the caller's, less what would make
// git or go-template-itos read anything but the scratch folder (a hook's
// GIT_DIR, the caller's GO_TEMPLATE_ITOS_ settings, CI's), and less itos-cc's
// ITOS_CC_ variables, the pinned mutation tool's, with no global or system git
// config and a fixed identity. The PATH is the caller's, unchanged.
func (b *binary) env() []string {
	var env []string
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(name, "GIT_") || strings.HasPrefix(name, "ITOS_CC_") ||
			strings.HasPrefix(name, "GO_TEMPLATE_ITOS_") || strings.HasPrefix(name, "GITHUB_") ||
			strings.HasPrefix(name, "GH_") || name == "CI" {
			continue
		}
		env = append(env, kv)
	}
	return append(env,
		"GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=go-template-itos features",
		"GIT_AUTHOR_EMAIL=features@localhost",
		"GIT_COMMITTER_NAME=go-template-itos features",
		"GIT_COMMITTER_EMAIL=features@localhost",
	)
}

// setEnv is env with the variable name, in any case, set to value.
func setEnv(env []string, name, value string) []string {
	var out []string
	for _, kv := range env {
		if n, _, _ := strings.Cut(kv, "="); !strings.EqualFold(n, name) {
			out = append(out, kv)
		}
	}
	return append(out, name+"="+value)
}
