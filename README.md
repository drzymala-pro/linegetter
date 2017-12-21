# linegetter
Package containing methods for fast random access reads from huge files, typically log files.

## Interfaces

### LineGetter
LineGetter implements random access to lines from io.ReadSeeker object.
Lines are separated with ASCII line feed character, 0x0A in hex.
```go
type LineGetter struct {...}
```

## Functions

### NewLineGetter
NewLineGetter returns a new LineGetter.
Upon creation of LineGetter, the whole ReadSeeker is scanned.
If the io.ReadSeeker is nil, nil is returned.
```go
func NewLineGetter(rs io.ReadSeeker) *LineGetter {...}
```

### GetLineCount
GetLineCount returns the number of lines available in the LineGetter.
```go
func (lg *LineGetter) GetLineCount() int64 {...}
```

### GetLine
GetLine returns the n-th line from the LineGetter.
 * If the line number is out of range, ErrInvalidArgument is returned.
 * If some error happens during reading, the error is returned and
   the resulting string does not contain the full expected length.
 * If the line length exceeds MaxLineLength, ErrLineTruncated is returned
   and the resulting string is truncated to MaxLineLength size.
```go
func (lg *LineGetter) GetLine(ln int64) (string, error) {...}
```

## Errors

```go
var (
    ErrInvalidArgument = errors.New("invalid argument")
    ErrLineTruncated   = errors.New("line truncated")
)
```

## Constants

```go
const (
    read_chunk_sz int64 = 16383 // 0x3FFF
    MaxLineLength int64 = read_chunk_sz
)
```

