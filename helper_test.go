package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestWrap(t *testing.T) {
	err := errors.New("original error")
	message := "wrapped error"
	want := Error{
		original: err,
		message:  message,
		eType:    TypeInternal,
	}
	e := Wrap(err, message)
	e.pcs = nil
	e.pc = 0
	if !reflect.DeepEqual(*e, want) {
		t.Errorf("New() = %v, want %v", *e, want)
	}

	err = New("original error of type *Error")
	message = "wrapped error"
	want = Error{
		original: err,
		message:  message,
		eType:    defaultErrType,
	}
	e = Wrap(err, message)
	e.pcs = nil
	e.pc = 0
	if !reflect.DeepEqual(*e, want) {
		t.Errorf("New() = %v, want %v", *e, want)
	}
}
func TestWrapf(t *testing.T) {
	err := errors.New("original error")
	format := "%s prefixed"
	message := "wrapped error"
	want := Error{
		original: err,
		message:  fmt.Sprintf(format, message),
		eType:    TypeInternal,
	}
	e := Wrapf(err, format, message)
	e.pcs = nil
	e.pc = 0
	if !reflect.DeepEqual(*e, want) {
		t.Errorf("New() = %v, want %v", *e, want)
	}

	err = New("original error of type *Error")
	format = "%s prefixed"
	message = "wrapped error"
	want = Error{
		original: err,
		message:  fmt.Sprintf(format, message),
		eType:    defaultErrType,
	}
	e = Wrapf(err, format, message)
	e.pcs = nil
	e.pc = 0
	if !reflect.DeepEqual(*e, want) {
		t.Errorf("New() = %v, want %v", *e, want)
	}
}
func TestErrWithoutTrace(t *testing.T) {
	type fields struct {
		original error
		message  string
		eType    errType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "No nested error",
			fields: fields{
				original: nil,
				message:  "hello friendly error msg",
				eType:    TypeInternal,
			},
			want: "hello friendly error msg",
		},
		{
			name: "Empty message",
			fields: fields{
				original: nil,
				message:  "",
				eType:    TypeInternal,
			},
			want: "unknown error occurred",
		},
		{
			name: "Nested error with message",
			fields: fields{
				original: &Error{
					original: &Error{
						message: "",
						eType:   TypeInputBody,
					},
					message: "hello nested err message",
					eType:   TypeInternal,
				},
				message: "",
				eType:   TypeInternal,
			},
			want: "hello nested err message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				original: tt.fields.original,
				message:  tt.fields.message,
				eType:    tt.fields.eType,
			}
			if got, _ := ErrWithoutTrace(e); got != tt.want {
				t.Errorf("Error.ErrWithoutTrace() = %v, want %v", got, tt.want)
			}
		})
	}
	err := errors.New("std error")

	_, isTypeErr := ErrWithoutTrace(err)
	if isTypeErr {
		t.Error("ErrWithoutTrace() should return false if error is not of type *Error")
	}
}
func TestType(t *testing.T) {
	type fields struct {
		original error
		message  string
		eType    errType
	}
	tests := []struct {
		name   string
		fields fields
		want   errType
	}{
		{
			name: "TypeInternal",
			fields: fields{
				eType: TypeInternal,
			},
			want: TypeInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				original: tt.fields.original,
				message:  tt.fields.message,
				eType:    tt.fields.eType,
			}
			if got := Type(e); got != tt.want {
				t.Errorf("Error.Type() = %v, want %v", got, tt.want)
			}
		})
	}
	err := errors.New("std error")
	if int(Type(err)) != -1 {
		t.Errorf("Type() should return -1 if error is not of type *Error")
	}
}
func TestTypeInt(t *testing.T) {
	err := NewWithType("error", TypeValidation)
	if TypeInt(err) != int(TypeValidation) {
		t.Error("TypeInt() should return int(TypeValidation)")
	}
}
func TestHasType(t *testing.T) {
	type args struct {
		err error
		et  errType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has required type, nested",
			args: args{
				err: ValidationErr(
					DuplicateErr(
						Internal("hello world"),
						DefaultMessage,
					),
					DefaultMessage,
				),
				et: TypeInternal,
			},
			want: true,
		},
		{
			name: "has required type, not nested",
			args: args{
				err: Internal("hello world"),
				et:  TypeInternal,
			},
			want: true,
		},
		{
			name: "does not have required type",
			args: args{
				err: ValidationErr(
					DuplicateErr(
						Internal("hello world"),
						DefaultMessage,
					),
					DefaultMessage,
				),
				et: TypeInputBody,
			},
			want: false,
		},
		{
			name: "*Error wrapped in external error",
			args: args{
				err: fmt.Errorf("unknown error %w", Internal("internal error")),
				et:  TypeInternal,
			},
			want: true,
		},
		{
			name: "other error type",
			args: args{
				err: fmt.Errorf("external error"),
				et:  TypeInputBody,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasType(tt.args.err, tt.args.et); got != tt.want {
				t.Errorf("HasType() = %v, want %v %s", got, tt.want, tt.args.err.Error())
			}
		})
	}
}

func TestContextErrors(t *testing.T) {
	err := context.Canceled
	werr := Wrap(err, "wrapped")
	if !HasType(werr, TypeContextCancelled) {
		t.Error("expected TypeContextCancelled")
	}

	code, _ := HTTPStatusCode(werr)
	if code != http.StatusRequestTimeout {
		t.Errorf("expected 408, got: %d", code)
	}

	gcode, _ := GRPCStatusCode(werr)
	if gcode != codes.Canceled {
		t.Errorf("expected %d/%s, got: %d", codes.Canceled, codes.Canceled, code)
	}

	err = context.DeadlineExceeded
	werr = Wrap(err, "wrapped")

	gcode, _ = GRPCStatusCode(werr)
	if gcode != codes.DeadlineExceeded {
		t.Errorf("expected %d/%s, got: %d", codes.DeadlineExceeded, codes.DeadlineExceeded, code)
	}
}
