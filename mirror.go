package errors

import "errors"

// Unwrap calls the Go builtin errors.UnUnwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is calls the Go builtin errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As calls the Go builtin errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Join returns an error that combines all the given errors.
// This is the exact implementation found in Go v1.20.
// It will be removed when Go >= v1.20 becomes the LTS version, and would just call the
// native Join after that.
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
