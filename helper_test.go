package errors

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
	e.fileLine = ""
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
	e.fileLine = ""
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
	e.fileLine = ""
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
	e.fileLine = ""
	if !reflect.DeepEqual(*e, want) {
		t.Errorf("New() = %v, want %v", *e, want)
	}
}
func TestHTTPStatusCodeMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
		want2 bool
	}{
		{
			name: "TypeInternal",
			args: args{
				err: Internal("unknown error occurred"),
			},
			want:  http.StatusInternalServerError,
			want1: "unknown error occurred",
			want2: true,
		},
		{
			name: "TypeInternal - Go builtin error type",
			args: args{
				err: errors.New("unknown error occurred"),
			},
			want:  http.StatusInternalServerError,
			want1: "unknown error occurred",
			want2: false,
		},
		{
			name: "TypeValidation",
			args: args{
				err: Validation("invalid email provided"),
			},
			want:  http.StatusUnprocessableEntity,
			want1: "invalid email provided",
			want2: true,
		},
		{
			name: "TypeInputBody",
			args: args{
				err: InputBody("invalid json provided"),
			},
			want:  http.StatusBadRequest,
			want1: "invalid json provided",
			want2: true,
		},
		{
			name: "TypeDuplicate",
			args: args{
				err: Duplicate("duplicate content detected"),
			},
			want:  http.StatusConflict,
			want1: "duplicate content detected",
			want2: true,
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				err: Unauthenticated("authentication required"),
			},
			want:  http.StatusUnauthorized,
			want1: "authentication required",
			want2: true,
		},
		{
			name: "TypeUnauthorized",
			args: args{
				err: Unauthorized("not authorized to access this resource"),
			},
			want:  http.StatusForbidden,
			want1: "not authorized to access this resource",
			want2: true,
		},
		{
			name: "TypeEmpty",
			args: args{
				err: Empty("empty content not expected"),
			},
			want:  http.StatusGone,
			want1: "empty content not expected",
			want2: true,
		},
		{
			name: "TypeNotFound",
			args: args{
				err: NotFound("requested resource not found"),
			},
			want:  http.StatusNotFound,
			want1: "requested resource not found",
			want2: true,
		},
		{
			name: "TypeMaximumAttempts",
			args: args{
				err: MaximumAttempts("exceeded maximum number of requests allowed"),
			},
			want:  http.StatusTooManyRequests,
			want1: "exceeded maximum number of requests allowed",
			want2: true,
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				err: SubscriptionExpired("your subscription has expired"),
			},
			want:  http.StatusPaymentRequired,
			want1: "your subscription has expired",
			want2: true,
		},
		{
			name: "TypeDownstreamDependencyTimedout",
			args: args{
				err: DownstreamDependencyTimedout("dependency timed out"),
			},
			want:  http.StatusInternalServerError,
			want1: "dependency timed out",
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := HTTPStatusCodeMessage(tt.args.err)
			if got != tt.want {
				t.Errorf("HTTPStatusCodeMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("HTTPStatusCodeMessage() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("HTTPStatusCodeMessage() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
func TestHTTPStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want2 bool
	}{
		{
			name: "TypeInternal",
			args: args{
				err: Internal("unknown error occurred"),
			},
			want:  http.StatusInternalServerError,
			want2: true,
		},
		{
			name: "TypeInternal - Go builtin error type",
			args: args{
				err: errors.New("unknown error occurred"),
			},
			want:  http.StatusInternalServerError,
			want2: false,
		},
		{
			name: "TypeValidation",
			args: args{
				err: Validation("invalid email provided"),
			},
			want:  http.StatusUnprocessableEntity,
			want2: true,
		},
		{
			name: "TypeInputBody",
			args: args{
				err: InputBody("invalid json provided"),
			},
			want:  http.StatusBadRequest,
			want2: true,
		},
		{
			name: "TypeDuplicate",
			args: args{
				err: Duplicate("duplicate content detected"),
			},
			want:  http.StatusConflict,
			want2: true,
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				err: Unauthenticated("authentication required"),
			},
			want:  http.StatusUnauthorized,
			want2: true,
		},
		{
			name: "TypeUnauthorized",
			args: args{
				err: Unauthorized("not authorized to access this resource"),
			},
			want:  http.StatusForbidden,
			want2: true,
		},
		{
			name: "TypeEmpty",
			args: args{
				err: Empty("empty content not expected"),
			},
			want:  http.StatusGone,
			want2: true,
		},
		{
			name: "TypeNotFound",
			args: args{
				err: NotFound("requested resource not found"),
			},
			want:  http.StatusNotFound,
			want2: true,
		},
		{
			name: "TypeMaximumAttempts",
			args: args{
				err: MaximumAttempts("exceeded maximum number of requests allowed"),
			},
			want:  http.StatusTooManyRequests,
			want2: true,
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				err: SubscriptionExpired("your subscription has expired"),
			},
			want:  http.StatusPaymentRequired,
			want2: true,
		},
		{
			name: "TypeDownstreamDependencyTimedout",
			args: args{
				err: DownstreamDependencyTimedout("dependency timed out"),
			},
			want:  http.StatusInternalServerError,
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := HTTPStatusCode(tt.args.err)
			if got != tt.want {
				t.Errorf("HTTPStatusCodeMessage(), %s, got = %v, want %v", tt.name, got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("HTTPStatusCodeMessage(), %s, got2 = %v, want %v", tt.name, got2, tt.want2)
			}
		})
	}
}

func TestErrWithoutTrace(t *testing.T) {
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
			if got, _ := ErrWithoutTrace(e); got != tt.want {
				t.Errorf("Error.Message() = %v, want %v", got, tt.want)
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

func TestWriteHTTP(t *testing.T) {
	type args struct {
		err error
		w   http.ResponseWriter
	}
	tests := []struct {
		name        string
		args        args
		wantMessage string
		wantStatus  int
	}{
		{
			name: "TypeInternal",
			args: args{
				err: Internal("system error"),
				w:   httptest.NewRecorder(),
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "system error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteHTTP(tt.args.err, tt.args.w)
			rr, _ := tt.args.w.(*httptest.ResponseRecorder)
			if rr != nil {
				if rr.Code != tt.wantStatus {
					t.Errorf(
						"WriteHTTP() got = %d, '%s', expected %d, '%s'",
						rr.Code,
						rr.Body.String(),
						tt.wantStatus,
						tt.wantMessage,
					)
				}
			}
		})
	}
}
