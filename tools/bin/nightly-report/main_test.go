package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// jobsPages is two pages of the API's job list, back to back as gh --paginate prints them: a
// green job, a skipped one, one still running, a red one that failed in two steps, a cancelled one
// with no step of its own failed, and a red one with a hostile name.
const jobsPages = `{"total_count":4,"jobs":[
 {"name":"platform (ubuntu-latest)","conclusion":"success","steps":[{"name":"Tests","conclusion":"success"}]},
 {"name":"skipped one","conclusion":"skipped","steps":[]},
 {"name":"report","conclusion":null,"steps":[{"name":"Run","conclusion":null}]},
 {"name":"explore","conclusion":"failure","steps":[
   {"name":"Set up job","conclusion":"success"},
   {"name":"Race detector","conclusion":"failure"},
   {"name":"Coverage floor","conclusion":"failure"},
   {"name":"Post","conclusion":"success"}]}]}
{"total_count":4,"jobs":[
 {"name":"platform (windows-latest)","conclusion":"cancelled","steps":[{"name":"Tests","conclusion":"skipped"}]},
 {"name":"x` + "`" + `y\u001b[31m","conclusion":"timed_out","steps":[{"name":"s` + "`" + `1\n","conclusion":"timed_out"}]}]}
`

func write(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func line(action, pkg, test, output string) string {
	data, err := json.Marshal(struct{ Action, Package, Test, Output string }{action, pkg, test, output})
	if err != nil {
		panic(err)
	}
	return string(data) + "\n"
}

func runProgram(t *testing.T, stdin string, args ...string) (status int, stdout, stderr string) {
	t.Helper()
	var out, errs bytes.Buffer
	status = run(args, strings.NewReader(stdin), &out, &errs)
	return status, out.String(), errs.String()
}

func TestUsageErrorsExitTwo(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[]}`)
	for name, args := range map[string][]string{
		"no mode":                    {},
		"unknown mode":               {"report"},
		"text with arguments":        {"text", "x"},
		"issue with no flags":        {"issue"},
		"issue with no url":          {"issue", "-branch", "b", "-jobs", jobs, "-results", dir},
		"issue with no jobs":         {"issue", "-branch", "b", "-url", "u", "-results", dir},
		"issue with no folder":       {"issue", "-branch", "b", "-url", "u", "-jobs", jobs},
		"issue with no branch":       {"issue", "-url", "u", "-jobs", jobs, "-results", dir},
		"issue with an unknown flag": {"issue", "-nope"},
		"issue with an argument":     {"issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", dir, "extra"},
	} {
		t.Run(name, func(t *testing.T) {
			status, stdout, stderr := runProgram(t, "", args...)
			if status != 2 || stdout != "" || stderr == "" {
				t.Errorf("run(%q) = %d, stdout %q, stderr %q; want 2, no stdout and a message", args, status, stdout, stderr)
			}
		})
	}
}

func TestTextPrintsEachEventsOutputAndKeepsLinesThatAreNotEvents(t *testing.T) {
	input := line("start", "p", "", "") +
		line("output", "p", "TestA", "=== RUN   TestA\n") +
		"# p [p.test]\n" +
		line("output", "p", "TestA", "--- FAIL: TestA (0.00s)\n") +
		"\n"
	status, stdout, _ := runProgram(t, input, "text")
	want := "=== RUN   TestA\n# p [p.test]\n--- FAIL: TestA (0.00s)\n\n"
	if status != 0 || stdout != want {
		t.Errorf("text = %d, %q; want 0, %q", status, stdout, want)
	}
}

func TestIssueNamesEachRedJobItsFailedStepsAndItsFailingTests(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", jobsPages)
	results := filepath.Join(dir, "results")
	// A failing parent is named by its failing subtest alone; a passing test, a skipped one and a
	// package that passed are not named; a package that failed to build is named by its output.
	write(t, results, "go-test-linux/linux.json",
		line("run", "example.com/m/a", "TestParent", "")+
			line("output", "example.com/m/a", "TestParent", "=== RUN   TestParent\n")+
			line("output", "example.com/m/a", "TestParent/sub", "    a_test.go:9: seed 7\n")+
			line("fail", "example.com/m/a", "TestParent/sub", "")+
			line("fail", "example.com/m/a", "TestParent", "")+
			line("pass", "example.com/m/a", "TestFine", "")+
			line("skip", "example.com/m/a", "TestSkipped", "")+
			line("fail", "example.com/m/a", "", "")+
			line("pass", "example.com/m/b", "", "")+
			line("output", "example.com/m/c", "", "c.go:1:1: syntax error\n")+
			line("fail", "example.com/m/c", "", "")+
			"not an event\n")
	write(t, results, "windows-1of2.json", line("pass", "example.com/m/a", "TestParent", ""))
	findings := write(t, dir, "findings.txt", "noise\nfuzz input: internal/x/testdata/fuzz/FuzzX/abc\nfuzz input: \n")

	status, stdout, stderr := runProgram(t, "", "issue", "-branch", "stack/go", "-url", "https://example.com/run/1",
		"-jobs", jobs, "-results", results, "-findings", findings)
	if status != 0 || stderr != "" {
		t.Fatalf("issue = %d, stderr %q; want 0 and none", status, stderr)
	}
	for _, want := range []string{
		"Nightly run https://example.com/run/1 is red on stack/go.\n",
		"### explore: failure\n\nFailed steps: `Race detector`, `Coverage floor`.\n",
		"### platform (windows-latest): cancelled\n\nNo step of it failed: it was cancelled, timed out or its runner stopped.\n",
		"### x'y: timed_out\n\nFailed steps: `s'1`.\n",
		"### Failing tests on linux\n\n- `example.com/m/a` `TestParent/sub`\n\n```\n    a_test.go:9: seed 7\n```\n",
		"- `example.com/m/c` failed with no failing test\n\n```\nc.go:1:1: syntax error\n```\n",
		"### Failing fuzz inputs\n\n- `internal/x/testdata/fuzz/FuzzX/abc`\n",
		"commit that file at that path",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("the body lacks %q:\n%s", want, stdout)
		}
	}
	for _, unwanted := range []string{
		"platform (ubuntu-latest)", "skipped one", "report", "`TestParent`", "TestFine", "TestSkipped", "example.com/m/b",
		"Failing tests on windows", "No go test result", "\x1b", "- `internal/x/testdata/fuzz/FuzzX/abc`\n- ",
	} {
		if strings.Contains(stdout, unwanted) {
			t.Errorf("the body names %q:\n%s", unwanted, stdout)
		}
	}
}

