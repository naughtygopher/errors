<p align="center"><img src="https://user-images.githubusercontent.com/1092882/87815217-d864a680-c882-11ea-9c94-24b67f7125fe.png" alt="errors gopher" width="256px"/></p>

[![](https://github.com/naughtygopher/errors/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/naughtygopher/errors/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/naughtygopher/errors.svg)](https://pkg.go.dev/github.com/naughtygopher/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/naughtygopher/errors?cache_updated=2025-07-04)](https://goreportcard.com/report/github.com/naughtygopher/errors)
[![Coverage Status](https://coveralls.io/repos/github/naughtygopher/errors/badge.svg?branch=master&cache_updated=2025-07-04)](https://coveralls.io/github/naughtygopher/errors?branch=master)
[![](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#error-handling)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/creativecreature/sturdyc/blob/master/LICENSE)

# Errors v1.3.1

Errors package is a drop-in replacement of the built-in Go errors package. It lets you create errors of 14 different types,
which should handle most of the use cases. Some of them are a bit too specific for web applications, but useful nonetheless.

Features of this package:

1. Multiple (14) error types
2. Easy handling of User friendly message(s)
3. Stacktrace - formatted, unfromatted, custom format (refer tests in errors_test.go)
4. Retrieve the Program Counters for the stacktrace
5. Retrieve runtime.Frames using `errors.RuntimeFrames(err error)` for the stacktrace
6. HTTP status code and user friendly message (wrapped messages are concatenated) for all error types
7. GRPC status code and user friendly message (wrapped messages are concatenated) for all error types
8. Helper functions to generate each error type
9. Helper function to get error Type, error type as int, check if error type is wrapped anywhere in chain
10. . `fmt.Formatter` support

In case of nested errors, the messages & errors are also looped through the full chain of errors.

### Available error types

1. TypeInternal - For internal system error. e.g. Database errors
2. TypeValidation - For validation error. e.g. invalid email address
3. TypeInputBody - For invalid input data. e.g. invalid JSON
4. TypeDuplicate - For duplicate content error. e.g. user with email already exists (when trying to register a new user)
5. TypeUnauthenticated - For not authenticated error
6. TypeUnauthorized - For unauthorized access error
7. TypeEmpty - For when an expected non-empty resource, is empty
8. TypeNotFound - For expected resource not found. e.g. user ID not found
9. TypeMaximumAttempts - For attempting the same action more than an allowed threshold
10. TypeSubscriptionExpired - For when a user's 'paid' account has expired
11. TypeDownstreamDependencyTimedout - For when a request to a downstream dependent service times out
12. TypeNotImplemented - For when the requested function cannot be fullfilled because of incapability
13. TypeContextTimedout - For when the Go context has timed out
14. TypeContextCancelled - For when the Go context has been cancelled

Helper functions are available for all the error types. Each of them have 3 helper functions, one which accepts only a string,
another which accepts an original error as well as a user friendly message, and one which accepts format string along with arguments.

All the dedicated error type functions are documented [here](https://pkg.go.dev/github.com/naughtygopher/errors?tab=doc#DownstreamDependencyTimedout).
Names are consistent with the error type, e.g. errors.Internal(string) and errors.InternalErr(error, string)

### User friendly messages

More often than not when writing APIs, we'd want to respond with an easier to undersand user friendly message.
Instead of returning the raw error, just log the raw error.

There are helper functions for all the error types. When in need of setting a friendly message for an existing error, there
are helper functions with the _suffix_ **'Err'**. All such helper functions accept the original error and a string.

```golang
package main

import (
	"fmt"

	"github.com/naughtygopher/errors"
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

	fmt.Println("err:", err)
	fmt.Println("\nerr.Error():", err.Error())

	fmt.Printf("\nformatted +v: %+v\n", err)
	fmt.Printf("\nformatted v: %v\n", err)
	fmt.Printf("\nformatted +s: %+s\n", err)
	fmt.Printf("\nformatted s: %s\n", err)

	_, msg, _ := errors.HTTPStatusCodeMessage(err)
	fmt.Println("\nmsg:", msg)
}
```

Output

```
err: bar is not happy

err.Error(): /Users/k.balakumaran/go/src/github.com/naughtygopher/errors/cmd/main.go:16: bar is not happy
hello world!bar is not happy

formatted +v: /Users/k.balakumaran/go/src/github.com/naughtygopher/errors/cmd/main.go:16: bar is not happy
hello world!bar is not happy

formatted v: bar is not happy

formatted +s: bar is not happy: hello world!

formatted s: bar is not happy

msg: bar is not happy
```

[Playground link](https://go.dev/play/p/OiLegJ9Xxc9)

### File & line number prefixed to errors

A common annoyance with Go errors which most people are aware of is, figuring out the origin of the error, especially when there are nested function calls.
Annotations help a lot by being able to provide contextual message to errors. e.g. `fmt.Errorf("database query returned error %w", err)`.
However in this package, the `Error() string` function (Go error interface method), prints the error prefixed by the filepath and line number. It'd look like `../Users/JohnDoe/apps/main.go:50 hello world` where 'hello world' is the error message.

### HTTP/GRPC status code & message

The functions `errors.HTTPStatusCodeMessage(error) (int, string, bool), errors.GRPCStatusCodeMessage(error) (int, string, bool)` returns the HTTP/GRPC status code, message, and a boolean value. The boolean is true, if the error is of type \*Error from this package. If error is nested, it unwraps and returns a single concatenated message. Sample described in the 'How to use?' section.

## How to use?

Other than the functions explained earlier in the _**User friendly messages**_ section, more examples are provided below.

```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/naughtygopher/errors"
	"github.com/naughtygopher/webgo/v6"
	"github.com/naughtygopher/webgo/v6/middleware/accesslog"
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
		fmt.Println(err.Error())
		// respond to request with friendly msg
		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		webgo.SendError(w, msg, status)
		return
	}

	webgo.R200(w, "yay!")
}

func routes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:    "home",
			Method:  http.MethodGet,
			Pattern: "/",
			Handlers: []http.HandlerFunc{
				handler,
			},
		},
	}
}

func main() {
	router := webgo.NewRouter(&webgo.Config{
		Host:         "",
		Port:         "8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}, routes()...)

	router.UseOnSpecialHandlers(accesslog.AccessLog)
	router.Use(accesslog.AccessLog)
	router.Start()
}
```

[webgo](https://github.com/naughtygopher/webgo) was used to illustrate the usage of the function, `errors.HTTPStatusCodeMessage`. It returns the appropriate http status code, user friendly message stored within, and a boolean value. Boolean value is `true` if the returned error of type \*Error.
Since we get the status code and message separately, when using any web framework, you can set values according to the respective framework's native functions. In case of Webgo, it wraps errors in a struct of its own. Otherwise, you could directly respond to the HTTP request by calling `errors.WriteHTTP(error,http.ResponseWriter)`.

Once the app is running, you can check the response by opening `http://localhost:8080` on your browser. Or on terminal

```bash
$ curl http://localhost:8080
{"errors":"we lost bar2!. bar2 was deceived by bar1 :(","status":500} // output
```

And the `fmt.Println(err.Error())` generated output on stdout would be:

```bash
/Users/username/go/src/errorscheck/main.go:28 /Users/username/go/src/errorscheck/main.go:20 sinking bar
```

## Benchmark [2025-07-03]

Macbook Air 13-inch, M3, 2024, Memory: 24 GB

```bash
$ go version
go version go1.24.4 darwin/arm64

$ go test -benchmem -bench .
goos: darwin
goarch: arm64
pkg: github.com/naughtygopher/errors
cpu: Apple M3
Benchmark_Internal-8                            	 3650916	       321.7 ns/op	    1104 B/op	       2 allocs/op
Benchmark_Internalf-8                           	 3155463	       378.9 ns/op	    1128 B/op	       3 allocs/op
Benchmark_InternalErr-8                         	 3866085	       312.2 ns/op	    1104 B/op	       2 allocs/op
Benchmark_InternalGetError-8                    	 1983544	       608.0 ns/op	    1576 B/op	       6 allocs/op
Benchmark_InternalGetErrorWithNestedError-8     	 2419369	       497.8 ns/op	    1592 B/op	       6 allocs/op
Benchmark_InternalGetMessage-8                  	 3815074	       316.1 ns/op	    1104 B/op	       2 allocs/op
Benchmark_InternalGetMessageWithNestedError-8   	 3470449	       342.2 ns/op	    1128 B/op	       3 allocs/op
Benchmark_HTTPStatusCodeMessage-8               	40540940	        29.12 ns/op	      16 B/op	       1 allocs/op
BenchmarkHasType-8                              	100000000	        11.44 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/naughtygopher/errors	13.805s
```

## Contributing

More error types, customization, features, multi-errors; PRs & issues are welcome!

## The gopher

The gopher used here was created using [Gopherize.me](https://gopherize.me/). Show some love to Go errors like our gopher lady here!
