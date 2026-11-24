package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/release/releasetest"
)

// The bump the commits since the last release ask for, the cases itos's
// tools/selftest/release-cut.ts proves over scratch histories.
func TestBumpOf(t *testing.T) {
	cases := []struct {
		name     string
		messages []string
		kind     string
	}{
		{"a feat", []string{"feat: a thing"}, "minor"},
		{"a fix", []string{"fix: a bug"}, "patch"},
		{"a scoped fix", []string{"fix(cli): a bug"}, "patch"},
		{"a feat and a fix", []string{"fix: a bug", "feat: a thing", "docs: say so"}, "minor"},
		{"a BREAKING-CHANGE: footer", []string{"feat: a thing", "fix: a bug\n\nWhy.\n\nBREAKING-CHANGE: a key goes"}, "major"},
		{"a BREAKING CHANGE: footer", []string{"fix: a bug\n\nBREAKING CHANGE: a key goes\nTask: T-1"}, "major"},
		{"a feat! header", []string{"feat!: a thing"}, "major"},
		{"a refactor! header", []string{"refactor(config)!: a key goes"}, "major"},
		{"docs, ci, test and build commits", []string{"docs: say so", "ci: run it", "test: prove it", "build: pin it"}, "none"},
		{"a BREAKING-CHANGE line outside the last paragraph", []string{"fix: a bug\n\nBREAKING-CHANGE: not a footer here\n\nTask: T-1"}, "patch"},
		{"no commit", nil, "none"},
		{"a header that is not a Conventional Commit", []string{"Update the readme"}, "patch"},
		{"a header with a type and no space after its colon", []string{"fix:a bug"}, "patch"},
		{"a header that is not a Conventional Commit beside a docs commit", []string{"docs: say so", "Update the readme"}, "patch"},
		{"a feat beside a header that is not a Conventional Commit", []string{"Update the readme", "feat: a thing"}, "minor"},
		{"a header that is not a Conventional Commit and a breaking one", []string{"Update the readme", "feat!: a thing"}, "major"},
		{"Git's default merge header", []string{"Merge branch 'main' into topic"}, "none"},
		{"GitHub's default merge header", []string{"Merge pull request #7 from owner/topic\n\nA title."}, "none"},
		{"a merge header with a type", []string{"docs: merge main into topic"}, "none"},
		{"a header that only starts with Merge", []string{"Merged the thing"}, "patch"},
	}
	for _, c := range cases {
		b := bumpOf(c.messages)
		if b.kind != c.kind || b.commits != len(c.messages) {
			t.Errorf("%s: bumpOf = %s over %d commit(s), want %s over %d", c.name, b.kind, b.commits, c.kind, len(c.messages))
		}
	}
}

// The next version from the last release's tag, 0.0.0 with none.
func TestBumped(t *testing.T) {
	cases := []struct{ tag, kind, want string }{
		{"", "minor", "0.1.0"},
		{"", "patch", "0.0.1"},
		{"", "major", "1.0.0"},
		{"v1.2.3", "minor", "1.3.0"},
		{"v1.2.3", "patch", "1.2.4"},
		{"v1.2.3", "major", "2.0.0"},
		{"v1.9.9", "minor", "1.10.0"},
	}
	for _, c := range cases {
		if got := bumped(c.tag, c.kind); got != c.want {
			t.Errorf("bumped(%q, %s) = %s, want %s", c.tag, c.kind, got, c.want)
		}
	}
}

// A version whose major go.mod's module path does not match is refused: v0
// and v1 from a path with no /vN suffix, vN from v2 from one ending in /vN.
func TestMismatch(t *testing.T) {
	const module = "github.com/itos-corp/go-template-itos"
	cases := []struct {
		next, module string
		refused      bool
	}{
		{"0.1.0", module, false},
		{"1.4.0", module, false},
		{"2.0.0", module, true},
		{"2.0.0", module + "/v2", false},
		{"3.0.0", module + "/v2", true},
		{"1.0.0", module + "/v2", true},
		{"5.0.0", "", false},
	}
	for _, c := range cases {
		problem := mismatch(c.next, c.module)
		if (problem != "") != c.refused {
			t.Errorf("mismatch(%s, %q) = %q, want refused %v", c.next, c.module, problem, c.refused)
		}
	}
	if problem := mismatch("2.0.0", module); !strings.Contains(problem, module+"/v2") {
		t.Errorf("the refusal does not name the path a v2 needs: %s", problem)
	}
}

// git log's records split into each commit's message, trimmed.
func TestMessages(t *testing.T) {
	log := "aaa\x1ffeat: a thing\n\nWhy.\n\x1e\nbbb\x1ffix: a bug\n\x1e\n"
	got := messages(log)
	if len(got) != 2 || got[0] != "feat: a thing\n\nWhy." || got[1] != "fix: a bug" {
		t.Errorf("messages = %q", got)
	}
}

