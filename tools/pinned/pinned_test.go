package pinned_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Adapted from github.com/donvargax/itos-container's tools/pinned/pinned_test.go.
//
// These tests run the real launcher with curl against a local release server. Only the
// platform and release metadata are replaced in a disposable copy of it; the repository's
// pins, the real tool cache and the release services are never touched.

func TestPinnedDownload(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statuses  []int
		corrupt   bool
		wantExit  int
		wantCalls int32
		wantText  string
	}{
		{"transient gateway timeout", []int{504, 200}, false, 7, 2, "an argument with spaces\n"},
		{"gateway outage outlasts early retries", []int{504, 504, 504, 504, 200}, false, 7, 5, "an argument with spaces\n"},
		{"exhausted gateway timeout", []int{504}, false, 2, 6, "pinned: cannot fetch"},
		{"missing release", []int{404}, false, 2, 1, "pinned: cannot fetch"},
		{"bad checksum after retry", []int{504, 200}, true, 2, 2, "SHA-256 did not match the pinned release"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			archive := toolArchive(t)
			body := archive
			if tc.corrupt {
				body = append(bytes.Clone(archive), []byte("corrupt release")...)
			}
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				i := int(calls.Add(1)) - 1
				status := tc.statuses[min(i, len(tc.statuses)-1)]
				// Keep outage tests short while exercising curl's real retry logic.
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(status)
				if status == http.StatusOK {
					_, _ = w.Write(body)
				}
			}))
			defer server.Close()
			root, launcher := fixtureLauncher(t, server.URL, archive)
			run := func() (int, string) { return runLauncher(t, launcher, nil, "golangci-lint", "an argument with spaces") }
			exit, output := run()
			if exit != tc.wantExit || !strings.Contains(output, tc.wantText) {
				t.Fatalf("exit %d, want %d and %q:\n%s", exit, tc.wantExit, tc.wantText, output)
			}
			if got := calls.Load(); got != tc.wantCalls {
				t.Fatalf("download attempts = %d, want %d", got, tc.wantCalls)
			}
			if tc.wantExit == 7 {
				// The tool's nonzero exit and exact argv survive a cache hit with the
				// release server unavailable.
				server.Close()
				exit, output = run()
				if exit != 7 || output != tc.wantText || calls.Load() != tc.wantCalls {
					t.Fatalf("cached tool: exit %d, calls %d, output %q", exit, calls.Load(), output)
				}
			} else {
				// A failed fetch leaves no file behind: no binary, no half-fetched archive.
				err := filepath.WalkDir(filepath.Join(root, ".tools"), func(path string, d os.DirEntry, err error) error {
					if err == nil && !d.IsDir() {
						t.Errorf("failed download left %s", path)
					}
					return err
				})
				if err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestPinnedRedirect(t *testing.T) {
	archive := toolArchive(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset":
			http.Redirect(w, r, "/archive", http.StatusFound)
		case "/archive":
			_, _ = w.Write(archive)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	_, launcher := fixtureLauncher(t, server.URL+"/asset", archive)
	path := toolPath(t, launcher, nil)
	// Consume the path as a caller's shell would.
	run := exec.CommandContext(t.Context(), "sh", "-c", `exec "$1" "$2"`, "fixture", path, "an argument with spaces")
	result, err := run.CombinedOutput()
	if run.ProcessState == nil || run.ProcessState.ExitCode() != 7 || string(result) != "an argument with spaces\n" {
		t.Fatalf("returned tool path is unusable: %v, output %q", err, result)
	}
}

func TestPinnedWithoutGitCachesInTree(t *testing.T) {
	archive := toolArchive(t)
	server := serveArchive(archive, new(atomic.Int32))
	defer server.Close()
	root, launcher := fixtureLauncher(t, server.URL, archive)
	path := toolPath(t, launcher, nil)
	// Git for Windows' shell prints paths Go cannot stat, so only the tail is compared.
	if want := "/.tools/golangci-lint-fixture/golangci-lint"; !strings.HasSuffix(path, want) {
		t.Fatalf("path = %q, want it to end in %q", path, want)
	}
	if _, err := os.Stat(filepath.Join(root, ".tools", "golangci-lint-fixture")); err != nil {
		t.Fatalf("the tool is not cached in the tree: %v", err)
	}
}

func TestPinnedSharesOneCacheAcrossWorktrees(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on the PATH")
	}
	archive := toolArchive(t)
	var calls atomic.Int32
	server := serveArchive(archive, &calls)
	defer server.Close()
	root, launcher := fixtureLauncher(t, server.URL, archive)
	git(t, root, "init", "-q")
	git(t, root, "-c", "user.name=t", "-c", "user.email=t@example.com", "commit", "-q", "--allow-empty", "-m", "start")
	// A caller's GIT_ variables name another repository: the launcher must still use this clone's.
	other := t.TempDir()
	env := []string{"GIT_DIR=" + other, "GIT_COMMON_DIR=" + other, "GIT_WORK_TREE=" + other}
	first := toolPath(t, launcher, env)
	if want := "/.git/itos/tools/golangci-lint-fixture/golangci-lint"; !strings.HasSuffix(first, want) {
		t.Fatalf("path = %q, want it to end in %q", first, want)
	}
	if _, err := os.Stat(filepath.Join(root, ".tools")); !os.IsNotExist(err) {
		t.Fatalf("a clone with git also wrote .tools: %v", err)
	}
	// A second worktree of the clone finds the same binary and downloads nothing.
	worktree := filepath.Join(t.TempDir(), "wt")
	git(t, root, "worktree", "add", "-q", "--detach", worktree)
	if err := os.MkdirAll(filepath.Join(worktree, "tools", "bin"), 0o700); err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile(launcher)
	if err != nil {
		t.Fatal(err)
	}
	second := filepath.Join(worktree, "tools", "bin", "pinned")
	if err := os.WriteFile(second, script, 0o600); err != nil {
		t.Fatal(err)
	}
	// The shell may spell the same folder differently (/var and /private/var, /tmp and /c/...),
	// so what proves the share is that no second download happened.
	if got := toolPath(t, second, nil); !strings.HasSuffix(got, "/.git/itos/tools/golangci-lint-fixture/golangci-lint") {
		t.Fatalf("second worktree path = %q, not in the clone's cache", got)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("download attempts = %d, want 1", got)
	}
}

// The real launcher's metadata for every tool and platform, with only uname replaced.
func TestPinnedMetadata(t *testing.T) {
	script := launcherSource(t)
	end := strings.Index(script, "\n# The checkout's own clone")
	if end < 0 {
		t.Fatal("cannot locate launcher metadata boundary")
	}
	script = strings.Replace(script[:end], `"$(uname -s)-$(uname -m)"`, `"$PINNED_PLATFORM"`, 1)
	script += "\nprintf '%s\\n' \"$version\" \"$asset\" \"$sum\" \"$url\" \"$member\"\n"
	// Every tool the launcher knows is checked: its case label, then the marker and version.
	label := regexp.MustCompile(`(?m)^([a-z][a-z-]*)\)$`)
	marker := regexp.MustCompile(`(?m)^([a-z][a-z-]*)\)\n(?:\t#.*\n)*?\t# renovate: depName=(\S+)\n(?:\t#.*\n)*\tversion=(\S+)$`)
	source := launcherSource(t)
	labels := label.FindAllStringSubmatch(source, -1)
	tools := marker.FindAllStringSubmatch(source, -1)
	if len(labels) < 3 || len(labels) != len(tools) {
		t.Fatalf("%d tools in the launcher, %d with a renovate marker and version", len(labels), len(tools))
	}
	for _, m := range tools {
		tool := struct{ name, dep, version string }{m[1], m[2], m[3]}
		for _, platform := range []string{"Linux-x86_64", "Linux-aarch64", "Darwin-x86_64", "Darwin-arm64"} {
			t.Run(tool.name+"/"+platform, func(t *testing.T) {
				cmd := exec.CommandContext(t.Context(), "sh", "-s", "--", tool.name)
				cmd.Stdin = strings.NewReader(script)
				cmd.Env = append(os.Environ(), "PINNED_PLATFORM="+platform)
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("metadata failed: %v\n%s", err, out)
				}
				f := strings.Split(strings.TrimSpace(string(out)), "\n")
				if len(f) != 5 {
					t.Fatalf("metadata = %q", out)
				}
				if f[0] != tool.version || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(f[2]) ||
					!strings.HasPrefix(f[3], "https://github.com/"+tool.dep+"/releases/download/v"+f[0]+"/") ||
					!strings.HasSuffix(f[3], "/"+f[1]) || f[4] == "" {
					t.Fatalf("metadata incomplete or inconsistent with its marker: %q", f)
				}
			})
		}
	}
}

