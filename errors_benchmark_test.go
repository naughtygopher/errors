package errors

import (
	"errors"
	"testing"
)

func Benchmark_Internal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Internal("hello world")
	}
}
func Benchmark_Internalf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Internalf("%s prefixed", "hello world")
	}
}

func Benchmark_InternalErr(b *testing.B) {
	err := errors.New("bad error")
	for i := 0; i < b.N; i++ {
		_ = InternalErr(err, "hello world")
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

func BenchmarkHasType(b *testing.B) {
	err := ValidationErr(
		DuplicateErr(
			SubscriptionExpiredErr(
				ValidationErr(
					InputBodyErr(
						Internal("hello world"),
						DefaultMessage,
					),
					DefaultMessage,
				),
				DefaultMessage,
			),
			DefaultMessage,
		),
		DefaultMessage,
	)
	for i := 0; i < b.N; i++ {
		if !HasType(err, TypeInternal) {
			b.Fatal("TypeInternal not found")
		}
	}
}
