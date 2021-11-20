package errstr

import (
	"errors"
	"fmt"
	"runtime"
)

type frame [3]uintptr

type errorStr struct {
	message string
	error   error
	frame   frame
}

func caller(skip int) (r frame) {
	runtime.Callers(skip+1, r[:])
	return
}

func (f frame) location() (function, file string, line int) {
	frames := runtime.CallersFrames(f[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

func (e *errorStr) Error() string {
	return e.message
}

func (e *errorStr) String() string {
	return e.message + " [" + e.frame.String() + "]"
}

type ErrorStr interface {
	error
	String() string
	Location() (function, file string, line int)
}

func (e *errorStr) Location() (function, file string, line int) {
	return e.frame.location()
}

func (f frame) String() (str string) {
	function, file, line := f.location()
	str = function
	if str == "" {
		str = "<unknown func>"
	}
	if file != "" {
		str += fmt.Sprintf(" %s:%d\n", file, line)
	}
	return
}

func (e *errorStr) Unwrap() error {
	return e.error
}

func (e *errorStr) Is(err error) bool {
	return errors.Is(e.error, err)
}

/*
Format formats new error
*/
func Format(skip int, f string, a ...interface{}) error {
	return &errorStr{fmt.Sprintf(f, a...), nil, caller(skip + 1)}
}

/*
Wrapf wraps error with formatted string
*/
func Wrapf(skip int, err error, f string, a ...interface{}) error {
	return &errorStr{fmt.Sprintf(f, a...), err, caller(skip + 1)}
}

/*
Wrap wraps error with string message
*/
func Wrap(err error, msg string, skip ...int) error {
	s := fnz(skip...)
	return &errorStr{msg, err, caller(s + 1)}
}

/*
New creates new error
*/
func New(msg string, skip ...int) error {
	s := fnz(skip...)
	return &errorStr{msg, nil, caller(s + 1)}
}

func Errorf(f string, a ...interface{}) error {
	return Format(1, f, a...)
}

/*
PanicMessage returns a message from the panic object
*/
func MessageOf(e interface{}) string {
	if p, ok := e.(error); ok {
		return p.Error()
	}
	return fmt.Sprint(e)
}

/*
PanicError returns an error from the panic object
*/
func ErrorOf(e interface{}) error {
	if p, ok := e.(error); ok {
		return p
	}
	return errors.New(fmt.Sprint(e))
}

/*
Fnz returns first non-zero int
*/
func fnz(a ...int) int {
	for _, x := range a {
		if x != 0 { return x }
	}
	return 0
}
