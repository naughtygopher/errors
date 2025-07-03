package errors

import (
	"google.golang.org/grpc/codes"
)

func grpcStatusCode(eT errType) codes.Code {
	status := codes.Unknown
	switch eT {
	case TypeInternal:
		{
			status = codes.Internal
		}
	case TypeValidation, TypeInputBody:
		{
			status = codes.InvalidArgument
		}
	case TypeDuplicate:
		{
			status = codes.AlreadyExists
		}
	case TypeUnauthenticated:
		{
			status = codes.Unauthenticated
		}
	case TypeUnauthorized:
		{
			status = codes.PermissionDenied
		}
	case TypeEmpty, TypeNotFound:
		{
			status = codes.NotFound
		}
	case TypeMaximumAttempts:
		{
			status = codes.ResourceExhausted
		}
	case TypeSubscriptionExpired:
		{
			status = codes.Unavailable
		}
	case TypeNotImplemented:
		{
			status = codes.Unimplemented
		}
	case TypeContextTimedout, TypeDownstreamDependencyTimedout:
		{
			status = codes.DeadlineExceeded
		}
	case TypeContextCancelled:
		{
			status = codes.Canceled
		}

	}

	return status
}

// GRPCStatusCodeMessage returns the appropriate GRPC status code, message, boolean for the error
// the boolean value is true if the error was of type *Error, false otherwise.
func GRPCStatusCodeMessage(err error) (codes.Code, string, bool) {
	code, isErr := GRPCStatusCode(err)
	msg, isErrMsg := Message(err)
	if msg == "" {
		msg = err.Error()
	}
	return code, msg, isErr && isErrMsg
}

// GRPCStatusCode returns appropriate GRPC response status code based on type of the error. The boolean
// is 'true' if the provided error is of type *Err. If joined error, boolean is true if all joined errors
// are of type *Error
// In case of joined errors, it'll return the status code of the last *Error
func GRPCStatusCode(err error) (codes.Code, bool) {
	derr, _ := err.(*Error)
	if derr != nil {
		return grpcStatusCode(derr.Type()), true
	}

	// Since TypeInternal is the default returned by getErrType, it is ignored.
	if et := getErrType(err); et != TypeInternal {
		return grpcStatusCode(et), false
	}

	jerr, _ := err.(*joinError)
	if jerr != nil {
		elen := len(jerr.errs)
		isErr := true
		for i := elen - 1; i >= 0; i-- {
			code, isE := GRPCStatusCode(jerr.errs[i])
			isErr = isE && isErr
			if isE {
				return code, isErr
			}
		}
	}

	return codes.Unknown, false
}
