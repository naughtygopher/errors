// Package errors helps in wrapping errors with custom type as well as a user friendly message. This is particularly useful when responding to APIs
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func Test_errType_Int(t *testing.T) {
	tests := []struct {
		name string
		e    errType
		want int
	}{
		{
			name: "valid",
			e:    TypeInternal,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Int(); got != tt.want {
				t.Errorf("errType.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	sampleErrContent := "foo bar"
	sampleError, file, line := sampleError(sampleErrContent, TypeInternal)

	type fields struct {
		original error
		message  string
		eType    errType
		fileLine string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "single error, no original error, no fileline",
			fields: fields{
				original: nil,
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "",
			},
			want: " hello world",
		},
		{
			name: "single error, no original error, has fileline",
			fields: fields{
				original: nil,
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "/home/user/main.go:60",
			},
			want: "/home/user/main.go:60 hello world",
		},
		{
			name: "with original error",
			fields: fields{
				original: errors.New("bad error"),
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "/home/user/main.go:60",
			},
			want: "/home/user/main.go:60 bad error",
		},
		{
			name: "with original error of type *Error",
			fields: fields{
				original: sampleError,
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "",
			},
			want: fmt.Sprintf(" %s:%d %s", file, line, sampleErrContent),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				original: tt.fields.original,
				message:  tt.fields.message,
				eType:    tt.fields.eType,
				fileLine: tt.fields.fileLine,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sampleError(content string, et errType) (*Error, string, int) {
	err := NewWithType(content, et)
	_, file, line, _ := runtime.Caller(0) // calling right here to get the correct filename and linenum
	line--                                // since the expected line number is -1 where New was called
	return err, file, line
}

func TestError_Type(t *testing.T) {
	type fields struct {
		original error
		message  string
		eType    errType
		fileLine string
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
				fileLine: tt.fields.fileLine,
			}
			if got := e.Type(); got != tt.want {
				t.Errorf("Error.Type() = %v, want %v", got, tt.want)
			}
		})
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
			SetDefaultType(tt.args.e)
			err := New(tt.args.message)
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
