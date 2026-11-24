package porttest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
)

// recorder is the Reporter the suite is run against here, so a suite that
// finds nothing wrong, and one that finds what a broken Files does wrong,
// are both seen.
type recorder struct{ errors []string }

func (*recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

// reads is a Files made of a function, so a test breaks one path of an
// otherwise good Files.
type reads func(ctx context.Context, path string) ([]byte, port.ReadFailure)

func (f reads) Read(ctx context.Context, path string) ([]byte, port.ReadFailure) { return f(ctx, path) }

var content = []byte("one two\nthree\n")

// good is the fake the suite is held to pass.
func good() Files {
	return Files{
		Data:     map[string][]byte{"present": content},
		Failures: map[string]port.ReadFailure{"unreadable": Unreadable("unreadable")},
	}
}

// open is a case over files, its paths the ones good answers for.
func open(files port.Files) func() FilesCase {
	return func() FilesCase {
		return FilesCase{Files: files, Present: "present", Content: content, Absent: "absent", Unreadable: "unreadable"}
	}
}

// breaking is good with one path's answer replaced, and with a done context's
// when onlyDone is set, so each check is shown failing alone.
func breaking(path string, onlyDone bool, answer func(data []byte) ([]byte, port.ReadFailure)) port.Files {
	return reads(func(ctx context.Context, p string) ([]byte, port.ReadFailure) {
		data, failure := good().Read(ctx, p)
		if p == path && (!onlyDone || ctx.Err() != nil) {
			return answer(data)
		}
		return data, failure
	})
}

func TestFilesContractPassesAGoodFiles(t *testing.T) {
	r := &recorder{}
	FilesContract(r, open(good()))
	if len(r.errors) != 0 {
		t.Errorf("a good Files was reported wrong: %q", r.errors)
	}
}

func TestFilesContractReportsEachWayAFilesIsWrong(t *testing.T) {
	shared := append([]byte(nil), content...)
	for _, test := range []struct {
		name  string
		files port.Files
		want  string // what the first report starts with
	}{
		{"bytes of another file", breaking("present", false, func([]byte) ([]byte, port.ReadFailure) {
			return []byte("other"), nil
		}), "Read(present)"},
		{"a failure for a present file", breaking("present", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, &port.Unreadable{Path: "present"}
		}), "Read(present)"},
		{"nothing read from a present file", breaking("present", false, func([]byte) ([]byte, port.ReadFailure) {
			return []byte{}, nil
		}), "Read(present)"},
		{"bytes the caller's change reaches", reads(func(_ context.Context, p string) ([]byte, port.ReadFailure) {
			if p == "present" {
				return shared, nil
			}
			return good().Read(context.Background(), p)
		}), "Read(present) after"},
		{"a missing file read without a failure", breaking("absent", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, nil
		}), "Read(absent)"},
		{"a missing file called unreadable", breaking("absent", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, &port.Unreadable{Path: "absent"}
		}), "Read(absent)"},
		{"a missing file named by another path", breaking("absent", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, &port.Missing{Path: "elsewhere"}
		}), "Read(absent)"},
		{"bytes with a missing file's failure", breaking("absent", false, func([]byte) ([]byte, port.ReadFailure) {
			return []byte("x"), &port.Missing{Path: "absent"}
		}), "Read(absent)"},
		{"an unreadable file read without a failure", breaking("unreadable", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, nil
		}), "Read(unreadable)"},
		{"an unreadable file called missing", breaking("unreadable", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, &port.Missing{Path: "unreadable"}
		}), "Read(unreadable)"},
		{"an unreadable file named by another path", breaking("unreadable", false, func([]byte) ([]byte, port.ReadFailure) {
			return nil, &port.Unreadable{Path: "elsewhere"}
		}), "Read(unreadable)"},
		{"bytes with an unreadable file's failure", breaking("unreadable", false, func([]byte) ([]byte, port.ReadFailure) {
			return []byte("x"), &port.Unreadable{Path: "unreadable"}
		}), "Read(unreadable)"},
		{"partial bytes on a done context", breaking("present", true, func([]byte) ([]byte, port.ReadFailure) {
			return []byte("one"), nil
		}), "Read(present) on a done context"},
		{"bytes with a failure on a done context", breaking("present", true, func(data []byte) ([]byte, port.ReadFailure) {
			return data, &port.Unreadable{Path: "present"}
		}), "Read(present) on a done context"},
		{"a missing file read on a done context", breaking("absent", true, func([]byte) ([]byte, port.ReadFailure) {
			return nil, nil
		}), "Read(absent) on a done context"},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := &recorder{}
			FilesContract(r, open(test.files))
			if len(r.errors) == 0 || !strings.HasPrefix(r.errors[0], test.want) {
				t.Errorf("reports = %q, want the first to start with %q", r.errors, test.want)
			}
		})
	}
}

// A Files that fails a done context with a failure, as one that waits does,
// and reads whole otherwise, is held to nothing more.
func TestFilesContractAcceptsAFailureOnADoneContext(t *testing.T) {
	r := &recorder{}
	FilesContract(r, open(breaking("present", true, func([]byte) ([]byte, port.ReadFailure) {
		return nil, &port.Unreadable{Path: "present"}
	})))
	if len(r.errors) != 0 {
		t.Errorf("a failure on a done context was reported wrong: %q", r.errors)
	}
}

// A one-byte file is the shortest one the suite can change in place, so the
// suite must change the byte at the start, never one past it.
func TestFilesContractHoldsAOneByteFile(t *testing.T) {
	r := &recorder{}
	FilesContract(r, func() FilesCase {
		return FilesCase{
			Files: Files{
				Data:     map[string][]byte{"present": []byte("x")},
				Failures: map[string]port.ReadFailure{"unreadable": Unreadable("unreadable")},
			},
			Present: "present", Content: []byte("x"), Absent: "absent", Unreadable: "unreadable",
		}
	})
	if len(r.errors) != 0 {
		t.Errorf("a one-byte file was reported wrong: %q", r.errors)
	}
}
