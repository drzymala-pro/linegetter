// linegetter contains methods for efficient, random access reading
// from huge files containing lot of lines, typically log files.
package linegetter

import (
	"io"
	"bufio"
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
	total_line_count int64
	reader_seeker    io.ReadSeeker
	line_index       []int64
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
	lg := LineGetter{ total_line_count: 0, reader_seeker: rs }
	err := lg.reindex()
	if err != nil {
		return nil, err
	}
	return &lg, nil
}

// GetLineCount returns the number of lines available in the LineGetter.
func (lg *LineGetter) GetLineCount() int64 {
	return lg.total_line_count
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
	_, err := lg.reader_seeker.Seek(start_idx, io.SeekStart)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, final_len)
	n, err := io.ReadFull(lg.reader_seeker, buffer)
	if err != nil {
		return string(buffer[:n]), io.ErrUnexpectedEOF
	}
	if truncated {
		return string(buffer), ErrLineTruncated
	}
	return string(buffer), nil
}

func (lg *LineGetter) reindex() error {
	var current_position int64 = 0
	lg.total_line_count = 0
	// Rewind to the beginning of reader
	_, err := lg.reader_seeker.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// Naive approach - read one byte at a time
	for {
		n, b, e := read_next_byte(lg.reader_seeker)
		switch {
		case e != nil:
			// Got read error
			return e
		case n == 0:
			// Got EOF
			lg.line_index[lg.total_line_count] = current_position
			return nil
		default:
			handle_character(lg)
		}
	}
}


// Read until given delimiter is found
func read_until(delimiter byte, reader *Reader) (length int64, last_byte byte, err error) {
	// Naive approach, read one byte at a time
	err = nil
	length = 0
	last_byte = 0
	var buffer [1]byte = {0}
	for {
		n, err = reader.Read(buffer)
		switch {
		case n == 1:
			length += 1
		case n == 0 && err == nil:
			err = io.EOF
		case err != nil:
			return length, last_byte, err
		case buffer[0] == delimiter:
			return length, last_byte, nil
		default:
			last_byte = buffer[0]
			continue
		}
	}
}



{
	b := make([]byte, 1)
	p := 0
	for {
		n, err := lg.reader_seeker.Read(b)
		switch {
		case err == io.EOF:
			/* Finished reading */
			return nil
		case err != nil:
			/* Read error occurred */
			lg.total_line_count = 0
			return err
		case n == 1:
			/* Read success, handle the read data */
			handle_next_byte(lg)
		default:
			/* No data, try again */
			retry_count += 1
			continue
		}
	}
}

func read_next_byte(reader *io.Reader) (byte, error) {

}

func handle_next_byte(lg *LineGetter) {
	if lg.total_line_count == 0 {
	/* We got at least one line, even if it won't end with '\n' */
	lg.total_line_count += 1
	}
	/* Read success */
	if b[0] == '\n' {
		/* We've got a line separator, i.e. new line */
		lg.total_line_count += 1
	}
	switch b[0] {
	case '\n':
		break
	default:
		/* Got end of line */
	}
	read(b[0])
	/* Increase the cursor index */
	p += 1
}

