# linegetter
Package containing methods for fast random access reads from huge files, typically log files.

## Interfaces

### LineGetter
LineGetter is an interface for retrieving concrete line from a file.
Lines can be accessed in random order. The delimiter for line separation
is the newline symbol. Newline symbol is platform specific.
```go
type LineGetter interface {
    GetLineCount() uint64
    GetLine(line uint64) (string, error)
}
```

## Functions

### NewLineGetter
NewLineGetter creates a LineGetter from an io.ReadSeeker.
Whole ReadSeeker is scanned and indexed for GetLine speed.
```go
func NewLineGetter(rs io.ReadSeeker) (*lineGetter, error) {...}
```

### GetLineCount
GetLineCount returns the number of lines available in the LineGetter.
Lines are split by newline sybmols. New line sybol is platform dependent.
```go
func (lg *lineGetter) GetLineCount() uint64 {...}
```

### GetLine
GetLine returns the given line without the trailing newline symbols.
If the line number is not in range then error is returned.
```go
func (lg *lineGetter) GetLine(line uint64) (string, error) {...}
```
