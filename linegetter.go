// linegetter contains methods for efficient, random access reading
// from huge files containing lot of lines, typically log files.
package linegetter

import (
	"io"
	"errors"
)

const (
	read_chunk_sz int64 = 16383 // 0x3FFF
	MaxLineLength int64 = read_chunk_sz
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrLineTruncated   = errors.New("line truncated")
)

// LineGetter implements random access to lines from io.ReadSeeker object.
// Lines are separated with ASCII line feed character, 0x0A in hex.
type LineGetter struct {
	total_line_count int64
	reader_seeker    io.ReadSeeker
	line_index       []int64
}

// NewLineGetter returns a new LineGetter.
// Upon creation of LineGetter, the whole ReadSeeker is scanned.
// If the io.ReadSeeker is nil, nil is returned.
func NewLineGetter(rs io.ReadSeeker) *LineGetter {
	if rs == nil {
		return nil
	}
	lg := LineGetter{ total_line_count: 0, reader_seeker: rs }
	lg.reindex()
	return &lg
}

// GetLineCount returns the number of lines available in the LineGetter.
func (lg *LineGetter) GetLineCount() int64 {
	return lg.total_line_count
}

// GetLine returns the n-th line from the LineGetter.
// * If the line number is out of range, ErrInvalidArgument is returned.
// * If some error happens during reading, the error is returned and
//   the resulting string does not contain the full expected length.
// * If the line length exceeds MaxLineLength, ErrLineTruncated is returned
//   and the resulting string is truncated to MaxLineLength size.
func (lg *LineGetter) GetLine(ln int64) (string, error) {
	if ln >= lg.total_line_count {
		return "", ErrInvalidArgument
	}
	var final_len int64
	var truncated bool  = false
	var start_idx int64 = lg.line_index[ln]
	var end_index int64 = lg.line_index[ln+1]
	if end_index - start_idx > MaxLineLength {
		truncated = true
		final_len = MaxLineLength
		end_index = start_idx + final_len
	} else {
		final_len = end_index - start_idx
	}
	buffer := make([]byte, final_len)
	_, err := lg.reader_seeker.Seek(start_idx, io.SeekStart)
	if err != nil {
		return "", err
	}
	n, err := io.ReadFull(lg.reader_seeker, buffer)
	if err != nil {
		return string(buffer[:n]), io.ErrUnexpectedEOF
	}
	if truncated {
		return string(buffer), ErrLineTruncated
	}
	return string(buffer), nil
}


func (lg *LineGetter) reindex() {
	lg.total_line_count = 0
	for {

	}
}
