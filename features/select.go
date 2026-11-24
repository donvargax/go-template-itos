package features

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// BinaryOnly tags a scenario only the binary driver can hold: an interrupt
// or a signal, the build's stamp, what the real disk does that the fake is
// not held to (decision 14). The in-process driver leaves it out, and the
// scenario is never copied.
const BinaryOnly = "@binary-only"

// Wip tags a scenario written and not yet implemented, which never runs.
const Wip = "@wip"

var scenarios = flag.String("scenarios", "", "run only the live scenarios with a tag this regular expression matches")

// Selection is the -scenarios expression the run was given.
func Selection() string { return *scenarios }

// Begin readies a run for -shard and -timings, and tells the caller to stop
// with an error when -shard is not i/n. Call it once, before the suite runs.
func Begin() error {
	if *shardFlag != "" {
		var err error
		if shardIndex, shardCount, err = parseShard(*shardFlag); err != nil {
			return err
		}
	}
	startTimings()
	return nil
}

// End writes the timings report when -timings names a file. Call it once,
// after the suite ran.
func End() error { return writeTimings() }

// TagFilter is godog's tag filter for a selection: every live scenario
// without one, else the live scenarios carrying a tag the expression
// matches. A scenario with any of the tags in leaveOut is never selected
// (the in-process driver's @binary-only). An expression that matches no tag
// is an error, so a selection never runs nothing.
func TagFilter(selection, dir string, leaveOut ...string) (string, error) {
	never := ""
	for _, tag := range append([]string{Wip}, leaveOut...) {
		never += "&&~" + tag
	}
	if selection == "" {
		return strings.TrimPrefix(never, "&&"), nil
	}
	pattern, err := regexp.Compile(selection)
	if err != nil {
		return "", err
	}
	tags, err := featureTags(dir)
	if err != nil {
		return "", err
	}
	var each []string
	for _, tag := range tags {
		if pattern.MatchString(tag) {
			each = append(each, tag+never)
		}
	}
	if len(each) == 0 {
		return "", fmt.Errorf("no scenario has a tag that %s matches", selection)
	}
	return strings.Join(each, ","), nil
}

// Every tag written in the feature files under dir, once each, sorted.
func featureTags(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.feature"))
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, file := range files {
		err := eachTagLine(file, func(tags []string, _ string) {
			for _, tag := range tags {
				seen[tag] = true
			}
		})
		if err != nil {
			return nil, err
		}
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags, nil
}

// eachTagLine calls visit for every tag line of a feature file with its tags,
// and for every Scenario line with the tags gathered above it and its name,
// the tags nil for a line that is not a scenario's.
func eachTagLine(file string, visit func(tags []string, scenario string)) error {
	f, err := os.Open(file) //nolint:gosec // the feature files of this tree
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	var pending []string
	lines := bufio.NewScanner(f)
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		switch {
		case strings.HasPrefix(line, "@"):
			var tags []string
			for _, word := range strings.Fields(line) {
				if strings.HasPrefix(word, "#") {
					break
				}
				if strings.HasPrefix(word, "@") {
					tags = append(tags, word)
				}
			}
			visit(tags, "")
			pending = append(pending, tags...)
		case strings.HasPrefix(line, "Scenario:"), strings.HasPrefix(line, "Scenario Outline:"):
			_, name, _ := strings.Cut(line, ":")
			visit(pending, strings.TrimSpace(name))
			pending = nil
		case strings.HasPrefix(line, "Feature:"):
			pending = nil
		}
	}
	return lines.Err()
}

// Listed is one live scenario: its ID, the feature file it is in and the
// drivers that run it.
type Listed struct {
	ID, File string
	Drivers  string // "both", or "binary" for a @binary-only scenario
}

var idTag = regexp.MustCompile(`^@(ID-[A-Z]+-\d+)$`)

// List lists the live scenarios of the feature files under dir, in file and
// then reading order, with no itos: a scenario with a @wip tag is not live,
// and a scenario with no @ID-... tag is an error, since nothing could name
// it. tools/bin/scenario-list prints it.
func List(dir string) ([]Listed, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.feature"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	var listed []Listed
	for _, file := range files {
		var problem error
		err := eachTagLine(file, func(tags []string, name string) {
			if name == "" {
				return
			}
			var id string
			drivers := "both"
			for _, tag := range tags {
				if m := idTag.FindStringSubmatch(tag); m != nil {
					id = m[1]
				}
				if tag == Wip {
					return
				}
				if tag == BinaryOnly {
					drivers = "binary"
				}
			}
			if id == "" && problem == nil {
				problem = fmt.Errorf("%s: the scenario %q has no @ID-... tag", filepath.Base(file), name)
			}
			listed = append(listed, Listed{ID: id, File: filepath.Base(file), Drivers: drivers})
		})
		if err != nil {
			return nil, err
		}
		if problem != nil {
			return nil, problem
		}
	}
	return listed, nil
}
