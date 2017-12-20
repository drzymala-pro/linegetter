package linegetter

import (
	"fmt"
	"io"
	"errors"
)

// LineGetter is an interface for retrieving concrete line from a file.
// Lines can be accessed in random order. The delimiter for line separation
// is the newline symbol. Newline symbol is platform specific.
type LineGetter interface {
	GetLineCount() uint64
	GetLine(line uint64) (string, error)
}

// NewLineGetter creates a LineGetter from an io.ReadSeeker.
// Whole ReadSeeker is scanned and indexed for GetLine speed.
func NewLineGetter(rs io.ReadSeeker) (*lineGetter, error) {
	if rs == nil {
		return nil, errors.New("Nil value as argument to NewLineGetter()")
	}
	lg := lineGetter{}
	/* TODO: Implement me! */
	return &lg, nil
}

// GetLineCount returns the number of lines available in the LineGetter.
// Lines are split by newline sybmols. New line sybol is platform dependent.
func (lg *lineGetter) GetLineCount() uint64 {
	return lg.lineCount
}

// GetLine returns the given line without the trailing newline symbols.
// If the line number is not in range then error is returned.
func (lg *lineGetter) GetLine(line uint64) (string, error) {
	if line >= lg.GetLineCount() {
		return "", errors.New(fmt.Sprintf("Line number out of range: %d/%d.", line, lg.GetLineCount()))
	}
	return "", nil
}

// lineGetter is the internal structure holding LineGetter data.
type lineGetter struct {
	lineCount uint64
	rs io.ReadSeeker
}
