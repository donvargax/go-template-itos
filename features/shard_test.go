package features

import (
	"strconv"
	"testing"

	"github.com/cucumber/godog"
)

// Harvested from itos's features/shard_test.go (github.com/donvargax/itos).

func TestParseShard(t *testing.T) {
	for text, ok := range map[string]bool{"1/1": true, "2/3": true, "3/3": true, "0/2": false, "3/2": false, "1": false, "a/b": false, "1/": false, "": false} {
		if _, _, err := parseShard(text); (err == nil) != ok {
			t.Errorf("parseShard(%q): %v", text, err)
		}
	}
	if i, n, err := parseShard("2/5"); err != nil || i != 2 || n != 5 {
		t.Errorf("parseShard(2/5) = %d, %d, %v", i, n, err)
	}
}

// Each scenario is in exactly one share, the same one on every call, and the
// shares are not lopsided.
func TestInShardPartitionsTheScenarios(t *testing.T) {
	defer func(i, n int) { shardIndex, shardCount = i, n }(shardIndex, shardCount)
	const count = 3
	var scenarios []*godog.Scenario
	for i := range 300 {
		scenarios = append(scenarios, &godog.Scenario{
			Uri:   "a.feature",
			Name:  "scenario " + strconv.Itoa(i/2),
			Steps: []*godog.Step{{Text: "step " + strconv.Itoa(i)}},
		})
	}
	sizes := make([]int, count)
	for _, s := range scenarios {
		in := 0
		for index := 1; index <= count; index++ {
			shardIndex, shardCount = index, count
			if inShard(s) {
				in++
				sizes[index-1]++
				if !inShard(s) {
					t.Fatalf("%s is in share %d once and not again", s.Name, index)
				}
			}
		}
		if in != 1 {
			t.Fatalf("%s is in %d shares", s.Name, in)
		}
	}
	for index, size := range sizes {
		if size < 60 || size > 140 {
			t.Errorf("share %d holds %d of %d scenarios", index+1, size, len(scenarios))
		}
	}
	shardIndex, shardCount = 1, 1
	if !inShard(scenarios[0]) {
		t.Error("a run with no shard skips a scenario")
	}
}
