package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func newerr(e error, message string, etype errType, skip int) *Error {
	pcs := make([]uintptr, 128)
	_ = runtime.Callers(skip, pcs)
	return &Error{
		original: e,
		message:  message,
		eType:    etype,
		pcs:      pcs,
		pc:       pcs[0] - 1,
	}
}

func newerrf(e error, etype errType, skip int, format string, args ...any) *Error {
	message := fmt.Sprintf(format, args...)
	return newerr(e, message, etype, skip)
}

func getErrType(err error) errType {
	e, _ := err.(*Error)
	if e != nil {
		return e.Type()
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return TypeContextTimedout
	}

	if errors.Is(err, context.Canceled) {
		return TypeContextCancelled
	}

	return TypeInternal
}

// Wrap is used to simply wrap an error with optional message; error type would be the
// default error type set using SetDefaultType; TypeInternal otherwise
// If the error being wrapped is already of type Error, then its respective type is used
func Wrap(original error, msg ...string) *Error {
	message := strings.Join(msg, ". ")
	return newerr(original, message, getErrType(original), 3)
}

func Wrapf(original error, format string, args ...any) *Error {
	return newerrf(original, getErrType(original), 4, format, args...)
}

// Deprecated: WrapWithMsg [deprecated, use `Wrap`] wrap error with a user friendly message
func WrapWithMsg(original error, msg string) *Error {
	return newerr(original, msg, getErrType(original), 3)
}

// NewWithType returns an error instance with custom error type
func NewWithType(msg string, etype errType) *Error {
	return newerr(nil, msg, etype, 3)
}

// NewWithTypef returns an error instance with custom error type. And formatted message
func NewWithTypef(etype errType, format string, args ...any) *Error {
	return newerrf(nil, etype, 4, format, args...)
}

// NewWithErrMsgType returns an error instance with custom error type and message
func NewWithErrMsgType(original error, message string, etype errType) *Error {
	return newerr(original, message, etype, 3)
}

// NewWithErrMsgTypef returns an error instance with custom error type and formatted message
func NewWithErrMsgTypef(original error, etype errType, format string, args ...any) *Error {
	return newerrf(original, etype, 4, format, args...)
}

// Internal helper method for creating internal errors
func Internal(message string) *Error {
	return newerr(nil, message, TypeInternal, 3)
}

// Internalf helper method for creating internal errors with formatted message
func Internalf(format string, args ...any) *Error {
	return newerrf(nil, TypeInternal, 4, format, args...)
}

// Validation is a helper function to create a new error of type TypeValidation
func Validation(message string) *Error {
	return newerr(nil, message, TypeValidation, 3)
}

// Validationf is a helper function to create a new error of type TypeValidation, with formatted message
func Validationf(format string, args ...any) *Error {
	return newerrf(nil, TypeValidation, 4, format, args...)
}

// InputBody is a helper function to create a new error of type TypeInputBody
func InputBody(message string) *Error {
	return newerr(nil, message, TypeInputBody, 3)
}

// InputBodyf is a helper function to create a new error of type TypeInputBody, with formatted message
func InputBodyf(format string, args ...any) *Error {
	return newerrf(nil, TypeInputBody, 4, format, args...)
}

// Duplicate is a helper function to create a new error of type TypeDuplicate
func Duplicate(message string) *Error {
	return newerr(nil, message, TypeDuplicate, 3)
}

// Duplicatef is a helper function to create a new error of type TypeDuplicate, with formatted message
func Duplicatef(format string, args ...any) *Error {
	return newerrf(nil, TypeDuplicate, 4, format, args...)
}

// Unauthenticated is a helper function to create a new error of type TypeUnauthenticated
func Unauthenticated(message string) *Error {
	return newerr(nil, message, TypeUnauthenticated, 3)
}

