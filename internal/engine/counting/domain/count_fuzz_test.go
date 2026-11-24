package counting_test

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/itos-corp/go-template-itos/internal/engine/counting/domain"
)

// FuzzCount holds Count to what no input may break, whatever bytes a file holds: the figures are
// bounded by the input's size, a final LF ends a line without adding one, and a word is told
// apart from the next by whitespace alone. Go runs the seeds with every go test; the nightly
// (tools/bin/nightly-fuzz) is where the fuzzer explores. A failing input it writes under
// testdata/fuzz is kept by committing that file, which every later run replays.
func FuzzCount(f *testing.F) {
	for _, seed := range []string{
		"", "x", "one\n", "\n\n", "one\r\ntwo\nthree", "a b,\tc!",
		"a\xffb", "\xff\xfe\x00", " \n \n", "tab\tseparated\twords",
	} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		got := counting.Count(data)

		// A line holds at least a byte, and a word a rune.
		if got.Lines < 0 || got.Lines > len(data) {
			t.Fatalf("Count(%q).Lines = %d, outside 0..%d", data, got.Lines, len(data))
		}
		if got.Words < 0 || got.Words > utf8.RuneCount(data) {
			t.Fatalf("Count(%q).Words = %d, more than the input's runes", data, got.Words)
		}
		// Lines are the LFs, and one more for a last line left unterminated.
		want := bytes.Count(data, []byte{'\n'})
		if len(data) > 0 && data[len(data)-1] != '\n' {
			want++
		}
		if got.Lines != want {
			t.Fatalf("Count(%q).Lines = %d, want %d", data, got.Lines, want)
		}
		// A final LF ends the last line and splits no word.
		if len(data) > 0 && data[len(data)-1] != '\n' {
			if ended := counting.Count(append(bytes.Clone(data), '\n')); ended != got {
				t.Fatalf("Count(%q) = %+v but with a final LF %+v", data, got, ended)
			}
		}
		// Whitespace between two inputs keeps their words apart, so the words add.
		other := []byte("x y\n")
		joined := counting.Count(append(append(bytes.Clone(data), ' '), other...))
		if joined.Words != got.Words+counting.Count(other).Words {
			t.Fatalf("Count(%q + space + %q).Words = %d, want %d", data, other, joined.Words, got.Words+counting.Count(other).Words)
		}
	})
}
