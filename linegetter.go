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
type LineGetter struct {
	read_skr io.ReadSeeker
	line_cnt int64
	line_pos []int64
}


// NewLineGetter returns a new LineGetter.
// If the io.ReadSeeker is nil, ErrInvalidArgument is returned.
// Upon creation of LineGetter, the whole ReadSeeker is scanned.
// If during scanning there is an error other than EOF, the error is returned.
// Resulting LineGetter is not nil only if there is no error.
func NewLineGetter(rs io.ReadSeeker) (*LineGetter, error) {
	if rs == nil {
		return nil, ErrInvalidArgument
	}
	lg := LineGetter{ line_cnt: 0, read_skr: rs }
	err := lg.reindex()
	if err != nil {
		return nil, err
	}
	return &lg, nil
}


// GetLineCount returns the number of lines available in the LineGetter.
func (lg *LineGetter) GetLineCount() int64 {
	return lg.line_cnt
}


// GetLine returns the n-th line from the LineGetter.
// Lines are separated with ASCII line feed character, 0x0A in hex.
// The line separator is not included in the resulting lines.
// * If the line number is out of range, ErrInvalidArgument is returned.
// * If some error happens during reading, the error is returned and
//   the resulting string does not contain the full expected length.
// * If the line length exceeds MaxLineLength, ErrLineTruncated is returned
//   and the resulting string is truncated to MaxLineLength size.
func (lg *LineGetter) GetLine(ln int64) (string, error) {
	if ln >= lg.line_cnt {
		return "", ErrInvalidArgument
	}
	var final_len int64
	var truncated bool  = false
	var start_idx int64 = lg.line_pos[ln]
	var end_index int64 = lg.line_pos[ln+1]
	if end_index - start_idx > MaxLineLength {
		truncated = true
		final_len = MaxLineLength
		end_index = start_idx + final_len
	} else {
		final_len = end_index - start_idx
	}
	_, err := lg.read_skr.Seek(start_idx, io.SeekStart)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, final_len)
	n, err := io.ReadFull(lg.read_skr, buffer)
	if err != nil {
		return string(buffer[:n]), io.ErrUnexpectedEOF
	}
	if truncated {
		return string(buffer), ErrLineTruncated
	}
	return string(buffer), nil
}


func (lg *LineGetter) reindex() error {
	// Reset the line getter and rewind the reader
	if err := lg.reset(); err != nil {
		return err
	}
	// Naive approach - scan one byte at a time
	var current_pos int64 = 0
	for {
		data, err := read_next_byte(lg.read_skr)
		switch err {
		case nil:
			current_pos += 1
			if data == '\n' {
				lg.line_cnt += 1
				lg.line_pos = append(lg.line_pos, current_pos)
			}
		case io.EOF:
			// Scanned the whole thing. Mark the length of the last line.
			// Add the position of cursor but do not increase line count.
			lg.line_pos = append(lg.line_pos, current_pos)
			return nil
		default:
			// Unexpected error
			return err
		}
	}
}


func (lg *LineGetter) reset() error {
	lg.line_cnt = 0
	lg.line_pos = []int64{}
	if _, err := lg.read_skr.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return nil
}


// read_next_byte returns one valid byte or error, but never both.
func read_next_byte(reader io.Reader) (byte, error) {
	var p []byte
	var n int
	var err error
	p = make([]byte, 1)
	for {
		n, err = reader.Read(p)
		switch {
		case n > 0:
			// If any data available, ignore errors
			return p[0], nil
		case err != nil:
			// If no data but error, return error
			return 0, err
		default:
			// Otherwise try reading again
			continue
		}
	}
}