// Unauthenticatedf is a helper function to create a new error of type TypeUnauthenticated, with formatted message
func Unauthenticatedf(format string, args ...any) *Error {
	return newerrf(nil, TypeUnauthenticated, 4, format, args...)

}

// Unauthorized is a helper function to create a new error of type TypeUnauthorized
func Unauthorized(message string) *Error {
	return newerr(nil, message, TypeUnauthorized, 3)
}

// Unauthorizedf is a helper function to create a new error of type TypeUnauthorized, with formatted message
func Unauthorizedf(format string, args ...any) *Error {
	return newerrf(nil, TypeUnauthorized, 4, format, args...)
}

// Empty is a helper function to create a new error of type TypeEmpty
func Empty(message string) *Error {
	return newerr(nil, message, TypeEmpty, 3)
}

// Emptyf is a helper function to create a new error of type TypeEmpty, with formatted message
func Emptyf(format string, args ...any) *Error {
	return newerrf(nil, TypeEmpty, 4, format, args...)
}

// NotFound is a helper function to create a new error of type TypeNotFound
func NotFound(message string) *Error {
	return newerr(nil, message, TypeNotFound, 3)
}

// NotFoundf is a helper function to create a new error of type TypeNotFound, with formatted message
func NotFoundf(format string, args ...any) *Error {
	return newerrf(nil, TypeNotFound, 4, format, args...)
}

// MaximumAttempts is a helper function to create a new error of type TypeMaximumAttempts
func MaximumAttempts(message string) *Error {
	return newerr(nil, message, TypeMaximumAttempts, 3)
}

// MaximumAttemptsf is a helper function to create a new error of type TypeMaximumAttempts, with formatted message
func MaximumAttemptsf(format string, args ...any) *Error {
	return newerrf(nil, TypeMaximumAttempts, 4, format, args...)
}

// SubscriptionExpired is a helper function to create a new error of type TypeSubscriptionExpired
func SubscriptionExpired(message string) *Error {
	return newerr(nil, message, TypeSubscriptionExpired, 3)
}

// SubscriptionExpiredf is a helper function to create a new error of type TypeSubscriptionExpired, with formatted message
func SubscriptionExpiredf(format string, args ...any) *Error {
	return newerrf(nil, TypeSubscriptionExpired, 4, format, args...)
}

// DownstreamDependencyTimedout is a helper function to create a new error of type TypeDownstreamDependencyTimedout
func DownstreamDependencyTimedout(message string) *Error {
	return newerr(nil, message, TypeDownstreamDependencyTimedout, 3)
}

// DownstreamDependencyTimedoutf is a helper function to create a new error of type TypeDownstreamDependencyTimedout, with formatted message
func DownstreamDependencyTimedoutf(format string, args ...any) *Error {
	return newerrf(nil, TypeDownstreamDependencyTimedout, 4, format, args...)
}

// NotImplemented is a helper function to create a new error of type TypeNotImplemented
func NotImplemented(message string) *Error {
	return newerr(nil, message, TypeNotImplemented, 3)
}

// NotImplementedf is a helper function to create a new error of type TypeNotImplemented, with formatted message
func NotImplementedf(format string, args ...any) *Error {
	return newerrf(nil, TypeNotImplemented, 4, format, args...)
}

// ContextCancelled is a helper function to create a new error of type TypeContextCancelled
func ContextCancelled(message string) *Error {
	return newerr(nil, message, TypeContextCancelled, 3)
}

// ContextCancelledf is a helper function to create a new error of type TypeContextCancelled, with formatted message
func ContextCancelledf(format string, args ...any) *Error {
	return newerrf(nil, TypeContextCancelled, 4, format, args...)
}

// TypeContextTimedout is a helper function to create a new error of type TypeContextTimedout
func ContextTimedout(message string) *Error {
	return newerr(nil, message, TypeContextTimedout, 3)
}

// ContextTimedoutf is a helper function to create a new error of type TypeContextTimedout, with formatted message
func ContextTimedoutf(format string, args ...any) *Error {
	return newerrf(nil, TypeContextTimedout, 4, format, args...)
}

