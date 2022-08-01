// (C) 2016-2017 SAP SE or an SAP affiliate company. All rights reserved.
package trc

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
)

// external

type tb interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

// RedirectTrcOutput redirects output to the returned TestOutput
// (TestOutput.Unplug() must be called to stop this)
func RedirectTrcOutput(t tb) TestOutput {
	to := &testOut{TestOutBase: TestOutBase{tb: t}}
	to.Plug()
	return to
}

// TestOutput collects all printed trace output.
// It provides means to read the output or check expected trace entries against
type TestOutput interface {
	// Unplug stops output redirection. Will fail if unread entries are found
	Unplug()
	// Resets the output to be empty. Can be used before unplug to ignore unread lines
	Reset()
	// Match reads the next trace line and checks if it matches the given expected entry.
	Match(e Entry)
	// MatchEmpty reads the next trace line and asserts that is is empty.
	MatchEmpty()
	// ReadLine reads and returns the next trace line
	ReadLine() string
}

// Entry represents an expected trace entry. Nil pointers match any value
type Entry struct {
	Topic *string
	Level *Level
	Info  *string
	Msg   *string
}

func SPtr(s string) *string {
	return &s
}

func NewTestOutBase(t tb) TestOutBase {
	return TestOutBase{tb: t}
}

type TestOutBase struct {
	tb
	bytes.Buffer
}

// internal
type testOut struct {
	TestOutBase
	origOut io.Writer
}

func (to *testOut) Plug() {
	to.origOut = out
	out = &to.Buffer
}

func (to *testOut) Unplug() {
	to.MatchEmpty()
	out = to.origOut
}

func (to *TestOutBase) ReadLine() string {
	var s string
	var s1 string
	var err error
Loop:
	for {
		s1, err = to.ReadString('\x1e')
		s += s1
		if err == nil {
			if len(s) > 2 && s[len(s)-2] == '\x1b' {
				continue Loop
			}
			var b byte
			b, err = to.ReadByte()
			if err == nil {
				if b != '\n' {
					err = errors.New("missing \\n after \\x1e")
				}
				s += string(b)
			}
		}
		if err != nil {
			to.Error("ERROR reading next trace line:", err)
			if len(s) > 0 {
				to.Logf("String read with error: >>%s<<", s)
				to.Log("Missing \\x1e\\n at end of trace line?")
			}
		}
		break
	}
	return s
}

func (to *TestOutBase) Reset() {
	to.Buffer.Reset()
}

func (to *TestOutBase) MatchEmpty() {
	if to.Len() > 0 {
		to.Error("Unexpected unread content in trace out buffer:\n" + to.String())
	}
}

var re = regexp.MustCompile(
	//                                                       ---1----------  -2-----               -3-----------
	// 2016-04   -13    20   :17   :51    .123456| +  0200  |ERROR          |msg     |v2node      |comp          |1   |func        |file        .go           (32  )
	`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{6}\|[+-]\d{4}\|([[:alpha:]]+[ ]*)\|([^|]*)\|[[:word:]]+\|([[:word:]]+)\|\d+\|[[:word:]]+\|[[:word:]]+\.[[:alpha:]]+\(\d+\)`,
)

var tpid = strconv.Itoa(os.Getpid())

func (to *TestOutBase) Match(e Entry) {
	if to.Failed() {
		to.Log("-Matching starts")
	}
	lev2sev := func(lev *Level) *string {
		if lev != nil {
			sev := severityFormat(*lev)
			return &sev
		}
		return nil
	}
	var msg *string
	if e.Info != nil || e.Msg != nil {
		var m string
		msg = &m
		if e.Msg != nil {
			m += *e.Msg
		}
		if e.Info != nil {
			if e.Msg != nil {
				m += " "
			}
			m += *e.Info
		}
	}
	expected := []*string{lev2sev(e.Level), msg, e.Topic}
	s := to.ReadLine()
	matchIdxs := re.FindStringSubmatchIndex(s)
	if len(matchIdxs) <= 0 {
		to.Errorf("ERROR: Unexpected trace line. Line did not match: >>%s<<", s)
	}
	var subs []string
	for i := 2; i < len(matchIdxs); i += 2 {
		subs = append(subs, s[matchIdxs[i]:matchIdxs[i+1]])
	}
	if subs == nil || len(subs) != len(expected) {
		to.Errorf("ERROR: Unexpected trace line. No subs matched: >>%s<<", s)
	}
	if len(subs) > 0 {
		for i, x := range subs {
			if expected[i] != nil {
				if *expected[i] != x {
					to.Errorf("ERROR: Unexpected value at position %d: %q instead of %q", i, x, *expected[i])
				}
			}
		}
	}
	if to.Failed() {
		to.Log("-Matching context:")
		to.Log(string(debug.Stack()))
		to.Logf("Actual trace line: >>%s<<", s)
		to.Logf("Escaped trace line: >>%q<<", s)
		to.Log("expected:", ptr2Str(expected))
		to.Log("Matched", len(subs), "subs")
		to.Log("-Matching finished")
	}
}

func ptr2Str(ptrs []*string) []string {
	ret := make([]string, len(ptrs))
	for i, ptr := range ptrs {
		if ptr == nil {
			ret[i] = "<nil>"
		} else {
			ret[i] = *ptr
		}
	}
	return ret
}
