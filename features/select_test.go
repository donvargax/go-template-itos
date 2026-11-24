package features

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const fixture = `@phase-1
Feature: fixture
  A description line.

  @ID-ONE-01 @both
  Scenario: runs under both drivers
    When it runs

  # a comment
  @ID-ONE-02 @binary-only # why
  Scenario: runs through the binary only
    When it runs

  @ID-ONE-03 @wip
  Scenario: is not written yet
    When it runs

  @ID-ONE-04
  Scenario Outline: an outline
    When it runs <x>
`

func fixtureDir(t *testing.T, text string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "one.feature"), []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestTagFilterWithoutASelectionLeavesOutWipAndTheGivenTags(t *testing.T) {
	for _, test := range []struct {
		leaveOut []string
		want     string
	}{
		{nil, "~@wip"},
		{[]string{BinaryOnly}, "~@wip&&~@binary-only"},
	} {
		got, err := TagFilter("", t.TempDir(), test.leaveOut...)
		if err != nil || got != test.want {
			t.Errorf("TagFilter(left out %v) = %q, %v, want %q", test.leaveOut, got, err, test.want)
		}
	}
}

func TestTagFilterTurnsAnExpressionIntoGodogsFilter(t *testing.T) {
	dir := fixtureDir(t, fixture)
	got, err := TagFilter(`^@ID-ONE-0[12]$`, dir, BinaryOnly)
	want := "@ID-ONE-01&&~@wip&&~@binary-only,@ID-ONE-02&&~@wip&&~@binary-only"
	if err != nil || got != want {
		t.Errorf("TagFilter = %q, %v, want %q", got, err, want)
	}
}

func TestTagFilterRefusesWhatMatchesNoTagOrIsNoExpression(t *testing.T) {
	dir := fixtureDir(t, fixture)
	if _, err := TagFilter(`^@ID-NONE-`, dir); err == nil || !strings.Contains(err.Error(), "no scenario has a tag") {
		t.Errorf("an expression matching no tag: %v", err)
	}
	if _, err := TagFilter(`(`, dir); err == nil {
		t.Error("a broken expression was accepted")
	}
	if _, err := TagFilter(`x`, filepath.Join(dir, "one.feature", "no")); err == nil {
		t.Error("an unreadable folder was accepted")
	}
}

func TestFeatureTagsListsEachTagOnceSortedAndStopsAtAComment(t *testing.T) {
	tags, err := featureTags(fixtureDir(t, fixture))
	want := []string{"@ID-ONE-01", "@ID-ONE-02", "@ID-ONE-03", "@ID-ONE-04", "@both", "@binary-only", "@phase-1", "@wip"}
	slices.Sort(want)
	if err != nil || !slices.Equal(tags, want) {
		t.Errorf("featureTags = %q, %v, want %q", tags, err, want)
	}
}

func TestListGivesTheLiveScenariosAndTheDriversThatRunThem(t *testing.T) {
	got, err := List(fixtureDir(t, fixture))
	want := []Listed{
		{"ID-ONE-01", "one.feature", "both"},
		{"ID-ONE-02", "one.feature", "binary"},
		{"ID-ONE-04", "one.feature", "both"},
	}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("List = %v, %v, want %v", got, err, want)
	}
}

func TestListRefusesAScenarioNothingCouldName(t *testing.T) {
	_, err := List(fixtureDir(t, "Feature: f\n  @count-cli\n  Scenario: nameless\n    When it runs\n"))
	if err == nil || !strings.Contains(err.Error(), "nameless") {
		t.Errorf("List = %v, want a refusal naming the scenario", err)
	}
}

func TestListOfTheRealFeaturesNamesEveryBinaryOnlyScenario(t *testing.T) {
	listed, err := List(".")
	if err != nil || len(listed) == 0 {
		t.Fatalf("List = %v, %v", listed, err)
	}
	seen := map[string]bool{}
	for _, one := range listed {
		if seen[one.ID] {
			t.Errorf("%s is listed twice", one.ID)
		}
		seen[one.ID] = true
	}
}

func TestBeginRefusesABadShardAndEndWritesNothingWithoutTimings(t *testing.T) {
	defer func(text string, i, n int) { *shardFlag, shardIndex, shardCount = text, i, n }(*shardFlag, shardIndex, shardCount)
	*shardFlag = "4/3"
	if err := Begin(); err == nil {
		t.Error("Begin accepted -shard=4/3")
	}
	*shardFlag = "2/3"
	if err := Begin(); err != nil || shardIndex != 2 || shardCount != 3 {
		t.Errorf("Begin() = %v with the share %d of %d", err, shardIndex, shardCount)
	}
	if err := End(); err != nil {
		t.Errorf("End() without -timings = %v", err)
	}
}