// InternalErr helper method for creation internal errors which also accepts an original error
func InternalErr(original error, message string) *Error {
	return newerr(original, message, TypeInternal, 3)
}

// InternalErr helper method for creation internal errors which also accepts an original error, with formatted message
func InternalErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeInternal, 4, format, args...)
}

// ValidationErr helper method for creation validation errors which also accepts an original error
func ValidationErr(original error, message string) *Error {
	return newerr(original, message, TypeValidation, 3)
}

// ValidationErr helper method for creation validation errors which also accepts an original error, with formatted message
func ValidationErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeValidation, 4, format, args...)
}

// InputBodyErr is a helper function to create a new error of type TypeInputBody which also accepts an original error
func InputBodyErr(original error, message string) *Error {
	return newerr(original, message, TypeInputBody, 3)
}

// InputBodyErrf is a helper function to create a new error of type TypeInputBody which also accepts an original error, with formatted message
func InputBodyErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeInputBody, 4, format, args...)
}

// DuplicateErr is a helper function to create a new error of type TypeDuplicate which also accepts an original error
func DuplicateErr(original error, message string) *Error {
	return newerr(original, message, TypeDuplicate, 3)
}

// DuplicateErrf is a helper function to create a new error of type TypeDuplicate which also accepts an original error, with formatted message
func DuplicateErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeDuplicate, 4, format, args...)
}

// UnauthenticatedErr is a helper function to create a new error of type TypeUnauthenticated which also accepts an original error
func UnauthenticatedErr(original error, message string) *Error {
	return newerr(original, message, TypeUnauthenticated, 3)
}

// UnauthenticatedErrf is a helper function to create a new error of type TypeUnauthenticated which also accepts an original error, with formatted message
func UnauthenticatedErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeUnauthenticated, 4, format, args...)
}

// UnauthorizedErr is a helper function to create a new error of type TypeUnauthorized which also accepts an original error
func UnauthorizedErr(original error, message string) *Error {
	return newerr(original, message, TypeUnauthorized, 3)
}

// UnauthorizedErrf is a helper function to create a new error of type TypeUnauthorized which also accepts an original error, with formatted message
func UnauthorizedErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeUnauthorized, 4, format, args...)
}

// EmptyErr is a helper function to create a new error of type TypeEmpty which also accepts an original error
func EmptyErr(original error, message string) *Error {
	return newerr(original, message, TypeEmpty, 3)
}

// EmptyErr is a helper function to create a new error of type TypeEmpty which also accepts an original error, with formatted message
func EmptyErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeEmpty, 4, format, args...)
}

// NotFoundErr is a helper function to create a new error of type TypeNotFound which also accepts an original error
func NotFoundErr(original error, message string) *Error {
	return newerr(original, message, TypeNotFound, 3)
}

// NotFoundErrf is a helper function to create a new error of type TypeNotFound which also accepts an original error, with formatted message
func NotFoundErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeNotFound, 4, format, args...)
}

// MaximumAttemptsErr is a helper function to create a new error of type TypeMaximumAttempts which also accepts an original error
func MaximumAttemptsErr(original error, message string) *Error {
	return newerr(original, message, TypeMaximumAttempts, 3)
}

// MaximumAttemptsErr is a helper function to create a new error of type TypeMaximumAttempts which also accepts an original error, with formatted message
func MaximumAttemptsErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeMaximumAttempts, 4, format, args...)
}

// SubscriptionExpiredErr is a helper function to create a new error of type TypeSubscriptionExpired which also accepts an original error
func SubscriptionExpiredErr(original error, message string) *Error {
	return newerr(original, message, TypeSubscriptionExpired, 3)
}

