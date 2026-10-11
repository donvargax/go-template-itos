package architecture_test

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The fixtures are independent packages in a temporary module, not broken
// mutations of the counting example. Every package/test compiles before lint.
// Local replacement modules supply mock/Rapid symbols without new dependencies.
func TestImportBoundaries(t *testing.T) {
	root := os.Getenv("ARCHITECTURE_ROOT")
	linter := os.Getenv("ARCHITECTURE_LINTER")
	if root == "" || linter == "" {
		t.Fatal("run tools/bin/architecture-check")
	}
	dir, err := os.MkdirTemp("", "T-9-architecture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	moduleBytes, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	module := strings.Fields(string(moduleBytes))[1]
	write := func(path, contents string) {
		t.Helper()
		if strings.HasSuffix(path, ".go") {
			formatted, err := format.Source([]byte(contents))
			if err != nil {
				t.Fatalf("format %s: %v", path, err)
			}
			contents = string(formatted)
		}
		path = filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// Decision 6's packages: the domain and its port subpackages, two
	// adapters, a slice named for its command, the UI, cmd and anything else.
	for _, path := range []string{"internal/counting", "internal/counting/port", "internal/counting/port/porttest", "internal/disk", "internal/stdin", "internal/count", "internal/cli", "cmd/probe", "internal/other"} {
		write(path+"/stub.go", "package probe\n\nconst Value = 1\n")
	}
	mocks := []string{"github.com/golang/mock/gomock", "go.uber.org/mock/gomock", "github.com/stretchr/testify/mock", "github.com/vektra/mockery/v2", "github.com/maxbrunsfeld/counterfeiter/v6"}
	const kong = "github.com/alecthomas/kong"
	mod := "module " + module + "\n\ngo 1.27\n"
	for i, path := range append(mocks, "pgregory.net/rapid", kong) {
		stub := fmt.Sprintf("stubs/library%d", i)
		version := "v0.0.0"
		if strings.HasSuffix(path, "/v2") {
			version = "v2.0.0"
		}
		if strings.HasSuffix(path, "/v6") {
			version = "v6.0.0"
		}
		mod += fmt.Sprintf("\nrequire %s %s\nreplace %s => ./%s\n", path, version, path, stub)
		write(stub+"/go.mod", "module "+path+"\n\ngo 1.27\n")
		write(stub+"/stub.go", "package probe\n\nconst Value = 1\n")
	}
	write("go.mod", mod)
	type expectation struct{ name, path, imported, rule string }
	var cases []expectation
	add := func(name, layer, imported, symbol, rule string, test bool) {
		path := layer + "/" + name + "/probe.go"
		body := "var Probe = dependency." + symbol + "\n"
		imports := fmt.Sprintf("import dependency %q\n", imported)
		if test {
			path = layer + "/" + name + "/probe_test.go"
			imports = fmt.Sprintf("import (\n dependency %q\n \"testing\"\n)\n", imported)
			body = "func TestProbe(t *testing.T) {\n _ = dependency." + symbol + "\n}\n"
			if imported == "os" && rule == "" {
				body = "func TestProbe(t *testing.T) {\n path := t.TempDir() + \"/fixture\"\n if err := dependency.WriteFile(path, []byte(\"fixture\"), 0o600); err != nil { t.Fatal(err) }\n data, err := dependency.ReadFile(path)\n if err != nil || string(data) != \"fixture\" { t.Fatalf(\"fixture: %q, %v\", data, err) }\n}\n"
			}
		}
		// Rapid's permitted case must be an immediate domain test, just as
		// the branch configuration (and its existing properties) require.
		if name == "rapid_allowed" {
			path = layer + "/" + name + "_test.go"
		}
		write(path, "package probe\n\n"+imports+"\n"+body)
		cases = append(cases, expectation{name, path, imported, rule})
	}
	// The domain arrow: the domain and its ports, code and tests, import
	// nothing of ours but the domain and its port subpackages, and no kong.
	for _, base := range []string{"internal/counting", "internal/counting/port", "internal/counting/port/porttest"} {
		for _, test := range []bool{false, true} {
			kind := "code"
			if test {
				kind = "test"
			}
			for _, target := range []string{"internal/disk", "internal/stdin", "internal/count", "internal/cli", "cmd/probe", "internal/other"} {
				add(strings.ReplaceAll(target, "/", "_")+"_"+kind, base, module+"/"+target, "Value", "domain", test)
			}
			add("kong_"+kind, base, kong, "Value", "domain", test)
			add("domain_allowed_"+kind, base, module+"/internal/counting", "Value", "", test)
			add("port_allowed_"+kind, base, module+"/internal/counting/port", "Value", "", test)
		}
		// The domain-io arrow: no I/O in the domain's own code; its tests
		// may read fixtures and use porttest's fakes.
		add("os_code", base, "os", "ReadFile", "domain-io", false)
		add("exec_code", base, "os/exec", "Command", "domain-io", false)
		add("net_code", base, "net", "Dial", "domain-io", false)
		add("http_code", base, "net/http", "Get", "domain-io", false)
		add("fixture_io_allowed", base, "os", "ReadFile", "", true)
		add("fake_allowed", base, module+"/internal/counting/port/porttest", "Value", "", true)
	}
	// The adapter arrow: an adapter's code imports nothing of ours but the
	// port, not the domain, another adapter, a slice, the UI, cmd or the
	// fakes; it may do I/O, and its tests may use the fakes.
	for _, target := range []string{"internal/counting", "internal/counting/port/porttest", "internal/disk", "internal/stdin", "internal/count", "internal/cli", "cmd/probe", "internal/other"} {
		add(strings.ReplaceAll(target, "/", "_"), "internal/disk", module+"/"+target, "Value", "infra", false)
	}
	add("port_allowed", "internal/disk", module+"/internal/counting/port", "Value", "", false)
	add("fake_allowed", "internal/disk", module+"/internal/counting/port/porttest", "Value", "", true)
	add("os_allowed", "internal/disk", "os", "ReadFile", "", false)
	add("fixture_io_allowed", "internal/disk", "os", "ReadFile", "", true)
	// What assembles them is free to import the domain and the adapter.
	for _, layer := range []string{"cmd/probe", "internal/other"} {
		add("domain_allowed", layer, module+"/internal/counting", "Value", "", false)
		add("disk_allowed", layer, module+"/internal/disk", "Value", "", false)
	}
	for i, imported := range mocks {
		for _, test := range []bool{false, true} {
			add(fmt.Sprintf("mock%d_%t", i, test), "internal/other", imported, "Value", "mocks", test)
		}
	}
	// Rapid only in the domain package's own tests, not its code, a
	// subpackage's tests (its ports and fakes) or anything else.
	add("rapid_allowed", "internal/counting", "pgregory.net/rapid", "Value", "", true)
	for _, layer := range []string{"internal/counting", "internal/counting/port", "internal/counting/port/porttest", "internal/disk", "internal/other"} {
		add("rapid_code", layer, "pgregory.net/rapid", "Value", "rapid", false)
		add("rapid_test", layer, "pgregory.net/rapid", "Value", "rapid", true)
	}
	compile := exec.Command("go", "test", "./...")
	compile.Dir = dir
	compile.Env = append(os.Environ(), "GOWORK=off", "GOPROXY=off")
	if output, err := compile.CombinedOutput(); err != nil {
		t.Fatalf("fixtures must compile and their tests pass before lint: %v\n%s", err, output)
	}
	t.Logf("all %d independent boundary fixtures compiled; fixture-I/O tests passed", len(cases))
	outputPath := filepath.Join(dir, "issues.json")
	// Absolute output paths avoid the config/fixture being on different
	// Windows drives. This changes presentation only, not any lint rule.
	lint := exec.Command(linter, "run", "--config", filepath.Join(root, ".golangci.yml"), "--output.json.path", outputPath, "--path-mode=abs", "--max-issues-per-linter=0", "--max-same-issues=0", "./...")
	lint.Dir = dir
	lint.Env = compile.Env
	output, lintErr := lint.CombinedOutput()
	if lintErr != nil {
		if exit, ok := lintErr.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
			t.Fatalf("linter failed to run, not an import refusal: %v\n%s", lintErr, output)
		}
	}
	var report struct {
		Issues []struct {
			FromLinter, Text string
			Pos              struct{ Filename string }
		}
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read actual linter JSON: %v\n%s", err, output)
	}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]bool)
	canonicalDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, issue := range report.Issues {
		if !filepath.IsAbs(issue.Pos.Filename) {
			t.Fatalf("linter did not emit an absolute path: %s", issue.Pos.Filename)
		}
		// macOS's temporary directory is reached through /var -> /private/var;
		// compare the physical files, not their different lexical spellings.
		filename, err := filepath.EvalSymlinks(issue.Pos.Filename)
		if err != nil {
			t.Fatal(err)
		}
		rel, err := filepath.Rel(canonicalDir, filename)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.ToSlash(rel)
		var expected *expectation
		for i := range cases {
			if cases[i].path == path {
				expected = &cases[i]
				break
			}
		}
		if expected == nil || expected.rule == "" || issue.FromLinter != "depguard" || !strings.Contains(issue.Text, expected.imported) || !strings.Contains(issue.Text, "from list '"+expected.rule+"'") {
			t.Errorf("unexpected diagnostic (not a passed negative): %s: %s: %s", path, issue.FromLinter, issue.Text)
			continue
		}
		seen[path] = true
	}
	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			if c.rule != "" && !seen[c.path] {
				t.Errorf("missing %s refusal of %s", c.rule, c.imported)
			}
		})
	}
}
