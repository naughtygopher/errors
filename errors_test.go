// Package errors helps in wrapping errors with custom type as well as a user friendly message. This is particularly useful when responding to APIs
package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {

	bar := func() error {
		return fmt.Errorf("hello %s", "world!")
	}

	foo := func() error {
		err := bar()
		if err != nil {
			return InternalErr(err, "bar is not happy")
		}
		return nil
	}

	err := foo()

	got := fmt.Sprintf("%+v", err)
	want := "errors/errors_test.go:21: bar is not happy\nhello world!"
	if !strings.Contains(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}

	got = fmt.Sprintf("%v", err)
	want = "bar is not happy"
	if !strings.Contains(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}

	got = fmt.Sprintf("%+s", err)
	want = "bar is not happy: hello world!"
	if !strings.Contains(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}

	got = fmt.Sprintf("%s", err)
	want = "bar is not happy"
	if !strings.Contains(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}

}
func TestErrorWithoutFileLine(t *testing.T) {
	err := New("error without file line")
	want := "error without file line"
	got := err.ErrorWithoutFileLine()
	if got != want {
		t.Errorf("ErrorWithoutFileLine() = %v\nwant %v", got, want)
	}

	err = Wrap(err, "wrapped error")
	want = "wrapped error: error without file line"
	got = err.ErrorWithoutFileLine()
	if got != want {
		t.Errorf("ErrorWithoutFileLine() = %v\nwant %v", got, want)
	}

	err = Wrap(errors.New("std err"), "wrapped std error")
	want = "wrapped std error: std err"
	got = err.ErrorWithoutFileLine()
	if got != want {
		t.Errorf("ErrorWithoutFileLine() = %v\nwant %v", got, want)
	}

	err = Wrap(errors.New("std err"), "")
	want = "std err"
	got = err.ErrorWithoutFileLine()
	if got != want {
		t.Errorf("ErrorWithoutFileLine() = %v\nwant %v", got, want)
	}

	err = New("")
	got = err.ErrorWithoutFileLine()
	if !strings.Contains(got, "errors/errors_test.go:") {
		t.Errorf("empty error should have fileline: %s", got)
	}
}
func TestNew(t *testing.T) {
	message := "friendly error message"
	want := Error{
		message: message,
		eType:   defaultErrType,
	}
	e := New(message)
	e.pcs = nil
	e.pc = 0

	if !reflect.DeepEqual(*e, want) {
		t.Errorf("New() = %v\nwant %v", *e, want)
	}
}

func TestErrorf(t *testing.T) {
	format := "%s prefixed"
	message := "friendly error message"
	want := Error{
		message: fmt.Sprintf(format, message),
		eType:   defaultErrType,
	}
	e := Errorf(format, message)
	e.pcs = nil
	e.pc = 0

	if !reflect.DeepEqual(*e, want) {
		t.Fail()
	}
}

func TestSetDefaultType(t *testing.T) {
	type args struct {
		message string
		e       errType
	}
	tests := []struct {
		name        string
		args        args
		wantErrType errType
	}{
		{
			name: "TypeInputBody",
			args: args{
				e: TypeInputBody,
			},
			wantErrType: TypeInputBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := defaultErrType
			SetDefaultType(tt.args.e)
			err := New(tt.args.message)
			// resetting to previous value to stop messing with the entire package
			SetDefaultType(before)
			if err.Type() != tt.wantErrType {
				t.Errorf(
					"New() = got type '%d', expected '%d",
					err.Type(),
					tt.wantErrType,
				)
			}
		})
	}
}
func TestStacktrace(t *testing.T) {
	err := errors.New("original error")
	e := Wrap(err, "wrapped error")
	got := Stacktrace(e)
	// silly way of verifying the stacktrace is correct, excluding filepaths
	strings.Contains(got, "errors.TestStacktrace(): wrapped error")
	strings.Contains(got, "errors/errors_test.go:76")
	strings.Contains(got, "original error")
}
func TestStacktraceNoFormat(t *testing.T) {
	err := errors.New("original error")
	e := Wrap(err, "wrapped error")
	got := strings.Join(StacktraceNoFormat(e), "#")
	// silly way of verifying the stacktrace is correct, excluding filepaths
	strings.Contains(got, "errors.TestStacktrace(): wrapped error")
	strings.Contains(got, "errors/errors_test.go:76")
	strings.Contains(got, "original error")
	if strings.Contains(got, "\n") {
		t.Error("StacktraceNoFormat() should not contain newlines")
	}
}
func TestStacktraceCustomFormat(t *testing.T) {
	err := errors.New("original error")
	e := Wrap(err, "wrapped error")
	msgFormat := "message: %m#"
	traceFormat := "function: %f|"
	got := StacktraceCustomFormat(msgFormat, traceFormat, e)
	want := "message: wrapped error#function: github.com/bnkamalesh/errors.TestStacktraceCustomFormat|function: testing.tRunner|message: original error#"
	if got != want {
		t.Errorf("StacktraceCustomFormat() = %v\nwant %v", got, want)
	}
}
