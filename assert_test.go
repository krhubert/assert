package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"
	"time"
)

func TestAreEqual(t *testing.T) {
	True(t, areEqual[*int](nil, nil))
	False(t, areEqual[any](nil, 0))
	True(t, areEqual(0, 0))
	False(t, areEqual(1, 0))
	True(t, areEqual([]byte("hello"), []byte("hello")))
	False(t, areEqual([]byte("hello"), []byte("-")))

	// interface Equal(V) bool
	now := time.Now()
	True(t, areEqual(now, now))
	True(t, areEqual(&now, &now))
	False(t, areEqual(now, time.Time{}))

	// dereference pointers
	a, b := 0, 0
	True(t, areEqual(&a, &b))
	True(t, areEqual([]*int{&a}, []*int{&b}))
}

func TestEqual(t *testing.T) {
	atb := new(assertTB)
	Equal(atb, 0, 0)
	atb.pass(t)

	atb = new(assertTB)
	Equal(atb, time.Now(), time.Time{})
	atb.fail(t, "expected equal")

	atb = new(assertTB)
	Equal(atb, bytes.NewReader([]byte("a")), bytes.NewReader(nil))
	atb.fail(t, "expected equal")

	Panic(t, func() {
		Equal(t, fmt.Errorf("0"), fmt.Errorf("0"))
	})
}

func TestNotEqual(t *testing.T) {
	atb := new(assertTB)
	NotEqual(atb, 0, 1)
	atb.pass(t)

	atb = new(assertTB)
	NotEqual(atb, 0, 0)
	atb.fail(t, "expected not equal, but got equal")

	Panic(t, func() {
		NotEqual(t, fmt.Errorf("0"), fmt.Errorf("0"))
	})
}

func TestError(t *testing.T) {
	atb := new(assertTB)
	Error(atb, fmt.Errorf("0"))
	atb.pass(t)

	atb = new(assertTB)
	Error(atb, nil)
	atb.fail(t, "expected error, got nil")
}

func TestNoError(t *testing.T) {
	atb := new(assertTB)
	NoError(atb, nil)
	atb.pass(t)

	atb = new(assertTB)
	NoError(atb, fmt.Errorf("0"))
	atb.fail(t, "unexpected error: 0")
}

func TestErrorContains(t *testing.T) {
	err := fmt.Errorf(
		"closed socket: %w %w",
		io.EOF,
		&fs.PathError{Op: "read", Path: "socket", Err: io.ErrClosedPipe},
	)

	atb := new(assertTB)
	ErrorContains(atb, err, "closed socket:")
	atb.pass(t)

	atb = new(assertTB)
	ErrorContains(atb, err, "closed socket: .*")
	atb.pass(t)

	atb = new(assertTB)
	ErrorContains(atb, err, io.EOF)
	atb.pass(t)

	atb = new(assertTB)
	ErrorContains(atb, err, io.ErrClosedPipe)
	atb.pass(t)

	atb = new(assertTB)
	var pathError *fs.PathError
	ErrorContains(atb, err, &pathError)
	atb.pass(t)

	atb = new(assertTB)
	ErrorContains(atb, err, "open socket")
	atb.fail(t, "unexpected error:")

	atb = new(assertTB)
	ErrorContains(atb, err, io.ErrNoProgress)
	atb.fail(t, "unexpected error:")

	atb = new(assertTB)
	var jsonError *json.SyntaxError
	ErrorContains(atb, err, &jsonError)
	atb.fail(t, "unexpected error:")

	atb = new(assertTB)
	ErrorContains(atb, nil, "")
	atb.fail(t, "error is nil")
}

func TestNil(t *testing.T) {
	atb := new(assertTB)
	Nil(atb, nil)
	atb.pass(t)

	atb = new(assertTB)
	Nil(atb, 0)
	atb.fail(t, "expected nil, got 0")

	Panic(t, func() {
		Nil(t, fmt.Errorf("0"))
	})
}