func TestPinnedBoundsDownloads(t *testing.T) {
	script := launcherSource(t)
	for _, flag := range []string{"--retry 5", "--connect-timeout 10", "--max-time 60"} {
		if !strings.Contains(script, flag) {
			t.Errorf("launcher's curl lacks %s", flag)
		}
	}
}

func serveArchive(archive []byte, calls *atomic.Int32) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		_, _ = w.Write(archive)
	}))
}

func runLauncher(t *testing.T, launcher string, env []string, args ...string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", append([]string{launcher}, args...)...)
	// The fixture server must remain local even on a host with proxies.
	cmd.Env = append(append(os.Environ(), "NO_PROXY=*", "no_proxy=*"), env...)
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil || cmd.ProcessState == nil {
		t.Fatalf("launcher did not finish: %v\n%s", err, output)
	}
	return cmd.ProcessState.ExitCode(), string(output)
}

func toolPath(t *testing.T, launcher string, env []string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", launcher, "--path", "golangci-lint")
	cmd.Env = append(append(os.Environ(), "NO_PROXY=*", "no_proxy=*"), env...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("--path failed: %v\n%s", err, &stderr)
	}
	return strings.TrimSpace(string(out))
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func launcherSource(t *testing.T) string {
	t.Helper()
	source, err := os.ReadFile("../bin/pinned")
	if err != nil {
		t.Fatal(err)
	}
	return string(source)
}

func fixtureLauncher(t *testing.T, url string, archive []byte) (string, string) {
	t.Helper()
	script := launcherSource(t)
	start := strings.Index(script, `case "$(uname -s)-$(uname -m)" in`)
	end := strings.Index(script, "\n# The checkout's own clone")
	if start < 0 || end <= start {
		t.Fatal("cannot locate launcher release metadata")
	}
	metadata := fmt.Sprintf("version=fixture\nasset=fixture.tar.gz\nmember=golangci-lint\nexe=\nurl=%s\nsum=%x\n", url, sha256.Sum256(archive))
	script = script[:start] + metadata + script[end:]
	root := t.TempDir()
	dir := filepath.Join(root, "tools", "bin")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	launcher := filepath.Join(dir, "pinned")
	if err := os.WriteFile(launcher, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, launcher
}

func toolArchive(t *testing.T) []byte {
	t.Helper()
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	body := []byte("#!/bin/sh\nprintf '%s\\n' \"$1\"\nexit 7\n")
	if err := tw.WriteHeader(&tar.Header{Name: "golangci-lint", Mode: 0o755, Size: int64(len(body))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return archive.Bytes()
}
