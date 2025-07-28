# Assert

[![PkgGoDev](https://pkg.go.dev/badge/github.com/krhubert/assert)](https://pkg.go.dev/github.com/krhubert/assert)

An assertion package with a minimal API, yet powerful enough to cover most use cases.

The `assert` package is designed to provide a more configurable and type-safe alternative to the popular `testify/assert` package.

## Motivation

Have you ever struggled with comparing complex structs in your tests? One that use `decimal.Decimal`, `time.Time`, or `uuid.UUID` types? Or perhaps you wanted to ignore automatically generated fields, such as `ID` or `CreatedAt`, in your struct comparisons?

`assert` allows for more flexible assertions, such as ignoring unexported fields or skipping empty fields in struct comparisons.

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

## Suite

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "io/fs"
    "testing"

 _ "github.com/mattn/go-sqlite3"
 "go.uber.org/mock/gomock"
    "github.com/krhubert/assert"
)

type ExampleTestSuite struct {
    assert.Suite

    i    int
    ctrl *gomock.Controller
    db  *sql.DB
}

func (s *ExampleTestSuite) Setup(t *testing.T) {
    s.i = 4
    s.ctrl = gomock.NewController(t)
    db, err := sql.Open("sqlite3", ":memory:")
    assert.NoError(t, err)
    s.db = db
}

func (s *ExampleTestSuite) Teardown(t *testing.T) {
    s.i = 0
    assert.NoError(t, s.db.Close())
}

func TestExample(t *testing.T) {
     ts := assert.Setup[ExampleTestSuite](t)
     row := ts.db.QueryRow("select date()")
     var date string
     assert.NoError(ts, row.Scan(&date))
     assert.NoError(ts, row.Err())
}
```

## Assert vs testify

```go
package assert

import (
 "testing"
 "time"

 "github.com/google/uuid"
 "github.com/shopspring/decimal"
 "github.com/stretchr/testify/require"
)

func TestExampleEqual(t *testing.T) {
  type User struct {
   Id        uuid.UUID
   Email     string
   CreatedAt time.Time
   Balance   decimal.Decimal

   active bool
  }

  loc, _ := time.LoadLocation("Europe/Warsaw")

  // db.CreateUser("test@example.com")
  createdAt := time.Now()
  user := User{
   Id:        uuid.New(),
   Email:     "test@example.com",
   CreatedAt: createdAt,
   Balance:   decimal.NewFromFloat(1),

   active: true,
  }

  want := User{
    Email:     "test@example.com",
    CreatedAt: createdAt.In(loc),
    Balance:   decimal.RequireFromString("1"),
  }
}
```

Running `go test` on:

<table>
<thead><tr><th>Assert</th><th>Testify</th></tr></thead>
<tbody>
<tr><td>

```go
assert.Equal(t, user, want, IgnoreUnexported(), SkipEmptyFields())
```

```text
PASS: TestExampleEqual (0.00s)
```

</td><td>

```go
require.Equal(t, user, want)
```

```text
--- FAIL: TestExampleEqual (0.00s)
    test_test.go:35:
                Error Trace:
                Error:          Not equal:
                                expected: assert.User{Id:uuid.UUID{0x66, 0x43, 0x33, 0x3b, 0xad, 0xf6, 0x48, 0xec, 0x9a, 0x7d, 0xff, 0x53, 0xc0, 0x90, 0x6e, 0xf1}, Email:"test@example.com", CreatedAt:time.Date(2025, time.July, 17, 13, 12, 17, 81207156, time.Local), Balance:decimal.Decimal{value:(*big.Int)(0xc0000b5b00), exp:0}, active:true}
                                actual  : assert.User{Id:uuid.UUID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, Email:"test@example.com", CreatedAt:time.Date(2025, time.July, 17, 13, 12, 17, 81207156, time.Location("Europe/Warsaw")), Balance:decimal.Decimal{value:(*big.Int)(0xc0000b5b20), exp:0}, active:false}

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -2,3 +2,3 @@
                                  Id: (uuid.UUID) (len=16) {
                                -  00000000  66 43 33 3b ad f6 48 ec  9a 7d ff 53 c0 90 6e f1  |fC3;..H..}.S..n.|
                                +  00000000  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
                                  },
                                @@ -6,6 +6,6 @@
                                  CreatedAt: (time.Time) {
                                -  wall: (uint64) 13985458619913019252,
                                -  ext: (int64) 1236572,
                                +  wall: (uint64) 81207156,
                                +  ext: (int64) 63888347537,
                                   loc: (*time.Location)({
                                -   name: (string) (len=5) "Local",
                                +   name: (string) (len=13) "Europe/Warsaw",
                                    zone: ([]time.zone) (len=11) {
                                @@ -1078,3 +1078,3 @@
                                  },
                                - active: (bool) true
                                + active: (bool) false
                                 }
                Test:           TestExampleEqual
```

</td></tr>
</tbody></table>