// The rule of what the binary is built from (itos-template's decision 23),
// over example repositories (its decision 19): each starts from
// releasetest.Example's v0.1.0, commits one change, and runs release-version
// in it. A module the binary links, or the toolchain, moved is a patch naming
// it; whatever leaves the binary as it was is nothing; a feat beside a moved
// module is a minor, as before; and with no tag the rule does not run.
func TestBinaryMoved(t *testing.T) {
	cases := []struct {
		name     string
		message  string
		edits    map[string][2]string // file: old, new ("" old writes the file whole)
		untagged bool
		bump     string
		says     string
	}{
		{name: "a module the binary links moved",
			edits: map[string][2]string{"go.mod": {"example.com/lib v1.0.0", "example.com/lib v1.1.0"}},
			bump:  "patch", says: "but the binary is built from what moved: example.com/lib v1.0.0 => ./lib to v1.1.0 => ./lib: 0.1.1"},
		{name: "a module only tests import moved",
			edits: map[string][2]string{"go.mod": {"example.com/testonly v1.0.0", "example.com/testonly v1.1.0"}},
			bump:  "none", says: "the binary links the same modules with the same toolchain: nothing to release"},
		{name: "a tool of the tool block moved",
			edits: map[string][2]string{"go.mod": {"example.com/tool v1.0.0", "example.com/tool v1.1.0"}},
			bump:  "none", says: "nothing to release"},
		{name: "a required module the binary never links moved",
			edits: map[string][2]string{"go.mod": {"example.com/unlinked v1.0.0", "example.com/unlinked v1.1.0"}},
			bump:  "none", says: "nothing to release"},
		{name: "the toolchain line moved",
			edits: map[string][2]string{"go.mod": {"toolchain go1.24.0", "toolchain go1.24.1"}},
			bump:  "patch", says: "but the binary is built from what moved: the toolchain go1.24.0 to go1.24.1: 0.1.1"},
		{name: "the toolchain line went and the go line moved",
			edits: map[string][2]string{"go.mod": {"go 1.24\n\ntoolchain go1.24.0\n", "go 1.25\n"}},
			bump:  "patch", says: "the toolchain go1.24.0 to go1.25: 0.1.1"},
		{name: "go.sum alone moved",
			edits: map[string][2]string{"go.sum": {"", "example.com/other v1.0.0 h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=\n"}},
			bump:  "none", says: "nothing to release"},
		{name: "a workflow's action moved",
			edits: map[string][2]string{".github/workflows/ci.yml": {"actions/checkout@v4", "actions/checkout@v5"}},
			bump:  "none", says: "nothing to release"},
		{name: "a feat beside a moved module", message: "feat: a thing",
			edits: map[string][2]string{"go.mod": {"example.com/lib v1.0.0", "example.com/lib v1.1.0"}},
			bump:  "minor", says: "1 feat(s), no breaking change: 0.2.0"},
		{name: "a moved module with no release yet", untagged: true,
			edits: map[string][2]string{"go.mod": {"example.com/lib v1.0.0", "example.com/lib v1.1.0"}},
			bump:  "none", says: "since no release (0.0.0), none a feat, a fix or a breaking change: nothing to release"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := releasetest.Example(t, !c.untagged)
			for file, edit := range c.edits {
				releasetest.Change(t, filepath.Join(repo, file), edit[0], edit[1])
			}
			message := c.message
			if message == "" {
				message = "build: move it"
			}
			releasetest.Commit(t, repo, message)
			t.Chdir(repo)
			var stdout, stderr bytes.Buffer
			if code := run(nil, &stdout, &stderr); code != 0 {
				t.Fatalf("exit %d: %s", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), "bump="+c.bump+"\n") {
				t.Errorf("stdout is not bump=%s:\n%s", c.bump, stdout.String())
			}
			if !strings.Contains(stderr.String(), c.says) {
				t.Errorf("stderr does not say %q:\n%s", c.says, stderr.String())
			}
		})
	}
}

// A shallow clone stops it before any rule: exit 2.
func TestShallow(t *testing.T) {
	repo := releasetest.Example(t, true)
	clone := filepath.Join(t.TempDir(), "clone")
	releasetest.Git(t, repo, "clone", "--quiet", "--depth", "1", "file://"+filepath.ToSlash(repo), clone)
	t.Chdir(clone)
	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr); code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "shallow clone") {
		t.Errorf("exit %d, stdout %q, stderr %q; want 2, nothing and the shallow clone named", code, stdout.String(), stderr.String())
	}
}