// SubscriptionExpiredErrf is a helper function to create a new error of type TypeSubscriptionExpired which also accepts an original error, with formatted message
func SubscriptionExpiredErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeSubscriptionExpired, 4, format, args...)
}

// DownstreamDependencyTimedoutErr is a helper function to create a new error of type TypeDownstreamDependencyTimedout which also accepts an original error
func DownstreamDependencyTimedoutErr(original error, message string) *Error {
	return newerr(original, message, TypeDownstreamDependencyTimedout, 3)
}

// DownstreamDependencyTimedoutErrf is a helper function to create a new error of type TypeDownstreamDependencyTimedout which also accepts an original error, with formatted message
func DownstreamDependencyTimedoutErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeDownstreamDependencyTimedout, 4, format, args...)
}

// NotImplementedErr is a helper function to create a new error of type TypeNotImplemented which also accepts an original error
func NotImplementedErr(original error, message string) *Error {
	return newerr(original, message, TypeNotImplemented, 3)
}

// NotImplementedErrf is a helper function to create a new error of type TypeNotImplemented which also accepts an original error, with formatted message
func NotImplementedErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeNotImplemented, 4, format, args...)
}

// ContextCancelledErr is a helper function to create a new error of type TypeContextCancelled which also accepts an original error
func ContextCancelledErr(original error, message string) *Error {
	return newerr(original, message, TypeContextCancelled, 3)
}

// ContextCancelledErrf is a helper function to create a new error of type TypeContextCancelled which also accepts an original error, with formatted message
func ContextCancelledErrf(original error, format string, args ...any) *Error {
	return newerrf(original, TypeContextCancelled, 4, format, args...)
}

// ErrWithoutTrace is a duplicate of Message, but with clearer name. The boolean is 'true' if the
// provided err is of type *Error
func ErrWithoutTrace(err error) (string, bool) {
	return Message(err)
}

// Message recursively concatenates all the messages set while creating/wrapping the errors. The boolean
// is 'true' if the provided error is of type *Err
func Message(err error) (string, bool) {
	derr, _ := err.(*Error)
	if derr != nil {
		return derr.Message(), true
	}

	jerr, _ := err.(*joinError)
	if jerr != nil {
		list := make([]string, 0, len(jerr.errs))
		isErr := true
		for i := range jerr.errs {
			msg, ok := Message(jerr.errs[i])
			isErr = isErr && ok
			if msg == "" {
				continue
			}
			list = append(list, msg)
		}
		return strings.Join(list, "\n"), isErr
	}

	return "", false
}

// WriteHTTP is a convenience method which will check if the error is of type *Error and
// respond appropriately
func WriteHTTP(err error, w http.ResponseWriter) {
	// INFO: consider sending back "unknown server error" message
	status, msg, _ := HTTPStatusCodeMessage(err)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg))
}

// Type returns the errType if it's an instance of *Error, -1 otherwise
// In case of joined error, it'll return the type of the last *Error
func Type(err error) errType {
	e, _ := err.(*Error)
	if e != nil {
		return e.Type()
	}
	je, _ := err.(*joinError)
	if je != nil {
		for i := len(je.errs) - 1; i >= 0; i-- {
			et := Type(je.errs[i])
			if et.Int() != -1 {
				return et
			}
		}
	}

	return errType(-1)
}

// Type returns the errType as integer if it's an instance of *Error, -1 otherwise
func TypeInt(err error) int {
	return Type(err).Int()
}

// HasType will check if the provided err type is available anywhere nested in the error
func HasType(err error, et errType) bool {
	if err == nil {
		return false
	}

	e, _ := err.(*Error)
	if e == nil {
		return HasType(errors.Unwrap(err), et)
	}

	je, _ := err.(*joinError)
	if je != nil {
		for i := 0; i < len(je.errs); i++ {
			if HasType(je.errs[i], et) {
				return true
			}
		}
	}

	if e.Type() == et {
		return true
	}

	return HasType(e.Unwrap(), et)
}
