// Command nightly-report is what the nightly's report job asks of go: the body of the Nightly red
// issue, and a plain-text view of go test -json for a job's log.
//
//	nightly-report issue -branch <name> -url <run's URL> -jobs <file> -results <folder> [-findings <file>]
//	nightly-report text < go-test.json
//
// issue judges what the run's jobs concluded and what their go test -json files recorded, never a
// job's log. The job list is what GitHub's API returns for one attempt of the run (gh api
// --paginate repos/<repo>/actions/runs/<id>/attempts/<n>/jobs): each job that concluded neither
// success nor skipped is red, named with the steps that failed in it. The results folder holds a
// go test -json file for each platform and for the explore job, named <label>.json (read from
// every level under it): the tests that failed in each, a failing test of subtests named by its
// failing subtests alone, each with its first and last lines, which holds a property's seed and a
// failing check's message, and a package that failed with no test failing (it did not build, or
// timed out) by its own output. The findings file, when given, holds the lines
// tools/bin/nightly-fuzz appends, "fuzz input: <path>" for each failing fuzz input, which the
// body names with how to commit it. What cannot be read (the job list, a results file) is said in
// the body and goes to stderr, never a failure of the program: the issue still names the run.
//
// text copies what go test would have printed from its -json events on stdin, so a job's log reads
// as go test's own while the events are kept for issue: the nightly's steps tee them to a file.
//
// The workflow opens, comments on and closes the issue with gh; this program only writes. It
// imports nothing but the standard library, as the other programs under tools/bin do.
//
// Harvested from donvargax/itos, tools/bin/nightly-report, which reads each job's log; this reads
// the job list and go test's own events instead. Exit status: 0 the body printed, 2 a usage error.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	keepLines  = 10    // the first and the last lines of a failing test's output kept
	keepBytes  = 2000  // and the most bytes of them
	bodyBytes  = 60000 // GitHub's body limit is 65536
	maxFailing = 25    // the failing tests named for each file
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return usage(stderr)
	}
	switch args[0] {
	case "text":
		if len(args) != 1 {
			return usage(stderr)
		}
		text(stdin, stdout)
		return 0
	case "issue":
		return issue(args[1:], stdout, stderr)
	}
	return usage(stderr)
}

func usage(stderr io.Writer) int {
	_, _ = fmt.Fprintln(stderr, "usage: nightly-report issue -branch <name> -url <run's URL> -jobs <file> -results <folder> [-findings <file>]")
	_, _ = fmt.Fprintln(stderr, "       nightly-report text < go-test.json")
	return 2
}

func issue(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("nightly-report issue", flag.ContinueOnError)
	flags.SetOutput(stderr)
	branch := flags.String("branch", "", "the branch the nightly ran on")
	url := flags.String("url", "", "the run's URL, which the body names")
	jobsFile := flags.String("jobs", "", "the attempt's jobs, as the API lists them")
	results := flags.String("results", "", "the folder of go test -json files")
	findings := flags.String("findings", "", "the file of fuzz input lines")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() > 0 || *branch == "" || *url == "" || *jobsFile == "" || *results == "" {
		return usage(stderr)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Nightly run %s is red on %s.\n", inline(*url), inline(*branch))

	jobs, err := readJobs(*jobsFile)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "nightly-report: the jobs could not be read: %v\n", err)
		fmt.Fprintf(&b, "\nIts jobs could not be read: see the run's log.\n\n%s\n", fence(firstLine(err)))
	}
	red := redJobs(jobs)
	if err == nil && len(red) == 0 {
		b.WriteString("\nNo job concluded red: see the run's log.\n")
	}
	for _, j := range red {
		fmt.Fprintf(&b, "\n### %s: %s\n", inline(j.Name), inline(j.Conclusion))
		if failed := failedSteps(j); len(failed) > 0 {
			fmt.Fprintf(&b, "\nFailed steps: %s.\n", strings.Join(failed, ", "))
		} else {
			b.WriteString("\nNo step of it failed: it was cancelled, timed out or its runner stopped.\n")
		}
	}

	files, err := readResults(*results)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "nightly-report: the go test results could not be read: %v\n", err)
		fmt.Fprintf(&b, "\nThe go test results could not be read.\n\n%s\n", fence(firstLine(err)))
	}
	named := 0
	for _, f := range files {
		if len(f.failures) == 0 {
			continue
		}
		named++
		fmt.Fprintf(&b, "\n### Failing tests on %s\n", inline(f.label))
		for i, fail := range f.failures {
			if i == maxFailing {
				fmt.Fprintf(&b, "\n... and %d more.\n", len(f.failures)-maxFailing)
				break
			}
			fmt.Fprintf(&b, "\n- %s\n\n%s\n", fail.name, fence(fail.output))
		}
	}
	if named == 0 && len(red) > 0 {
		b.WriteString("\nNo go test result records a failing test: the failure is in a step's own output.\n")
	}

	if inputs := readFindings(*findings, stderr); len(inputs) > 0 {
		b.WriteString("\n### Failing fuzz inputs\n\n")
		for _, in := range inputs {
			fmt.Fprintf(&b, "- %s\n", in)
		}
		b.WriteString("\nThe job's log prints each input's bytes. To keep one as a regression, commit that file at that path: go test replays it in every run.\n")
	}

	body := b.String()
	if len(body) > bodyBytes {
		body = strings.ToValidUTF8(body[:bodyBytes], "") + "\n\n... cut: see the run's log.\n"
	}
	_, _ = fmt.Fprint(stdout, body)
	return 0
}

