[![Build and Test](https://github.com/42z-io/httpie/actions/workflows/build_test.yml/badge.svg)](https://github.com/42z-io/httpie/actions/workflows/build_test.yml)
[![Coverage Status](https://coveralls.io/repos/github/42z-io/httpie/badge.svg?branch=master)](https://coveralls.io/github/42z-io/httpie?branch=master)
[![Docs](https://img.shields.io/badge/API-docs?label=docs&color=blue&link=https%3A%2F%2Fpkg.go.dev%2Fgithub.com%2F42z-io%2Fhttpie)](https://pkg.go.dev/github.com/42z-io/httpie)
[![License](https://img.shields.io/badge/License-MIT-blue)](https://github.com/42z-io/confik/blob/main/LICENSE) [![Version](https://img.shields.io/github/tag/42z-io/httpie?include_prereleases=&sort=semver&color=blue)](https://github.com/42z-io/httpie/releases/)

Opinionated middleware, and helper functions for HTTP based applications.

Why httpie? So you can have your pie and eat it too...?

![Logo](logo.png)


# Middleware

## TransactionalMiddleware

The transaction middleware will embed a transaction (sql.Tx) into your context. 

You must provide a function which provides the sql.Tx to the middleware:

```go
middleware := httpie.TransactionalMiddleware(func(ctx context.Context) (*sql.Tx, error) {
  return db.BeginTx(ctx, nil)
})
```

You can then access the transaction from the context:

```go
func getTx(ctx context.Context) *sql.Tx {
  tx := context.Value(httpie.TransactionCtxKey)
  if tx != nil {
    return tx.(*sql.Tx)
  }
  return nil
}
```

**Note:** If a transaction is not present then your repository / service layer should either acquire one itself, or not use a transaction and rely on your normal `DB.Query` style calls.

The transaction will only be created when the HTTP request is: PUT, POST, DELETE

The transaction will be rolled back if the HTTP status is >= 400

The transaction will be automatically comitted if the HTTP status is < 400

## Logging Middleware

The logging middleware will use slog to record requests and responses.

You need to provide it with a `slog.Logger` and any configuration. The default logs a lot of common values.

```go
middleware := httpie.LoggingMiddleware(slog.Default(), httpie.LoggingOpts{
  LogRequest: true,
  LogResponse: true,
  OnResponse: httpie.DefaultLogResponse,
  OnRequest: httpie.DefaultLogRequest,
})
```

You can customize the response and request logging by providing your own OnResponse and OnRequest handlers.

# Helpers

There are various other helpers for reading/writing JSON and handling errors.

```go
WriteErr(w http.ResponseWriter, err error) error
WriteOk(w http.ResponseWriter, data T) error
WriteOkOrErr(w http.ResponseWriter, data T, err error)
ReadJson(r *http.Request, data *T) error
GetQueryParamIntDefault(r *http.Request, key string, defaultValue int) (int, error)
GetQueryParamListDefault(r *http.Request, key string, defaultValue []string) ([]string, error)
GetQueryParamDefault(r *http.Request, key string, defaultValue string) (string, error)
```

# Validation

There is a slightly modified version of `github.com/go-playground/validator/v10` that has a secure password validator (`securepassword`)

You can use this validator as follows:

```go
type MyStruct struct {
  password string `validate:"required,securepassword"`
}
```

You can execute the validator by running:

```
err := httpie.Validate(MyStruct{password: "test"})
```

`err` will be an `ErrHttpValidation` if validation fails. See [Validation Error](#validation-errors)


# Errors

There is a non standard error system in place that is useful for mapping common errors that may occur in the repository or service layer.

These errors are below and will map to standard HTTP errors when using `WriteErr`. Any other `error` passed to `WriteErr` will trigger a 500 internal service error.

| Error Name | Status Code | Message | Purpose |
| ---------- | ----------- | ------- | ------- |
| ErrBadRequest | 400 | Bad Request | Used to indicate the request was malformed |
| ErrUnauthorized | 401 | Unauthorized | Used to indicate the request has missing or invalid authorization |
| ErrForbidden | 403 | Forbidden | The user is authenticated but not authorized for the resource |
| ErrNotFound | 404 | Not Found | The resource was not found |
| ErrConflict | 409 | Conflict | The resource already exists |
| ErrInternal | 500 | Internal Server Error | There was an unexepected error |

These errors are not meant to be comprehensive, it is useful to have errors that may occur in the service layer (like not finding an object) be able to propagate with the correct http error codes.

You can implement new error codes like this:

```go
var ErrMyError = httpie.NewErrHttp(status_code, "error_message")
```

The errors when rendered using `WriteErr` will be in JSON format:

```json
{
  "message": "not found"
}
```

## Validation Errors

There is a special variant of `ErrHttp` called `ErrHttpValidation`. This includes some extra information for returning a `map[string]string` of errors.

These errors are meant for an API which understands the failures.

You can use it as follows:

```go
err := NewHttpErrValidation(map[string]string{"field":"error_code"})
```

When passed to `WriteErr` it will return a 400 bad request with the following JSON output:

```json
{
  "message": "validation failed",
  "errors": {
    "field": "error_code"
  }
}
```

# Watched Response Writer

In middleware you often want to be able to look at the response, and optionally override it before actually writing it to the client.

There is a `WatchedResponseWriter` that is a simple wrapper around `http.ResponseWriter`

It will delay actually writing any requests to the response until `Apply()` is called. It can also be `Reset()` if the middleware determines it want's to send something else.

This is used by the `TransactionalMiddleware` to ensure we send an internal server error if a `tx.Commit()` fails.

It is also used by the `LoggingMiddleware` to capture the HTTP status code.

**Note:** This naively uses a buffer to capture the written bytes, it's likely not a problem but for something high performance this could be an issue [just a theory]

You can use it in middleware like this:

```go
// Middleware that will convert any http status code >=400 into a 500 internal server error
func MyMiddlewareHandler(w http.ResponseWriter, r *http.Request) {
  // Create a watched response writer
  ww := httpie.NewWatchedResponseWriter(w)
  
  // We need to apply the response when we are ready
  defer ww.Apply()
  
  // Call the http handler
  next.ServeHTTP(ww, r)

  // Detect some error and do something different
  if ww.StatusCode() >= 400 {
    // Reset the response buffer and status code
    ww.Reset()

    // Send a totally different message
    ww.WriteHeader(500)
    ww.Write("internal server error")
    // OR httpie.WriteErr(ww, ErrInternal)
  }
}
```

# Context

The middleware in this package tends to inject things into the context. You often need to be able to pull things out of the context.

## GetContextValue

You can get a typed object out of the context as follows:

```go
type myKeyType int
var uniqueCtxKey myKeyType = 0

type MyStruct struct {
  MyValue int
}

myStruct := MyStruct {
  MyValue: 42,
}

// Assign to the context
ctx := context.WithValue(context.Background(), uniqueCtxKey, &myStruct)

// Get from the context
ctx = httpie.GetContextValue[MyStruct](ctx, ctxKey)
```

The value will be `nil` if it did not exist in the context or if the type was incorrect.

# Clock Service

The clock service is a simple wrapper around `time.Now().UTC()`. It's purpose is to allow fine-grained mocking in your
service layer by including the clock service as a dependency.

You can make use of `ClockServiceMock` to mock time in your service layer during testing.

To use the clock service:

```go
type MyService struct {
  clockService IClockService
}

func (s *MyService) GetNow() time.Time {
  return s.clockService.Now()
}

func NewMyService(clockService IClockService) *MyService {
  return &MyService{
    clockService,
  }
}

cs := new(ClockService)
myService := NewMyService(clockService)
myService.GetNow()
```