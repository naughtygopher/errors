<p align="center"><img src="https://user-images.githubusercontent.com/1092882/87815217-d864a680-c882-11ea-9c94-24b67f7125fe.png" alt="errors gopher" width="256px"/></p>

[![](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/bnkamalesh/errors)
[![Maintainability](https://api.codeclimate.com/v1/badges/errors/maintainability)](https://codeclimate.com/github/bnkamalesh/errors/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/bnkamalesh/errors)](https://goreportcard.com/report/github.com/bnkamalesh/errors)
[![Build Status](https://travis-ci.org/bnkamalesh/errors.svg?branch=master)](https://travis-ci.org/bnkamalesh/errors)
[![codecov](https://codecov.io/gh/bnkamalesh/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/bnkamalesh/errors)

# Errors

Errors package is a drop-in replacement of the built-in Go errors package with no external dependencies. It lets you create errors of 11 different types which should handle most of the use cases. Some of them are a bit too specific for web applications, but useful nonetheless. Following are the primary features of this package:

1. Custom error types
2. User friendly message
3. File & line number prefixed to errors

In case of nested errors, the messages (in case of nesting with this package's error) & errors are also looped through.

### Custom error types

1. TypeInternal - is the error type for when there is an internal system error. e.g. Database errors
2. TypeValidation - is the error type for when there is a validation error. e.g. invalid email address
3. TypeInputBody - is the error type for when the input data is invalid. e.g. invalid JSON
4. TypeDuplicate - is the error type for when there's duplicate content. e.g. user with email already exists (when trying to register a new user)
5. TypeUnauthenticated - is the error type when trying to access an authenticated API without authentication
6. TypeUnauthorized - is the error type for when there's an unauthorized access attempt
7. TypeEmpty - is the error type for when an expected non-empty resource, is empty
8. TypeNotFound - is the error type for an expected resource is not found. e.g. user ID not found
9. TypeMaximumAttempts - is the error type for attempting the same action more than an allowed threshold
10. TypeSubscriptionExpired - is the error type for when a user's 'paid' account has expired
11. TypeDownstreamDependencyTimedout - is the error type for when a request to a downstream dependent service times out

### User friendly messages

More often than not, when writing APIs, we'd want to respond with an easier to undersand user friendly message. Instead of returning the raw error. And log the raw error.

There are helper functions for all the error types, when in need of setting a friendly message, there are helper functions have a suffix 'Err'. All such helper functions accept the original error and a string.

```golang
package main
import(
    "fmt"
    "github.com/bnkamalesh/errors"
)

func Bar() error {
    return fmt.Errorf("hello %s", "world!")
}

func Foo() error {
    err := Bar()
    if err != nil {
        return errors.InternalErr(err, "bar is not happy")
    }
    return nil
}

func main() {
    err := Foo()
    fmt.Println(err)
    _,msg,_ := errors.HTTPStatusCodeMessage(err)
    fmt.Println(msg)
}
```

### File & line number prefixed to errors

A common annoyance with Go errors which most people are aware of is, figuring out the origin of the error, especially when there are nested function calls. Ever since error annotation was introduced in Go, a lot of people have tried using it to trace out an errors origin by giving function names, contextual message etc in it. e.g. `fmt.Errorf("database query returned error %w", err)`. This errors package, whenever you call the Go error interface's `Error() string` function, it'll print the error prefixed by the filepath and line number. It'd look like `../Users/JohnDoe/apps/main.go:50 hello world` where 'hello world' is the error message.

## How to use?

A sample was already shown in the user friendly message section, following one would show 1-2 scenarios.

```golang
import (
    "fmt"
    "log"

    "github.com/bnkamalesh/webgo/v4"
    "github.com/bnkamalesh/errors"
)

func bar() error {
    return fmt.Errorf("%s %s", "sinking", "bar")
}

func bar2() error {
    err := bar()
    if err != nil {
        return errors.InternalErr(err, "bar2 was deceived by bar1 :(")
    }
    return nil
}

func foo() error {
    err := bar2()
    if err != nil {
        return errors.InternalErr(err, "we lost bar2!")
    }
    return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    err := foo()
    if err != nil {
        // log the error on your server for troubleshooting
        log.Println(err.Error())
        // respond to request with friendly msg
        status, msg, _ := errors.HTTPStatusCodeMessage(err)
        webgo.SendError(w, msg, status)
        return
    }

    webgo.R200(w, "yay!")
}

func routes() []*webgo.Route {
	return []*webgo.Route{
		&webo.Route{
			Name: "home",
			Method: http.http.MethodGet,
			Pattern: "/",
			Handlers: []http.HandlerFunc{
				handler,
			},
		},
	}
}

func main() {
	router := webgo.NewRouter(*webgo.Config{
		Host:         "",
		Port:         "8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}, routes())

	router.UseOnSpecialHandlers(middleware.AccessLog)
	router.Use(middleware.AccessLog)
	router.Start()
}
```

[webgo](https://github.com/bnkamalesh/webgo) was used to illustrate a very specific function, `errors.HTTPStatusCodeMessage`. This just returns the appropriate http status code, user friendly message stored within and a boolean value. Boolean value is `true` if the returned error is of this package's error type.
Since we get the status code and message separately, when using any web framework, you can set values according to the respective framework's native functions. In case of Webgo, it wraps errors in a struct of its own. 

Otherwise, you could directly respond to the HTTP request by calling `WriteHTTP(error,http.ResponseWriter)`.

## Contributing

If more error types, customization etc. are required, PRs & issues are welcome!

## The gopher

The gopher used here was created using [Gopherize.me](https://gopherize.me/). Show some love to Go errors like our gopher lady here!