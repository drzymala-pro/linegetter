package linegetter

import (
	"testing"
	"bytes"
	"io"
)

var mock io.ReadSeeker


func TestInvalidParameter(t *testing.T) {
	mock = nil
	ilg := NewLineGetter(mock)
	if ilg != nil {
		t.Fatalf("Creating LineGetter with invalid argument does not return nil.")
	}
}

func TestEmptyReader(t *testing.T) {
	mock = bytes.NewReader([]byte(""))
	elg := NewLineGetter(mock)
	c := elg.GetLineCount()
	if c != 0 {
		t.Fatalf("Empty LineGetter returns non zero line count: %v", c)
	}
	line, err := elg.GetLine(0)
	if line != "" {
		t.Fatalf("Empty LineGetter has returned a non empty line: %v", line)
	}
	if err != io.EOF {
		t.Fatalf("Empty LineGetter has not returned EOF error.")
	}
}


func TestSingleByteReader(t *testing.T) {
	mock = bytes.NewReader([]byte("G"))
	lg := NewLineGetter(mock)
	c := lg.GetLineCount()
	if c != 1 {
		t.Fatalf("LineGetter returned wrong number of lines: %v", c)
	}
	line, err := lg.GetLine(0)
	if line != "G" {
		t.Fatalf("LineGetter has returned wrong line: %v", line)
	}
	if err != nil {
		t.Fatalf("LineGetter has returned error: %v", err)
	}
}


