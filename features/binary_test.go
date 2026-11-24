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
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
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
// git, itos or go-template-itos read anything but the scratch folder (a
// hook's GIT_DIR, the caller's GO_TEMPLATE_ITOS_ settings, CI's), with no
// global or system git config and a fixed identity, and the caller's PATH
// without its claude, its itos, its itos extensions or a git that is an itos
// (callerPath).
func (b *binary) env() []string {
	var env []string
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(name, "GIT_") || strings.HasPrefix(name, "ITOS_") ||
			strings.HasPrefix(name, "GO_TEMPLATE_ITOS_") || strings.HasPrefix(name, "GITHUB_") ||
			strings.HasPrefix(name, "GH_") || name == "CI" || strings.EqualFold(name, "PATH") {
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
		"PATH="+callerPath(),
	)
}

// The programs no scenario reaches on the caller's PATH: its claude, its itos
// and its itos extensions (itos-*). An itos installed on the machine would
// run its own policy where a scenario meant only go-template-itos (itos's
// T-087), and itos runs and lists every extension it finds, so one installed
// on the machine would change what a scenario sees (itos's T-080).
var hiddenAlways = []string{"claude", "itos", "itos-*"}

var (
	hiddenPathsMu sync.Mutex
	hiddenPaths   = map[string]string{}
	callerPathDir string
)

// callerPath is the caller's PATH with none of hiddenAlways on it, so no
// scenario reaches the Claude Code or the itos of the machine it runs on:
// each folder holding one is replaced by a folder of links to everything else
// in it, or left out where links cannot be made (windows without the right to
// make them).
func callerPath() string { return pathHiding(hiddenAlways...) }

// pathHiding is the caller's PATH with none of the programs named on it, nor
// a git that is an itos, made once a run for each set of names. A git shim of
// the caller's (itos git-shim install) would run the global itos for a
// scenario's git command; it is told from the real git by the rule itos
// itself goes by (isItos, git_test.go), so every git a scenario's commands
// start is the real one, as the steps' own are (gitBin, itos's T-104).
func pathHiding(names ...string) string {
	hiddenPathsMu.Lock()
	defer hiddenPathsMu.Unlock()
	key := strings.Join(names, " ")
	if text, ok := hiddenPaths[key]; ok {
		return text
	}
	var folders []string
	for i, folder := range filepath.SplitList(os.Getenv("PATH")) {
		hidden := func(e os.DirEntry) bool {
			return isOneOf(e, names) || isOneOf(e, []string{"git"}) && isItos(filepath.Join(folder, e.Name()))
		}
		entries, err := os.ReadDir(folder)
		if err != nil || !slices.ContainsFunc(entries, hidden) {
			folders = append(folders, folder)
			continue
		}
		if callerPathDir == "" {
			if callerPathDir, err = os.MkdirTemp("", "go-template-itos-features-path-"); err != nil {
				continue
			}
		}
		links := filepath.Join(callerPathDir, strconv.Itoa(len(hiddenPaths)), strconv.Itoa(i))
		if linkAllBut(folder, links, entries, hidden) == nil {
			folders = append(folders, links)
		}
	}
	text := strings.Join(folders, string(os.PathListSeparator))
	hiddenPaths[key] = text
	return text
}

// removeCallerPath removes the folders callerPath and pathHiding made.
func removeCallerPath() {
	if callerPathDir != "" {
		_ = os.RemoveAll(callerPathDir)
	}
}

// isOneOf is whether a folder's entry is a program of the names the PATH
// would find: the name, or on windows the name with any extension, in any
// case. A name ending in * is a prefix, so itos-* is every itos extension.
func isOneOf(e os.DirEntry, names []string) bool {
	name := e.Name()
	same := func(a, b string) bool { return a == b }
	starts := strings.HasPrefix
	if runtime.GOOS == "windows" {
		name = strings.TrimSuffix(name, filepath.Ext(name))
		same = strings.EqualFold
		starts = func(s, prefix string) bool {
			return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
		}
	}
	return slices.ContainsFunc(names, func(n string) bool {
		if prefix, ok := strings.CutSuffix(n, "*"); ok {
			return starts(name, prefix)
		}
		return same(name, n)
	})
}

// linkAllBut makes links a folder of links to every entry of folder but the
// hidden ones.
func linkAllBut(folder, links string, entries []os.DirEntry, hidden func(os.DirEntry) bool) error {
	if err := os.MkdirAll(links, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		if hidden(e) {
			continue
		}
		if err := os.Symlink(filepath.Join(folder, e.Name()), filepath.Join(links, e.Name())); err != nil {
			return err
		}
	}
	return nil
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
