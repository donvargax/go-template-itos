package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/count/domain"
	"github.com/itos-corp/go-template-itos/internal/count/port"
	"github.com/itos-corp/go-template-itos/internal/count/port/porttest"
	"pgregory.net/rapid"
)

func TestCountCountsLogicalLinesAndUnicodeWhitespaceWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  domain.Result
	}{
		{name: "empty", input: "", want: domain.Result{}},
		{name: "one unterminated line", input: "one", want: domain.Result{Lines: 1, Words: 1}},
		{name: "terminated line", input: "one\n", want: domain.Result{Lines: 1, Words: 1}},
		{name: "empty lines", input: "\n\n", want: domain.Result{Lines: 2}},
		{name: "mixed line endings", input: "one\r\ntwo\nthree", want: domain.Result{Lines: 3, Words: 3}},
		{name: "unicode whitespace and punctuation", input: "a\u2003b,\tc!", want: domain.Result{Lines: 1, Words: 3}},
		{name: "invalid UTF-8 is not whitespace", input: string([]byte{'a', 0xff, 'b'}), want: domain.Result{Lines: 1, Words: 1}},
		{name: "long line", input: strings.Repeat("x", 128*1024), want: domain.Result{Lines: 1, Words: 1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := domain.Count([]byte(test.input)); got != test.want {
				t.Errorf("Count(%q) = %+v, want %+v", test.input, got, test.want)
			}
		})
	}
}

func TestFileReadsThroughFilesAndReturnsContentsCount(t *testing.T) {
	files := porttest.Files{Data: map[string][]byte{"notes": []byte("one\r\ntwo words\n")}}
	got, err := domain.File(files, "notes")
	if err != nil || got != (domain.Result{Lines: 2, Words: 3}) {
		t.Fatalf("File() = %+v, %v; want 2 lines and 3 words", got, err)
	}
}

func TestFileReturnsTypedMissingAndUnreadableFailures(t *testing.T) {
	for _, test := range []struct {
		name        string
		files       port.Files
		wantMissing bool
	}{
		{name: "missing", files: porttest.Files{}, wantMissing: true},
		{name: "unreadable", files: porttest.Files{Failures: map[string]port.ReadFailure{"notes": porttest.Unreadable("notes")}}, wantMissing: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := domain.File(test.files, "notes")
			if err == nil {
				t.Fatal("File() returned no failure")
			}
			var missing *port.Missing
			var unreadable *port.Unreadable
			if test.wantMissing && !errors.As(err, &missing) {
				t.Errorf("failure %T = %v, want typed Missing", err, err)
			}
			if !test.wantMissing && !errors.As(err, &unreadable) {
				t.Errorf("failure %T = %v, want typed Unreadable", err, err)
			}
		})
	}
}

func TestCountAddsAcrossNewlineTerminatedPieces(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		pieces := rapid.SliceOfN(rapid.StringMatching(`[a-z\p{Zs}\t]{0,30}`), 0, 12).Draw(t, "pieces")
		var joined string
		var want domain.Result
		for _, piece := range pieces {
			terminated := piece + "\n"
			part := domain.Count([]byte(terminated))
			want.Lines += part.Lines
			want.Words += part.Words
			joined += terminated
		}
		if got := domain.Count([]byte(joined)); got != want {
			t.Fatalf("Count(joined pieces) = %+v, sum of pieces = %+v", got, want)
		}
	})
}
