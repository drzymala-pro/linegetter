package linegetter

import (
	"testing"
	"bytes"
	"io"
)

var mock io.ReadSeeker


func make_line_getter_or_die(t *testing.T, rs io.ReadSeeker) *LineGetter {
	lg, err := NewLineGetter(rs)
	if err != nil {
		t.Fatalf("Creating LineGetter returned error.")
	}
	if lg == nil {
		t.Fatalf("Creating LineGetter returned nil.")
	}
	return lg
}


func TestInvalidParameter(t *testing.T) {
	mock = nil
	ilg, err := NewLineGetter(mock)
	if err == nil {
		t.Fatalf("Creating LineGetter with invalid argument does not return error.")
	}
	if ilg != nil {
		t.Fatalf("Creating LineGetter with invalid argument does not return nil.")
	}
}

type emptyReaderSeeker struct {}
func (r *emptyReaderSeeker) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}
func (r *emptyReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func TestEmptyReader(t *testing.T) {
	mock := &emptyReaderSeeker{}
	elg := make_line_getter_or_die(t, mock)
	c := elg.GetLineCount()
	if c != 0 {
		t.Fatalf("Empty LineGetter returns non zero line count: %v", c)
	}
	line, err := elg.GetLine(0)
	if line != "" {
		t.Fatalf("Empty LineGetter has returned a non empty line: %v", line)
	}
	if err == nil {
		t.Fatalf("Empty LineGetter has not returned error.")
	}
}


func TestSingleByteReader(t *testing.T) {
	mock = bytes.NewReader([]byte("G"))
	lg := make_line_getter_or_die(t, mock)
	c := lg.GetLineCount()
	if c != 1 {
		t.Fatalf("LineGetter returned wrong number of lines: %v", c)
	}
	line, err := lg.GetLine(1)
	if err != nil {
		t.Fatalf("LineGetter has returned error: %v", err)
	}
	if line != "G" {
		t.Fatalf("LineGetter has returned wrong line: \"%s\"", line)
	}
}


func TestSingleLineReader(t *testing.T) {
	text := "aaaaaaaaaaa"
	mock = bytes.NewReader([]byte(text))
	lg := make_line_getter_or_die(t, mock)
	c := lg.GetLineCount()
	if c != 1 {
		t.Fatalf("LineGetter returned wrong number of lines: %v", c)
	}
	line, err := lg.GetLine(1)
	if err != nil {
		t.Fatalf("LineGetter has returned error: %v", err)
	}
	if line != text {
		t.Fatalf("LineGetter has returned wrong line: \"%s\"", line)
	}
}


func TestDoubleLineReader(t *testing.T) {
	text := "aaaaaaaaaaa\nbbbbbbbbbb"
	mock = bytes.NewReader([]byte(text))
	lg := make_line_getter_or_die(t, mock)
	c := lg.GetLineCount()
	if c != 2 {
		t.Fatalf("LineGetter returned wrong number of lines: %v", c)
	}
	line1, err1 := lg.GetLine(1)
	if err1 != nil {
		t.Fatalf("LineGetter has returned error: %v", err1)
	}
	if line1 != "aaaaaaaaaaa" {
		t.Fatalf("LineGetter has returned wrong line: \"%s\"", line1)
	}
	line2, err2 := lg.GetLine(2)
	if err2 != nil {
		t.Fatalf("LineGetter has returned error: %v", err2)
	}
	if line2 != "bbbbbbbbbb" {
		t.Fatalf("LineGetter has returned wrong line: \"%s\"", line2)
	}
}

