// linegetter contains methods for efficient, random access reading
// from huge files containing lot of lines, typically log files.
package linegetter

import (
	"io"
	"errors"
)

const (
	read_chunk_sz := 16383 // 0x3FFF
	MaxLineLength := read_chunk_sz
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrLineTruncated   = errors.New("line truncated")
)

// LineGetter implements random access to lines from io.ReadSeeker object.
// Lines are separated with ASCII line feed character, 0x0A in hex.
type LineGetter struct {
	total_line_count uint64
	reader_seeker    io.ReadSeeker
}

// NewLineGetter returns a new LineGetter.
// Upon creation of LineGetter, the whole ReadSeeker is scanned.
// If the io.ReadSeeker is nil, ErrInvalidArgument error is returned.
func NewLineGetter(rs *io.ReadSeeker) *LineGetter {
	if rs == nil {
		return nil, ErrInvalidArgument
	}
	lg := LineGetter{ total_line_count: 0, reader_seeker: rs }
	lg.reindex()
	return &lg, nil
}

// GetLineCount returns the number of lines available in the LineGetter.
func (lg *LineGetter) GetLineCount() uint64 {
	return lg.total_line_count
}

// GetLine returns the n-th line from the LineGetter.
// * If the line number is out of range, ErrInvalidArgument is returned.
// * If some error happens during reading, the error is returned and
//   the resulting string does not contain the full expected length.
// * If the line length exceeds MaxLineLength, ErrLineTruncated is returned
//   and the resulting string is truncated to MaxLineLength size.
func (lg *LineGetter) GetLine(ln uint64) (string, error) {
	if ln >= lg.total_line_count {
		return "", ErrInvalidArgument
	}
	final_len := 0
	truncated := false
	start_idx := lg.li[ln]
	end_index := lg.li[ln+1]
	if end_index - start_idx > MaxLineLength {
		truncated = true
		final_len = MaxLineLength
		end_index = start_idx + final_len
	} else {
		final_len = end_index - start_idx
	}
	buffer = make([]byte, final_len)
	_, err := lg.reader_seeker.Seek(start_idx, io.SeekStart)
	if err != nil {
		return "", err
	}
	n, err = io.ReadFull(lg.reader_seeker, buffer)
	if err != nil {
		return buffer[:n], io.ErrUnexpectedEOF
	}
	if truncated {
		return buffer, ErrLineTruncated
	}
	return buffer, nil
}


func (lg *LineGetter) reindex() {
	lg.lineCount = 0
	for {

	}
}