func TestIssueSaysWhenNoTestFailedInARedRun(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[{"name":"explore","conclusion":"failure","steps":[{"name":"Fuzz","conclusion":"failure"}]}]}`)
	// The results folder is not there: the night died before any test wrote its events.
	status, stdout, stderr := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", filepath.Join(dir, "none"))
	if status != 0 || stderr != "" {
		t.Fatalf("issue = %d, stderr %q; want 0 and none", status, stderr)
	}
	want := "Nightly run u is red on b.\n\n### explore: failure\n\nFailed steps: `Fuzz`.\n\nNo go test result records a failing test: the failure is in a step's own output.\n"
	if stdout != want {
		t.Errorf("body = %q, want %q", stdout, want)
	}
}

func TestIssueSaysWhenTheJobsCannotBeRead(t *testing.T) {
	dir := t.TempDir()
	for name, jobs := range map[string]string{
		"missing":   filepath.Join(dir, "none.json"),
		"not JSON":  write(t, dir, "bad.json", "<html>\nbad gateway"),
		"truncated": write(t, dir, "cut.json", `{"jobs":[{"name":"a"`),
	} {
		t.Run(name, func(t *testing.T) {
			status, stdout, stderr := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", dir)
			if status != 0 {
				t.Fatalf("issue = %d, want 0", status)
			}
			if !strings.HasPrefix(stdout, "Nightly run u is red on b.\n\nIts jobs could not be read: see the run's log.\n\n```\n") {
				t.Errorf("body = %q, want the run named and the jobs said unreadable", stdout)
			}
			if !strings.Contains(stderr, "nightly-report: the jobs could not be read: ") {
				t.Errorf("stderr = %q, want the error", stderr)
			}
			if strings.Contains(stdout, "bad gateway") {
				t.Errorf("body = %q, want the error's first line alone", stdout)
			}
		})
	}
}

