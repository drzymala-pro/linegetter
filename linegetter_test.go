package linegetter

import (
	"testing"
	"bytes"
	"io"
)


var linegetter_tests = []struct {
	idx int      // Index to show in error logs
	inp string   // input buffer in the form of string
	exp []string // expected lines as strings
}{
	{0,  "",      []string{""}},
	{1,  " ",     []string{" "}},
	{2,  "   ",   []string{"   "}},
	{3,  "a",     []string{"a"}},
	{4,  "abc",   []string{"abc"}},
	{5,  "\n",    []string{"", ""}},
	{6,  "abc\n", []string{"abc", ""}},
	{7,  "\nabc", []string{"", "abc"}},
	{8,  "\n\nabc", []string{"", "", "abc"}},
	{9,  "abc\n\n", []string{"abc", "", ""}},
	{10, "\n\nabc  \n\n", []string{"", "", "abc  ", "", ""}},
	{11, "\n\n  abc\n\n", []string{"", "", "  abc", "", ""}},
	{12, "abc\ndef\nghi",       []string{"abc", "def", "ghi"}},
	{13, "abc\ndef\nghi\n",     []string{"abc", "def", "ghi", ""}},
	{14, "abc\ndef\nghi\n\n\n", []string{"abc", "def", "ghi", "", "", ""}},
	{15, "abc\n\ndef\n\nghi",   []string{"abc", "", "def", "", "ghi"}},
	{16, "\nabc\ndef\nghi",     []string{"", "abc", "def", "ghi"}},
	{17, "\nabc\ndef\nghi\n",   []string{"", "abc", "def", "ghi", ""}},
	{18, "\nabc\n \ndef\nghi\n",   []string{"", "abc", " ", "def", "ghi", ""}},
	{19, "abc\ndef\nghi\n\n \n\n", []string{"abc", "def", "ghi", "", " ", "", ""}},
}


func TestInvalidParameter(t *testing.T) {
	var mock io.ReadSeeker = nil
	ilg, err := NewLineGetter(mock)
	if err == nil {
		t.Fatalf("Creating LineGetter with invalid argument does not return error.")
	}
	if ilg != nil {
		t.Fatalf("Creating LineGetter with invalid argument does not return nil.")
	}
}


func TestTableData(t *testing.T) {
	var e error
	var expct []string
	var line_number int64
	var input, got_line string
	var readskr io.ReadSeeker
	var linegetter *LineGetter
	for _, v := range linegetter_tests {
		input = v.inp
		expct = v.exp
		readskr = bytes.NewReader([]byte(input))
		linegetter, e = NewLineGetter(readskr)
		if e != nil {
			t.Fatalf("TestTableData idx:%v - Creating LineGetter returned error: %v", v.idx, e)
		}
		if c := linegetter.GetLineCount(); c != int64(len(expct)) {
			t.Fatalf("TestTableData idx:%v - LineGetter returned wrong number of lines: %v instead of: %v", v.idx, c, len(expct))
		}
		for j, w := range expct {
			line_number = int64(j) + 1
			got_line, e = linegetter.GetLine(line_number)
			if e != nil {
				t.Fatalf("TestTableData idx:%v - Got error: \"%v\" when reading line %v", v.idx, e, line_number)
			}
			if w != got_line {
				t.Fatalf("TestTableData idx:%v - Got wrong line number %v: \"%v\" instead of \"%v\"", v.idx, line_number, got_line, w)
			}
		}
	}
}