// Before 1.0 a breaking change makes a minor, not a major: a 0.y.z version
// promises no stable interface, so 1.0.0 is never a commit's side effect. From
// 1.0 a breaking change makes a major, as before (here refused, the example's
// module path having no /v2). The rule is itos-cc's (go-template-itos#2).
func TestBreakingBeforeOne(t *testing.T) {
	cases := []struct {
		name, message, tag string
		untagged           bool
		code               int
		stdout, says       string
	}{
		{name: "from v0.1.0", message: "feat!: a key goes",
			stdout: "next=0.2.0\nbump=minor\n", says: "1 breaking change(s), before 1.0 a minor: 0.2.0"},
		{name: "with no release yet", message: "fix: a bug\n\nBREAKING-CHANGE: a key goes", untagged: true,
			stdout: "next=0.1.0\nbump=minor\n", says: "1 breaking change(s), before 1.0 a minor: 0.1.0"},
		{name: "from v1.0.0", message: "refactor!: a key goes", tag: "v1.0.0",
			code: 1, says: "refusing v2.0.0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := releasetest.Example(t, !c.untagged)
			if c.tag != "" {
				releasetest.Git(t, repo, "tag", c.tag)
			}
			releasetest.Change(t, filepath.Join(repo, "cmd/go-template-itos/key.go"), "", "package main\n")
			releasetest.Commit(t, repo, c.message)
			t.Chdir(repo)
			var stdout, stderr bytes.Buffer
			if code := run(nil, &stdout, &stderr); code != c.code {
				t.Fatalf("exit %d, want %d: %s", code, c.code, stderr.String())
			}
			if !strings.Contains(stdout.String(), c.stdout) {
				t.Errorf("stdout does not hold %q:\n%s", c.stdout, stdout.String())
			}
			if !strings.Contains(stderr.String(), c.says) {
				t.Errorf("stderr does not say %q:\n%s", c.says, stderr.String())
			}
		})
	}
}

// The commits release-version read as patches for want of a Conventional Commit
// header, each named by its short hash and header; a merge with Git's or
// GitHub's default header is skipped, since its parents are read.
func TestUnconventional(t *testing.T) {
	log := "aaaaaaaa1111\x1fUpdate the readme\n\nWhy.\n\x1e\n" +
		"bbbbbbb2222222\x1ffeat: a thing\n\x1e\n" +
		"ccccccc3333333\x1fMerge pull request #7 from owner/topic\n\nA title.\n\x1e\n" +
		"ddddddd4444444\x1fMerge branch 'main' into topic\n\x1e\n" +
		"eeeeeee5555555\x1ffix:a bug\n\x1e\n"
	want := []string{"aaaaaaa Update the readme", "eeeeeee fix:a bug"}
	got := unconventional(log)
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("unconventional = %q, want %q", got, want)
	}
	if got := unconventional("abc\x1fshort\n\x1e"); len(got) != 1 || got[0] != "abc short" {
		t.Errorf("a hash shorter than seven characters: %q", got)
	}
}

// A commit whose header is not a Conventional Commit releases a patch and is
// named on stderr by its short hash and header, instead of failing the
// release (decision 16); a typed commit that releases nothing still releases
// nothing; a default merge header is skipped.
func TestNotConventional(t *testing.T) {
	cases := []struct {
		name    string
		message string
		bump    string
		says    []string
		not     string
	}{
		{name: "a free-form header", message: "Update the readme", bump: "patch",
			says: []string{"1 commit(s) that are not Conventional Commits, read as patches: 0.1.1",
				"read as a patch, not a Conventional Commit: %s Update the readme"}},
		{name: "a typed commit that releases nothing", message: "docs: say so", bump: "none",
			not: "read as a patch"},
		{name: "a default merge header", message: "Merge pull request #7 from owner/topic", bump: "none",
			not: "read as a patch"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := releasetest.Example(t, true)
			releasetest.Change(t, filepath.Join(repo, "README.md"), "", "text\n")
			releasetest.Commit(t, repo, c.message)
			t.Chdir(repo)
			short := strings.TrimSpace(gitOut(t, repo, "rev-parse", "--short=7", "HEAD"))
			var stdout, stderr bytes.Buffer
			if code := run(nil, &stdout, &stderr); code != 0 {
				t.Fatalf("exit %d: %s", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), "bump="+c.bump+"\n") {
				t.Errorf("stdout is not bump=%s:\n%s", c.bump, stdout.String())
			}
			for _, says := range c.says {
				says = strings.ReplaceAll(says, "%s", short)
				if !strings.Contains(stderr.String(), says) {
					t.Errorf("stderr does not say %q:\n%s", says, stderr.String())
				}
			}
			if c.not != "" && strings.Contains(stderr.String(), c.not) {
				t.Errorf("stderr says %q:\n%s", c.not, stderr.String())
			}
		})
	}
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return string(out)
}
