[![PkgGoDev](https://pkg.go.dev/badge/github.com/krhubert/assert)](https://pkg.go.dev/github.com/krhubert/assert)
# Assert

A simple assertion package with a minimal API, yet powerful enough to cover most use cases.

## Usage

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "io/fs"
    "testing"

    "github.com/krhubert/assert"
)

func TestXYZ(t *testing.T) {
    // assert checks if value implements Equal() bool method
    // and uses it to compare values. This makes it possible
    // to compare e.g. time.Time values.
    now := time.Now()
    assert.Equal(t, now, now)

    // assert is type safe, so this will not compile
    // assert.Equal(t, 1, "1")

    // assert do not compare pointers, but values they point to
    a, b := 1, 1
    assert.Equal(t, &a, &b)

    // assert checks for errors
    assert.Error(t, errors.New("error")) 
    assert.NoError(t, nil)

    // assert checks if error contains a target value
    // like string, error or struct
    err := fmt.Errorf(
        "closed socket: %w %w",
        io.EOF,
        &fs.PathError{Op: "read", Path: "socket", Err: io.ErrClosedPipe},
    )
    assert.ErrorContains(t, err, "closed socket")
    assert.ErrorContains(t, err, io.EOF)
    var pathError *fs.PathError
    assert.ErrorContains(atb, err, &pathError)

    // assert checks if function panics or not
    assert.Panic(t, func() { panic(0) })
    assert.NotPanic(t, func() { })

    // assert can be used to check errors in defer functions
    defer assert.Defer(t, func() error { return file.Close() })()

    // assert can be used to check if a value is of a specific type
    gs := TypeAssert[fmt.GoStringer](t, time.Time{}) 
}
```