// job is a job of the attempt, as the jobs call lists it.
type job struct {
	Name       string `json:"name"`
	Conclusion string `json:"conclusion"` // empty while it runs
	Steps      []struct {
		Name       string `json:"name"`
		Conclusion string `json:"conclusion"`
	} `json:"steps"`
}

// readJobs decodes the jobs file: gh --paginate prints one JSON object a page, back to back.
func readJobs(path string) ([]job, error) {
	data, err := os.ReadFile(path) //nolint:gosec // the job list the workflow wrote, named by its own flag
	if err != nil {
		return nil, err
	}
	var jobs []job
	dec := json.NewDecoder(strings.NewReader(string(data)))
	for {
		var page struct {
			Jobs []job `json:"jobs"`
		}
		if err := dec.Decode(&page); errors.Is(err, io.EOF) {
			return jobs, nil
		} else if err != nil {
			return jobs, err
		}
		jobs = append(jobs, page.Jobs...)
	}
}

// redJobs are the jobs that concluded neither success nor skipped; one still running has no
// conclusion yet, as the report job itself, and is left out.
func redJobs(jobs []job) []job {
	var red []job
	for _, j := range jobs {
		if j.Conclusion != "" && j.Conclusion != "success" && j.Conclusion != "skipped" {
			red = append(red, j)
		}
	}
	return red
}

// failedSteps names the steps of a job that failed or were cancelled, quoted.
func failedSteps(j job) []string {
	var names []string
	for _, s := range j.Steps {
		if s.Conclusion == "failure" || s.Conclusion == "cancelled" || s.Conclusion == "timed_out" {
			names = append(names, "`"+strings.TrimSpace(inline(s.Name))+"`")
		}
	}
	return names
}

// event is one line of go test -json.
type event struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Output  string `json:"Output"`
}

// failure is a failing test, or a package that failed with none, by its name and output's tail.
type failure struct {
	name   string
	output string
}

// resultFile is one go test -json file: its label is its name less the extension.
type resultFile struct {
	label    string
	failures []failure
}

// readResults reads every .json file under dir, in name order. A missing folder is no results,
// since a night that died before any test ran kept none; a file that cannot be read is an error
// and the others are still read.
func readResults(dir string) ([]resultFile, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			paths = append(paths, path)
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	sort.Strings(paths)
	var files []resultFile
	for _, path := range paths {
		f, openErr := os.Open(path) //nolint:gosec // a results file the workflow downloaded, found under the folder its flag names
		if openErr != nil {
			err = errors.Join(err, openErr)
			continue
		}
		label := strings.TrimSuffix(filepath.Base(path), ".json")
		files = append(files, resultFile{label: label, failures: failures(f)})
		_ = f.Close()
	}
	return files, err
}

