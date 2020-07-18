package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

func TestHelperFnsForAllTypes(t *testing.T) {
	type args struct {
		message string
		eType   errType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TypeInternal",
			args: args{
				eType:   TypeInternal,
				message: "internal error",
			},
		},
		{
			name: "TypeValidation",
			args: args{
				eType:   TypeValidation,
				message: "validation error",
			},
		},
		{
			name: "TypeInputBody",
			args: args{
				eType:   TypeInputBody,
				message: "invalid input body",
			},
		},
		{
			name: "TypeDuplicate",
			args: args{
				eType:   TypeDuplicate,
				message: "duplicate contennt",
			},
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				eType:   TypeUnauthenticated,
				message: "not authenticated",
			},
		},
		{
			name: "TypeUnauthorized",
			args: args{
				eType:   TypeUnauthorized,
				message: "not authorized",
			},
		},
		{
			name: "TypeEmpty",
			args: args{
				eType:   TypeEmpty,
				message: "empty content",
			},
		},
		{
			name: "TypeNotFound",
			args: args{
				eType:   TypeNotFound,
				message: "resource not found",
			},
		},
		{
			name: "TypeMaximumAttempts",
			args: args{
				eType:   TypeMaximumAttempts,
				message: "exceeded maximum number of allowed attempts",
			},
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				eType:   TypeSubscriptionExpired,
				message: "subscription expired",
			},
		},
		{
			name: "TypeDownstreamDependencyTimedout",
			args: args{
				eType:   TypeDownstreamDependencyTimedout,
				message: "downstream dependency call timed out",
			},
		},
	}
	for _, tt := range tests {
		var want, got *Error
		switch tt.args.eType {
		case TypeInternal:
			{
				got = Internal(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeInternal)
			}
		case TypeValidation:
			{
				got = Validation(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeValidation)
			}
		case TypeInputBody:
			{
				got = InputBody(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeInputBody)
			}
		case TypeDuplicate:
			{
				got = Duplicate(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeDuplicate)
			}
		case TypeUnauthenticated:
			{
				got = Unauthenticated(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeUnauthenticated)
			}
		case TypeUnauthorized:
			{
				got = Unauthorized(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeUnauthorized)
			}

		case TypeEmpty:
			{
				got = Empty(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeEmpty)
			}

		case TypeNotFound:
			{
				got = NotFound(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeNotFound)
			}

		case TypeMaximumAttempts:
			{
				got = MaximumAttempts(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeMaximumAttempts)
			}
		case TypeSubscriptionExpired:
			{
				got = SubscriptionExpired(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeSubscriptionExpired)
			}
		case TypeDownstreamDependencyTimedout:
			{
				got = DownstreamDependencyTimedout(tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(nil, tt.args.message, file, line, TypeDownstreamDependencyTimedout)
			}

		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("%v() = %v, want %v", tt.args.eType, got, want)
		}
	}
}

