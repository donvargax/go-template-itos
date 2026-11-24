package porttest

import (
	"bytes"
	"context"
	"errors"

	"github.com/itos-corp/go-template-itos/internal/engine/counting/port"
)

// Reporter is the small test API the contract suites need, so a *testing.T
// satisfies it and the suites import no testing package.
type Reporter interface {
	Helper()
	Errorf(string, ...any)
}

// FilesCase is what an implementation hands FilesContract: the Files under
// test and three paths it answers for. Present holds Content; Absent does not
// exist; Unreadable exists and cannot be read as a file.
type FilesCase struct {
	Files      port.Files
	Present    string
	Content    []byte
	Absent     string
	Unreadable string
}

// FilesContract holds every port.Files to the same behavior: the bytes read,
// each typed failure and what a done context does. The in-memory fake and the
// disk adapter both run it unchanged, which is why the fake may stand in for
// the adapter. open returns a fresh case for each call.
func FilesContract(t Reporter, open func() FilesCase) {
	t.Helper()
	checkBytesRead(t, open())
	checkReadsAreIndependent(t, open())
	checkMissing(t, open())
	checkUnreadable(t, open())
	checkDoneContext(t, open())
}

func checkBytesRead(t Reporter, c FilesCase) {
	t.Helper()
	got, failure := c.Files.Read(context.Background(), c.Present)
	if failure != nil || !bytes.Equal(got, c.Content) {
		t.Errorf("Read(present) = %q, %v, want %q and no failure", got, failure, c.Content)
	}
}

func checkReadsAreIndependent(t Reporter, c FilesCase) {
	t.Helper()
	first, failure := c.Files.Read(context.Background(), c.Present)
	if failure != nil || len(first) == 0 {
		t.Errorf("Read(present) = %q, %v, want content and no failure", first, failure)
		return
	}
	first[0] ^= 0xff
	second, failure := c.Files.Read(context.Background(), c.Present)
	if failure != nil || !bytes.Equal(second, c.Content) {
		t.Errorf("Read(present) after the caller changed an earlier result = %q, %v, want %q", second, failure, c.Content)
	}
}

func checkMissing(t Reporter, c FilesCase) {
	t.Helper()
	got, failure := c.Files.Read(context.Background(), c.Absent)
	var missing *port.Missing
	if got != nil || !errors.As(failure, &missing) || missing.Path != c.Absent {
		t.Errorf("Read(absent) = %q, %T %v, want no bytes and Missing for %q", got, failure, failure, c.Absent)
	}
}

func checkUnreadable(t Reporter, c FilesCase) {
	t.Helper()
	got, failure := c.Files.Read(context.Background(), c.Unreadable)
	var unreadable *port.Unreadable
	if got != nil || !errors.As(failure, &unreadable) || unreadable.Path != c.Unreadable {
		t.Errorf("Read(unreadable) = %q, %T %v, want no bytes and Unreadable for %q", got, failure, failure, c.Unreadable)
	}
}

// checkDoneContext holds what every implementation owes a done context: a read
// that ends is whole or failed, never bytes with a failure, and a failure is
// typed. Whether a read that cannot wait acts on the context is the
// implementation's own; none here may return a partial result.
func checkDoneContext(t Reporter, c FilesCase) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, failure := c.Files.Read(ctx, c.Present)
	switch {
	case failure == nil && !bytes.Equal(got, c.Content):
		t.Errorf("Read(present) on a done context = %q, want %q or a failure", got, c.Content)
	case failure != nil && got != nil:
		t.Errorf("Read(present) on a done context = %q with failure %v, want one or the other", got, failure)
	}
	if _, failure := c.Files.Read(ctx, c.Absent); failure == nil {
		t.Errorf("Read(absent) on a done context succeeded, want a failure")
	}
}