func TestNotNil(t *testing.T) {
	atb := new(assertTB)
	NotNil(atb, 0)
	atb.pass(t)

	atb = new(assertTB)
	NotNil(atb, nil)
	atb.fail(t, "expected not nil, got nil")

	Panic(t, func() {
		NotNil(t, fmt.Errorf("0"))
	})
}

func TestZero(t *testing.T) {
	atb := new(assertTB)
	Zero(atb, 0)
	atb.pass(t)

	atb = new(assertTB)
	Zero(atb, 1)
	atb.fail(t, "expected zero, got 1")
}

func TestNotZero(t *testing.T) {
	atb := new(assertTB)
	NotZero(atb, 1)
	atb.pass(t)

	atb = new(assertTB)
	NotZero(atb, 0)
	atb.fail(t, "expected not zero, got 0")
}

func TestLen(t *testing.T) {
	atb := new(assertTB)
	Len(atb, []int{1, 2, 3}, 3)
	atb.pass(t)

	atb = new(assertTB)
	Len(atb, [3]int{1, 2, 3}, 3)
	atb.pass(t)

	atb = new(assertTB)
	Len(atb, map[int]int{1: 1}, 1)
	atb.pass(t)

	atb = new(assertTB)
	Len(atb, "hello", 5)
	atb.pass(t)

	atb = new(assertTB)
	Len(atb, make(chan int), 0)
	atb.pass(t)

	atb = new(assertTB)
	Len(atb, []int{1, 2, 3}, 2)
	atb.fail(t, "expected length 2, got 3")
}

func TestTrue(t *testing.T) {
	atb := new(assertTB)
	True(atb, true)
	atb.pass(t)

	atb = new(assertTB)
	True(atb, false)
	atb.fail(t, "expected true, got false")
}

func TestFalse(t *testing.T) {
	atb := new(assertTB)
	False(atb, false)
	atb.pass(t)

	atb = new(assertTB)
	False(atb, true)
	atb.fail(t, "expected false, got true")
}

func TestPanic(t *testing.T) {
	atb := new(assertTB)
	Panic(atb, func() { panic(0) })
	atb.pass(t)

	atb = new(assertTB)
	Panic(atb, func() {})
	atb.fail(t, "expected panic, got nothing")
}

func TestNotPanic(t *testing.T) {
	atb := new(assertTB)
	NotPanic(atb, func() {})
	atb.pass(t)

	atb = new(assertTB)
	NotPanic(atb, func() { panic(0) })
	atb.fail(t, "unexpected panic: 0")
}

func TestDefer(t *testing.T) {
	atb := new(assertTB)
	func() {
		defer Defer(atb, func() error { return nil })
	}()
	atb.pass(t)

	atb = new(assertTB)
	func() {
		fn := func() error { return fmt.Errorf("0") }
		defer Defer(atb, fn)()
	}()
	atb.fail(t, "unexpected defer error: 0")
}

func TestTypeAssert(t *testing.T) {
	atb := new(assertTB)
	TypeAssert[int](atb, 0)
	atb.pass(t)

	atb = new(assertTB)
	TypeAssert[io.Reader](atb, &bytes.Buffer{})
	atb.pass(t)

	atb = new(assertTB)
	TypeAssert[string](atb, 0)
	atb.fail(t, "assertion string.(int) failed")
}

type assertTB struct {
	helper  bool
	failed  bool
	format  string
	args    []any
	message string
}

func (atb *assertTB) Helper() {
	atb.helper = true
}

func (atb *assertTB) Fatalf(format string, args ...any) {
	atb.failed = true
	atb.format = format
	atb.args = args
	atb.message = fmt.Sprintf(format, args...)
}

func (atb *assertTB) pass(t testing.TB) {
	t.Helper()
	if !atb.helper {
		t.Fatal("Helper not called")
	}
	if atb.failed {
		t.Fatalf("expected pass, got failed")
	}
}

func (atb *assertTB) fail(t testing.TB, message string) {
	t.Helper()
	if !atb.helper {
		t.Fatal("Helper not called")
	}

	if !atb.failed {
		t.Fatalf("expected failed, got pass")
	}

	if !strings.Contains(atb.message, message) {
		t.Fatalf("expected message %q, got %q", message, atb.message)
	}
}
