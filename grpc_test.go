package errors

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestGRPCStatusCodeMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  codes.Code
		want1 string
		want2 bool
	}{
		{
			name: "TypeInternal",
			args: args{
				err: Internalf("unknown error occurred:%s", "formatted"),
			},
			want:  codes.Internal,
			want1: "unknown error occurred:formatted",
			want2: true,
		},
		{
			name: "TypeInternal - Go builtin error type",
			args: args{
				err: errors.New("unknown error occurred"),
			},
			want:  codes.Unknown,
			want1: "unknown error occurred",
			want2: false,
		},
		{
			name: "TypeValidation",
			args: args{
				err: Validationf("invalid email provided:%s", "formatted"),
			},
			want:  codes.InvalidArgument,
			want1: "invalid email provided:formatted",
			want2: true,
		},
		{
			name: "TypeInputBody",
			args: args{
				err: InputBodyf("invalid json provided:%s", "formatted"),
			},
			want:  codes.InvalidArgument,
			want1: "invalid json provided:formatted",
			want2: true,
		},
		{
			name: "TypeDuplicate",
			args: args{
				err: Duplicatef("duplicate content detected:%s", "formatted"),
			},
			want:  codes.AlreadyExists,
			want1: "duplicate content detected:formatted",
			want2: true,
		},
		{
			name: "TypeUnauthenticated",
			args: args{
				err: Unauthenticatedf("authentication required:%s", "formatted"),
			},
			want:  codes.Unauthenticated,
			want1: "authentication required:formatted",
			want2: true,
		},
		{
			name: "TypeUnauthorized",
			args: args{
				err: Unauthorizedf("not authorized to access this resource:%s", "formatted"),
			},
			want:  codes.PermissionDenied,
			want1: "not authorized to access this resource:formatted",
			want2: true,
		},
		{
			name: "TypeEmpty",
			args: args{
				err: Emptyf("empty content not expected:%s", "formatted"),
			},
			want:  codes.NotFound,
			want1: "empty content not expected:formatted",
			want2: true,
		},
		{
			name: "TypeNotFound",
			args: args{
				err: NotFoundf("requested resource not found:%s", "formatted"),
			},
			want:  codes.NotFound,
			want1: "requested resource not found:formatted",
			want2: true,
		},
		{
			name: "TypeMaximumAttempts",
			args: args{
				err: MaximumAttemptsf("exceeded maximum number of requests allowed:%s", "formatted"),
			},
			want:  codes.ResourceExhausted,
			want1: "exceeded maximum number of requests allowed:formatted",
			want2: true,
		},
		{
			name: "TypeSubscriptionExpired",
			args: args{
				err: SubscriptionExpiredf("your subscription has expired:%s", "formatted"),
			},
			want:  codes.Unavailable,
			want1: "your subscription has expired:formatted",
			want2: true,
		},
		{
			name: "TypeDownstreamDependencyTimedout",
			args: args{
				err: DownstreamDependencyTimedoutf("dependency timed out:%s", "formatted"),
			},
			want:  codes.DeadlineExceeded,
			want1: "dependency timed out:formatted",
			want2: true,
		},
		{
			name: "TypeNotImplemented",
			args: args{
				err: NotImplementedf("feature not implemented:%s", "formatted"),
			},
			want:  codes.Unimplemented,
			want1: "feature not implemented:formatted",
			want2: true,
		},
		{
			name: "TypeContextCancelled",
			args: args{
				err: ContextCancelledf("context cancelled:%s", "formatted"),
			},
			want:  codes.Canceled,
			want1: "context cancelled:formatted",
			want2: true,
		},
		{
			name: "TypeContextTimedout",
			args: args{
				err: ContextTimedoutf("context timed out:%s", "formatted"),
			},
			want:  codes.DeadlineExceeded,
			want1: "context timed out:formatted",
			want2: true,
		},
		{
			name: "joined: TypeContextTimedout",
			args: args{
				err: Join(errors.New("base error"), ContextTimedout("context timed out")),
			},
			want:  codes.DeadlineExceeded,
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
			want:  codes.Canceled,
			want1: "context timed out\ncontext cancelled",
			want2: false,
		},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := GRPCStatusCodeMessage(tt.args.err)
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
