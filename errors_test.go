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
			want: ": hello world",
		},
		{
			name: "single error, no original error, has fileline",
			fields: fields{
				original: nil,
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "/home/user/main.go:60",
			},
			want: "/home/user/main.go:60: hello world",
		},
		{
			name: "with original error",
			fields: fields{
				original: errors.New("bad error"),
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "/home/user/main.go:60",
			},
			want: "/home/user/main.go:60: hello world\nbad error",
		},
		{
			name: "with original error of type *Error",
			fields: fields{
				original: sampleError,
				message:  "hello world",
				eType:    TypeInternal,
				fileLine: "",
			},
			want: fmt.Sprintf(": hello world\n%s:%d: %s", file, line, sampleErrContent),
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
				t.Errorf("Error.Error() = '%v', want '%v'", got, tt.want)
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

func Benchmark_Internal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Internal("hello world")
	}
}
func Benchmark_Internalf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Internalf("%s prefixed", "hello world")
	}
}

func Benchmark_InternalErr(b *testing.B) {
	err := errors.New("bad error")
	for i := 0; i < b.N; i++ {
		InternalErr(err, "hello world")
	}
}

func Benchmark_InternalGetError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Internal("hello world").Error()
	}
}
func Benchmark_InternalGetErrorWithNestedError(b *testing.B) {
	err := errors.New("bad error")
	for i := 0; i < b.N; i++ {
		_ = InternalErr(err, "hello world").Error()
	}
}

func Benchmark_InternalGetMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Internal("hello world").Message()
	}
}

func Benchmark_InternalGetMessageWithNestedError(b *testing.B) {
	err := New("bad error")
	for i := 0; i < b.N; i++ {
		_ = InternalErr(err, "hello world").Message()
	}
}

func Benchmark_HTTPStatusCodeMessage(b *testing.B) {
	// SubscriptionExpiredErr is the slowest considering it's the last item in switch case
	err := SubscriptionExpiredErr(SubscriptionExpired("old"), "expired")
	for i := 0; i < b.N; i++ {
		_, _, _ = HTTPStatusCodeMessage(err)
	}
}

func TestError_Message(t *testing.T) {
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
			name: "No nested error",
			fields: fields{
				original: nil,
				message:  "hello friendly error msg",
				eType:    TypeInternal,
				fileLine: "errors.go:87",
			},
			want: "hello friendly error msg",
		},
		{
			name: "Empty message",
			fields: fields{
				original: nil,
				message:  "",
				eType:    TypeInternal,
				fileLine: "errors.go:87",
			},
			want: "errors.go:87: unknown error occurred",
		},
		{
			name: "Nested error with message",
			fields: fields{
				original: &Error{
					original: &Error{
						message:  "",
						eType:    TypeInputBody,
						fileLine: "errors.go:87",
					},
					message:  "hello nested err message",
					eType:    TypeInternal,
					fileLine: "errors.go:87",
				},
				message:  "",
				eType:    TypeInternal,
				fileLine: "errors.go:87",
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
				fileLine: tt.fields.fileLine,
			}
			if got := e.Message(); got != tt.want {
				t.Errorf("Error.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