// failures reads go test -json events and returns what failed: each failing test with no failing
// subtest, in the order they failed, and each package that failed with no test failing in it.
// A line that is not an event (the build's own output, a truncated last line) is skipped.
func failures(r io.Reader) []failure {
	type key struct{ pkg, test string }
	output := map[key][]string{}
	var failed []key
	failedPkg := map[string]bool{}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var e event
		if json.Unmarshal(scanner.Bytes(), &e) != nil || e.Package == "" {
			continue
		}
		k := key{e.Package, e.Test}
		switch e.Action {
		case "output":
			output[k] = append(output[k], e.Output)
		case "fail":
			if e.Test == "" {
				failedPkg[e.Package] = true
			} else {
				failed = append(failed, k)
			}
		}
	}
	var out []failure
	hasTest := map[string]bool{}
	for _, k := range failed {
		hasTest[k.pkg] = true
		leaf := true
		for _, other := range failed {
			if other.pkg == k.pkg && strings.HasPrefix(other.test, k.test+"/") {
				leaf = false
				break
			}
		}
		if leaf {
			out = append(out, failure{"`" + inline(k.pkg) + "` `" + inline(k.test) + "`", tail(output[k])})
		}
	}
	var pkgs []string
	for pkg := range failedPkg {
		if !hasTest[pkg] {
			pkgs = append(pkgs, pkg)
		}
	}
	sort.Strings(pkgs)
	for _, pkg := range pkgs {
		out = append(out, failure{"`" + inline(pkg) + "` failed with no failing test", tail(output[key{pkg, ""}])})
	}
	return out
}

// tail is a test's output cut to its first and last lines: what it logged before failing is at
// the start, and what it printed last (a panic, a seed) at the end, and the middle is the long
// part. The most bytes kept are the output's first and last halves of keepBytes.
func tail(chunks []string) string {
	lines := strings.Split(strings.TrimRight(strings.Join(chunks, ""), "\n"), "\n")
	if len(lines) > 2*keepLines {
		left := len(lines) - 2*keepLines
		lines = append(append(lines[:keepLines:keepLines], fmt.Sprintf("... %d lines left out ...", left)), lines[len(lines)-keepLines:]...)
	}
	text := strings.Join(lines, "\n")
	if len(text) > keepBytes {
		text = strings.ToValidUTF8(text[:keepBytes/2], "") + "\n... cut ...\n" + strings.ToValidUTF8(text[len(text)-keepBytes/2:], "")
	}
	return text
}

// readFindings reads the fuzz input lines of the findings file, the paths alone; with no file
// named, or none written (no fuzz target failed), there are none.
func readFindings(path string, stderr io.Writer) []string {
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path) //nolint:gosec // the findings file the workflow wrote, named by its own flag
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			_, _ = fmt.Fprintf(stderr, "nightly-report: the findings could not be read: %v\n", err)
		}
		return nil
	}
	var inputs []string
	for _, line := range strings.Split(string(data), "\n") {
		if p, ok := strings.CutPrefix(strings.TrimSpace(line), "fuzz input: "); ok && p != "" {
			inputs = append(inputs, "`"+inline(p)+"`")
		}
	}
	return inputs
}

// text prints each event's output as go test prints it. A line that is not an event is copied
// as it is, so the build's own failures are not lost.
func text(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var e event
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			_, _ = fmt.Fprintln(w, scanner.Text())
			continue
		}
		_, _ = fmt.Fprint(w, e.Output)
	}
}

var (
	// An escape sequence: a CSI, an OSC or another string to its BEL or ST or the line's end,
	// any other escape, or a lone ESC.
	escape  = regexp.MustCompile(`\x1b(?:\[[0-?]*[ -/]*[@-~]|[\]PX^_][^\x07\x1b]*(?:\x07|\x1b\\)?|[ -/]*[0-~])?`)
	control = regexp.MustCompile(`[\x00-\x08\x0b-\x1f\x7f]`)
)

// clean drops escape sequences, carriage returns and other control characters but tab and LF, so
// nothing a failing test printed moves the cursor or hides text in the issue.
func clean(s string) string {
	return control.ReplaceAllString(escape.ReplaceAllString(s, ""), "")
}

// inline is text safe in one line of the body: cleaned, on one line, its backticks made quotes.
func inline(s string) string {
	return strings.NewReplacer("\n", " ", "\t", " ", "`", "'").Replace(clean(s))
}

// fence is text in a code block no content of its own can close: its backtick runs are made quotes.
func fence(s string) string {
	return "```\n" + strings.ReplaceAll(clean(s), "```", "'''") + "\n```"
}

// firstLine is an error's text to its first line break, trimmed.
func firstLine(err error) string {
	line, _, _ := strings.Cut(strings.TrimSpace(err.Error()), "\n")
	return strings.TrimSpace(line)
}
