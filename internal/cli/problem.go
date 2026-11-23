package cli

import (
	"bytes"
	"encoding/json"
)

// The contract with scripts is the exit code and the --json object less its
// "message" and "fix" keys. Plain output is for people and may change in any
// release (docs/CLI.md).

// Problem is one reason a run did not succeed. Code is its exit code; Rule
// is a stable, dotted lower-case ID scripts switch on; Message and Fix are
// sentences for people. Its subject holds what the problem is about, such
// as a file or a flag, as keys of their own beside rule, message and fix in
// --json, so no script parses a message.
type Problem struct {
	Code    int
	Rule    string
	Message string
	Fix     string
	subject map[string]any
}

// NewProblem is a problem of its code, rule, message and fix, about nothing
// yet.
func NewProblem(code int, rule, message, fix string) *Problem {
	return &Problem{Code: code, Rule: rule, Message: message, Fix: fix}
}

// With sets the subject key to value and returns p. rule, message, fix and
// exit are reserved: the problem's own keys, which a subject never hides, so
// With of one of them sets a key --json never prints.
func (p *Problem) With(key string, value any) *Problem {
	if p.subject == nil {
		p.subject = map[string]any{}
	}
	p.subject[key] = value
	return p
}

// Subject is the value of the subject key, or nil.
func (p *Problem) Subject(key string) any { return p.subject[key] }

// Error is the problem's message.
func (p *Problem) Error() string { return p.Message }

// MarshalJSON is the problem as --json gives it: the subject's keys, then
// rule, message, fix and exit over them, one line. Its text is as written,
// <, > and & among it.
func (p *Problem) MarshalJSON() ([]byte, error) {
	out := map[string]any{}
	for k, v := range p.subject {
		out[k] = v
	}
	out["rule"], out["message"], out["fix"], out["exit"] = p.Rule, p.Message, p.Fix, p.Code
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// plain is the problem as a line for people: its message, then its fix.
func (p *Problem) plain() string {
	if p.Fix == "" {
		return p.Message
	}
	return p.Message + ". " + p.Fix
}