func TestHelperFnsForAllTypesWithOriginalError(t *testing.T) {
	originalErr := errors.New("error returned by some other package")
	type args struct {
		message string
		eType   errType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TypeInternal",
			args: args{
				eType:   TypeInternal,
				message: "internal error",
			},
		},
		{
			name: "TypeValidation",
			args: args{
				eType:   TypeValidation,
				message: "validation error",
			},
		},
		{
			name: "TypeInputBody",
			args: args{
				eType:   TypeInputBody,
				message: "invalid input body",
			},
		},
		{
			name: "TypeDuplicate",
			args: args{
				eType:   TypeDuplicate,
				message: "duplicate contennt",
			},
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				eType:   TypeUnauthenticated,
				message: "not authenticated",
			},
		},
		{
			name: "TypeUnauthorized",
			args: args{
				eType:   TypeUnauthorized,
				message: "not authorized",
			},
		},
		{
			name: "TypeEmpty",
			args: args{
				eType:   TypeEmpty,
				message: "empty content",
			},
		},
		{
			name: "TypeNotFound",
			args: args{
				eType:   TypeNotFound,
				message: "resource not found",
			},
		},
		{
			name: "TypeMaximumAttempts",
			args: args{
				eType:   TypeMaximumAttempts,
				message: "exceeded maximum number of allowed attempts",
			},
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				eType:   TypeSubscriptionExpired,
				message: "subscription expired",
			},
		},
		{
			name: "TypeDownstreamDependencyTimedout",
			args: args{
				eType:   TypeDownstreamDependencyTimedout,
				message: "downstream dependency call timed out",
			},
		},
	}
	for _, tt := range tests {
		var want, got *Error
		switch tt.args.eType {
		case TypeInternal:
			{
				got = InternalErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeInternal)
			}
		case TypeValidation:
			{
				got = ValidationErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeValidation)
			}
		case TypeInputBody:
			{
				got = InputBodyErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeInputBody)
			}
		case TypeDuplicate:
			{
				got = DuplicateErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeDuplicate)
			}
		case TypeUnauthenticated:
			{
				got = UnauthenticatedErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeUnauthenticated)
			}
		case TypeUnauthorized:
			{
				got = UnauthorizedErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeUnauthorized)
			}

		case TypeEmpty:
			{
				got = EmptyErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeEmpty)
			}

		case TypeNotFound:
			{
				got = NotFoundErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeNotFound)
			}

		case TypeMaximumAttempts:
			{
				got = MaximumAttemptsErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeMaximumAttempts)
			}
		case TypeSubscriptionExpired:
			{
				got = SubscriptionExpiredErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeSubscriptionExpired)
			}
		case TypeDownstreamDependencyTimedout:
			{
				got = DownstreamDependencyTimedoutErr(originalErr, tt.args.message)
				_, file, line, _ := runtime.Caller(0)
				line--
				want = newerr(originalErr, tt.args.message, file, line, TypeDownstreamDependencyTimedout)
			}

		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Validation() = %v, want %v", got, want)
		}
	}
}

func TestNewWithErrMsgType(t *testing.T) {
	type args struct {
		message string
		err     error
		eType   errType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TypeInternal",
			args: args{
				err:     errors.New("Go builtin internal error"),
				eType:   TypeInternal,
				message: "internal error",
			},
		},
		{
			name: "TypeValidation",
			args: args{
				eType:   TypeValidation,
				message: "validation error",
			},
		},
		{
			name: "TypeInputBody",
			args: args{
				eType:   TypeInputBody,
				message: "invalid input body",
			},
		},
		{
			name: "TypeDuplicate",
			args: args{
				eType:   TypeDuplicate,
				message: "duplicate contennt",
			},
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				eType:   TypeUnauthenticated,
				message: "not authenticated",
			},
		},
		{
			name: "TypeUnauthorized",
			args: args{
				eType:   TypeUnauthorized,
				message: "not authorized",
			},
		},
		{
			name: "TypeEmpty",
			args: args{
				eType:   TypeEmpty,
				message: "empty content",
			},
		},
		{
			name: "TypeNotFound",
			args: args{
				eType:   TypeNotFound,
				message: "resource not found",
			},
		},
		{
			name: "TypeMaximumAttempts",
			args: args{
				eType:   TypeMaximumAttempts,
				message: "exceeded maximum number of allowed attempts",
			},
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				eType:   TypeSubscriptionExpired,
				message: "subscription expired",
			},
		},
		{
			name: "TypeDownstreamDependencyTimedout",
			args: args{
				eType:   TypeDownstreamDependencyTimedout,
				message: "downstream dependency call timed out",
			},
		},
	}
	for _, tt := range tests {
		got := NewWithErrMsgType(tt.args.err, tt.args.message, tt.args.eType)
		_, file, line, _ := runtime.Caller(0)
		line--
		want := newerr(tt.args.err, tt.args.message, file, line, tt.args.eType)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewWithErrMsgType() = %v, want %v", got, want)
		}
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