func TestIssueSaysWhenNoJobIsRed(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[{"name":"a","conclusion":"success"},{"name":"b","conclusion":"skipped"}]}`)
	_, stdout, _ := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", dir)
	if want := "Nightly run u is red on b.\n\nNo job concluded red: see the run's log.\n"; stdout != want {
		t.Errorf("body = %q, want %q", stdout, want)
	}
}

func TestIssueKeepsAFailingTestsOutputToItsTailAndSafe(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[{"name":"j","conclusion":"failure"}]}`)
	var events strings.Builder
	for i := range 30 {
		events.WriteString(line("output", "p", "TestLong", fmt.Sprintf("line %d\n", i)))
	}
	events.WriteString(line("output", "p", "TestLong", "```\x1b[31mred\x1b[0m\r\x00\n"))
	events.WriteString(line("fail", "p", "TestLong", ""))
	events.WriteString(line("output", "p", "TestBig", strings.Repeat("é", 3000)))
	events.WriteString(line("fail", "p", "TestBig", ""))
	write(t, dir, "r/linux.json", events.String())
	_, stdout, _ := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", filepath.Join(dir, "r"))

	if !strings.Contains(stdout, "line 9\n... 11 lines left out ...\nline 21\n") || strings.Contains(stdout, "line 10\n") {
		t.Errorf("the output of TestLong is not cut to its first and last %d lines:\n%s", keepLines, stdout)
	}
	if !strings.Contains(stdout, "'''red\n") || strings.Contains(stdout, "\x1b") || strings.Contains(stdout, "\r") || strings.Contains(stdout, "\x00") {
		t.Errorf("the output of TestLong is not made safe:\n%q", stdout)
	}
	section := stdout[strings.Index(stdout, "`p` `TestBig`"):]
	if big := strings.Count(section, "é"); big != keepBytes/2 || !strings.Contains(section, "\n... cut ...\n") {
		t.Errorf("TestBig keeps %d runes, want the first and last %d bytes, %d runes, with the cut said", big, keepBytes/2, keepBytes/2)
	}
}

func TestIssueNamesAtMostMaxFailingTestsOfAFileAndCutsAHugeBody(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[{"name":"j","conclusion":"failure"}]}`)
	var events strings.Builder
	for i := range maxFailing + 3 {
		events.WriteString(line("fail", "p", fmt.Sprintf("Test%02d", i), ""))
	}
	write(t, dir, "r/linux.json", events.String())
	_, stdout, _ := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", filepath.Join(dir, "r"))
	if !strings.Contains(stdout, "`Test24`") || strings.Contains(stdout, "`Test25`") || !strings.Contains(stdout, "\n... and 3 more.\n") {
		t.Errorf("the failing tests are not cut to %d with the rest counted:\n%s", maxFailing, stdout)
	}

	// A red job's name is as long as the body may be.
	long := write(t, dir, "long.json", fmt.Sprintf(`{"jobs":[{"name":%q,"conclusion":"failure"}]}`, strings.Repeat("n", 2*bodyBytes)))
	_, stdout, _ = runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", long, "-results", dir)
	if !strings.HasSuffix(stdout, "\n\n... cut: see the run's log.\n") || len(stdout) > bodyBytes+100 {
		t.Errorf("a body of %d bytes is not cut near %d", len(stdout), bodyBytes)
	}
}

func TestIssueReportsAnUnreadableResultsFileAndStillReadsTheOthers(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[{"name":"j","conclusion":"failure"}]}`)
	results := filepath.Join(dir, "r")
	write(t, results, "good.json", line("fail", "p", "TestGood", ""))
	bad := write(t, results, "bad.json", "")
	if err := os.Chmod(bad, 0); err != nil {
		t.Fatal(err)
	}
	if f, err := os.Open(bad); err == nil {
		_ = f.Close()
		t.Skip("the file is readable whatever its mode (running as root)")
	}
	_, stdout, stderr := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", results)
	if !strings.Contains(stdout, "The go test results could not be read.") || !strings.Contains(stdout, "`p` `TestGood`") {
		t.Errorf("body = %q, want the unreadable file said and the other read", stdout)
	}
	if !strings.Contains(stderr, "the go test results could not be read") {
		t.Errorf("stderr = %q, want the error", stderr)
	}
}

func TestIssueWarnsOfUnreadableFindingsButNotOfMissingOnes(t *testing.T) {
	dir := t.TempDir()
	jobs := write(t, dir, "jobs.json", `{"jobs":[{"name":"j","conclusion":"failure"}]}`)
	_, _, stderr := runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", dir, "-findings", filepath.Join(dir, "none"))
	if stderr != "" {
		t.Errorf("a findings file not written warns: %q", stderr)
	}
	// A folder cannot be read as a file.
	_, _, stderr = runProgram(t, "", "issue", "-branch", "b", "-url", "u", "-jobs", jobs, "-results", dir, "-findings", dir)
	if !strings.Contains(stderr, "the findings could not be read") {
		t.Errorf("stderr = %q, want the findings error", stderr)
	}
}
