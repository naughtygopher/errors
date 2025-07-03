package errors

import (
	"net/http"
)

func httpStatusCode(eT errType) int {
	status := http.StatusInternalServerError
	switch eT {
	case TypeValidation:
		{
			status = http.StatusUnprocessableEntity
		}
	case TypeInputBody:
		{
			status = http.StatusBadRequest
		}
	case TypeDuplicate:
		{
			status = http.StatusConflict
		}
	case TypeUnauthenticated:
		{
			status = http.StatusUnauthorized
		}
	case TypeUnauthorized:
		{
			status = http.StatusForbidden
		}
	case TypeEmpty:
		{
			status = http.StatusGone
		}
	case TypeNotFound:
		{
			status = http.StatusNotFound
		}
	case TypeMaximumAttempts:
		{
			status = http.StatusTooManyRequests
		}
	case TypeSubscriptionExpired:
		{
			status = http.StatusPaymentRequired
		}
	case TypeNotImplemented:
		{
			status = http.StatusNotImplemented
		}
	case TypeContextTimedout, TypeContextCancelled:
		{
			status = http.StatusRequestTimeout
		}
	}

	return status
}

// HTTPStatusCodeMessage returns the appropriate HTTP status code, message, boolean for the error
// the boolean value is true if the error was of type *Error, false otherwise.
func HTTPStatusCodeMessage(err error) (int, string, bool) {
	code, isErr := HTTPStatusCode(err)
	msg, isErrMsg := Message(err)
	if msg == "" {
		msg = err.Error()
	}
	return code, msg, isErr && isErrMsg
}

// HTTPStatusCode returns appropriate HTTP response status code based on type of the error. The boolean
// is 'true' if the provided error is of type *Err. If joined error, boolean is true if all joined errors
// are of type *Error
// In case of joined errors, it'll return the status code of the last *Error
func HTTPStatusCode(err error) (int, bool) {
	derr, _ := err.(*Error)
	if derr != nil {
		return httpStatusCode(derr.Type()), true
	}

	// Since TypeInternal is the default returned by getErrType, it is ignored.
	if et := getErrType(err); et != TypeInternal {
		return httpStatusCode(et), false
	}

	jerr, _ := err.(*joinError)
	if jerr != nil {
		elen := len(jerr.errs)
		isErr := true
		for i := elen - 1; i >= 0; i-- {
			code, isE := HTTPStatusCode(jerr.errs[i])
			isErr = isE && isErr
			if isE {
				return code, isErr
			}
		}
	}

	return http.StatusInternalServerError, false
}
