package porttest

import (
	"errors"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
)

func TestReadReturnsAnIndependentCopy(t *testing.T) {
	data := []byte("notes")
	files := Files{Data: map[string][]byte{"input": data}}
	got, failure := files.Read(t.Context(), "input")
	if failure != nil {
		t.Fatal(failure)
	}
	got[0] = 'N'
	if string(files.Data["input"]) != "notes" || string(data) != "notes" {
		t.Errorf("Read() exposed mutable fake state: map=%q source=%q", files.Data["input"], data)
	}
}

func TestReadReturnsAConfiguredFailure(t *testing.T) {
	want := Unreadable("input")
	files := Files{
		Data:     map[string][]byte{"input": []byte("notes")},
		Failures: map[string]port.ReadFailure{"input": want},
	}
	got, failure := files.Read(t.Context(), "input")
	if failure != want || got != nil {
		t.Errorf("Read() = %q, %v, want nil and the configured failure", got, failure)
	}
}

func TestReadReportsAMissingPathAsMissing(t *testing.T) {
	_, failure := Files{}.Read(t.Context(), "absent")
	var missing *port.Missing
	if !errors.As(failure, &missing) || missing.Path != "absent" {
		t.Errorf("Read() failure = %T %v, want Missing for absent", failure, failure)
	}
}

func TestUnreadableNamesThePath(t *testing.T) {
	var unreadable *port.Unreadable
	if failure := Unreadable("input"); !errors.As(failure, &unreadable) || unreadable.Path != "input" {
		t.Errorf("Unreadable() = %T %v, want Unreadable for input", failure, failure)
	}
}

func TestFakeFilesHoldTheFilesContract(t *testing.T) {
	FilesContract(t, func() FilesCase {
		content := []byte("one two\nthree\n")
		return FilesCase{
			Files: Files{
				Data:     map[string][]byte{"present": content},
				Failures: map[string]port.ReadFailure{"unreadable": Unreadable("unreadable")},
			},
			Present:    "present",
			Content:    content,
			Absent:     "absent",
			Unreadable: "unreadable",
		}
	})
}
