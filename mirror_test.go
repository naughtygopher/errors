package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestUnwrap(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "wrapped",
			args: args{err: NewWithErrMsgType(
				New("level 0"),
				"level 1",
				TypeInternal,
			)},
			wantErr: true,
		},
		{
			name: "no wrapping",
			args: args{err: NewWithErrMsgType(
				nil,
				"level 1",
				TypeInternal,
			)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIs(t *testing.T) {
	target := New("level 0")
	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has err in the chain",
			args: args{
				err: NewWithErrMsgType(
					NewWithErrMsgType(
						target,
						"level 1",
						TypeInputBody,
					),
					"level 2",
					TypeNotFound,
				),
				target: target,
			},
			want: true,
		},
		{
			name: "no matching err in the chain",
			args: args{
				err: NewWithErrMsgType(
					NewWithErrMsgType(
						New("level 0"),
						"level 1",
						TypeInputBody,
					),
					"level 2",
					TypeNotFound,
				),
				target: target,
			},
			want: false,
		},
		{
			name: "external target error",
			args: args{
				err: NewWithErrMsgType(
					nil,
					"level 2",
					TypeNotFound,
				),
				target: errors.New("external"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAs(t *testing.T) {
	// ref: https://github.com/golang/go/issues/37625#issuecomment-594045710
	err := fmt.Errorf("fmt error")
	target := &Error{}
	if errors.As(err, &target) {
		t.Error("As() = true, want false")
	}

	err = New("type *Error")
	if !errors.As(err, &target) {
		t.Error("As() = false, want true")
	}

	out := target.Error()
	want := "/errors/mirror_test.go:121: type *Error"
	if !strings.Contains(out, want) {
		t.Errorf("Error() = %s, want %s", out, want)
	}
}

func TestJoin(t *testing.T) {
	joined := Join(
		errors.New("[1] std"),
		New("[2] custom"),
		errors.New("[3] std"),
		Validation("[4] validation"),
	)
	got := fmt.Sprint(joined)
	if !strings.Contains(got, "[1] std") ||
		!strings.Contains(got, "mirror_test.go:136: [2] custom") ||
		!strings.Contains(got, "[3] std") ||
		!strings.Contains(got, "mirror_test.go:138: [4] validation") {
		t.Error(got)
	}

	msg, ok := Message(joined)
	if ok {
		t.Errorf(
			"Expected: false, got: %v",
			ok,
		)
	}
	expectedMsg := `[2] custom
[4] validation`
	if msg != expectedMsg {
		t.Error(msg)
	}

	expectedCode := http.StatusUnprocessableEntity
	code, msg, ok := HTTPStatusCodeMessage(joined)
	if ok {
		t.Errorf("expected false, got: %v", ok)
	}
	if expectedMsg != msg ||
		expectedCode != code {
		t.Errorf(
			"Expected msg: %s, expected code: %d, got msg: %s, got code: %d",
			expectedMsg, expectedCode, msg, code,
		)
	}
}
