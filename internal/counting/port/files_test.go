package port

import "testing"

func TestFailuresNameTheirPath(t *testing.T) {
	tests := []struct {
		failure ReadFailure
		want    string
	}{
		{&Missing{Path: "in"}, "missing file: in"},
		{&Unreadable{Path: "in"}, "unreadable file: in"},
	}
	for _, test := range tests {
		if got := test.failure.Error(); got != test.want {
			t.Errorf("Error() = %q, want %q", got, test.want)
		}
	}
}
