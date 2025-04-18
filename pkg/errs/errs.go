package errs

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type StackError struct {
	msg string
	pcs []uintptr
}

func (se *StackError) Error() string {
	sb := &strings.Builder{}
	sb.WriteString(se.msg)
	sb.WriteString("\nStack trace:\n")
	frames := runtime.CallersFrames(se.pcs)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(sb, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return sb.String()
}

func (se *StackError) Msg() string {
	return se.msg
}

// Unwrap so you can still use errors.Is/As
func (se *StackError) Unwrap() error { return nil }

func New(format string, args ...any) error {
	// Skip 2 frames: runtime.Callers + this New()
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])
	se := &StackError{
		msg: fmt.Sprintf(format, args...),
		pcs: pcs[:n],
	}
	return se
}

func Wrap(err error, format string, args ...any) error {
	format += ": %w"
	args = append(args, err)
	return New(format, args...)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
