package errors

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
		{
			name: "TypeNotImplemented",
			args: args{
				err: NotImplemented("feature not implemented"),
			},
			want:  http.StatusNotImplemented,
			want1: "feature not implemented",
			want2: true,
		},
		{
			name: "TypeContextCancelled",
			args: args{
				err: context.Canceled,
			},
			want:  http.StatusRequestTimeout,
			want1: "context canceled",
			want2: false,
		},
		{
			name: "TypeContextTimedout",
			args: args{
				err: context.DeadlineExceeded,
			},
			want:  http.StatusRequestTimeout,
			want1: "context deadline exceeded",
			want2: false,
		},
		{
			name: "joined: TypeContextTimedout",
			args: args{
				err: Join(errors.New("base error"), ContextTimedout("context timed out")),
			},
			want:  http.StatusRequestTimeout,
			want1: "context timed out",
			want2: false,
		},
		{
			name: "joined: TypeContextTimedout",
			args: args{
				err: Join(
					errors.New("base error"),
					ContextTimedout("context timed out"),
					ContextCancelled("context cancelled"),
				),
			},
			want:  http.StatusRequestTimeout,
			want1: "context timed out\ncontext cancelled",
			want2: false,
		},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := HTTPStatusCodeMessage(tt.args.err)
			if got != tt.want {
				t.Errorf(
					"[%d/%s] HTTPStatusCodeMessage() got = %v, want %v",
					idx, tt.name, got, tt.want,
				)
			}
			if got1 != tt.want1 {
				t.Errorf(
					"[%d/%s] HTTPStatusCodeMessage() got1 = %v, want %v",
					idx, tt.name, got1, tt.want1,
				)
			}
			if got2 != tt.want2 {
				t.Errorf("[%d/%s] HTTPStatusCodeMessage() got2 = %v, want %v",
					idx, tt.name, got2, tt.want2,
				)
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
