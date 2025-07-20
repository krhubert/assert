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

func Test_equal(t *testing.T) {
	type T struct {
		t  *time.Time
		ts []*time.Time
	}

	True(t, equal[*int](nil, nil))
	False(t, equal[any](nil, 0))
	True(t, equal(0, 0))
	False(t, equal(1, 0))
	True(t, equal([]byte("hello"), []byte("hello")))
	False(t, equal([]byte("hello"), []byte("-")))
	//
	// interface Equal(V) bool
	now := time.Now()
	now1 := now.Add(1)
	utc := now.In(time.UTC)

	wlc := Must(time.LoadLocation("Europe/Warsaw"))
	waw := now.In(wlc)

	True(t, equal(utc, waw))
	True(t, equal(&utc, &waw))
	True(t, equal(
		[]*T{{t: &utc}},
		[]*T{{t: &waw}},
	))
	True(t, equal(
		&T{t: &utc},
		&T{t: &waw},
	))
	True(t, equal(
		&T{t: &utc, ts: []*time.Time{&utc}},
		&T{t: &waw, ts: []*time.Time{&waw}},
	))
	True(t, equal(
		[]T{{t: &utc, ts: []*time.Time{&utc}}},
		[]T{{t: &waw, ts: []*time.Time{&waw}}},
	))
	False(t, equal(
		[]T{{t: &utc, ts: []*time.Time{&utc}}},
		[]T{{t: &waw, ts: []*time.Time{&now1}}},
	))
	False(t, equal(now, time.Time{}))

	// dereference pointers
	a, b := 0, 0
	True(t, equal(&a, &b))
	True(t, equal([]*int{&a}, []*int{&b}))
}

func TestEqual(t *testing.T) {
	atb := &assertTB{TB: t}
	Equal(atb, 0, 0)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Equal(atb, time.Now(), time.Time{})
	atb.fail(t, "expected equal")

	atb = &assertTB{TB: t}
	Equal(atb, bytes.NewReader([]byte("a")), bytes.NewReader(nil))
	atb.fail(t, "expected equal")

	Panic(t, func() {
		Equal(t, fmt.Errorf("0"), fmt.Errorf("0"))
	})
}

func TestEqualUnexported(t *testing.T) {
	type T struct {
		A int
		b int
	}

	atb := &assertTB{TB: t}
	got := T{A: 1, b: 2}
	want := T{A: 1, b: 3}
	Equal(atb, got, want, IgnoreUnexported())
	atb.pass(t)
}

func TestEqualSkipEmptyFields(t *testing.T) {
	type T struct {
		A int
		b int
		C time.Time
		D []int
	}

	atb := &assertTB{TB: t}
	got := T{A: 1, b: 2, C: time.Now(), D: []int{1}}
	want := T{b: 2}
	Equal(atb, got, want, SkipEmptyFields())
	atb.pass(t)
}

func TestNotEqual(t *testing.T) {
	atb := &assertTB{TB: t}
	NotEqual(atb, 0, 1)
	atb.pass(t)

	atb = &assertTB{TB: t}
	NotEqual(atb, 0, 0)
	atb.fail(t, "expected not equal, but got equal")

	Panic(t, func() {
		NotEqual(t, fmt.Errorf("0"), fmt.Errorf("0"))
	})
}

func TestError(t *testing.T) {
	atb := &assertTB{TB: t}
	Error(atb, fmt.Errorf("0"))
	atb.pass(t)

	atb = &assertTB{TB: t}
	Error(atb, nil)
	atb.fail(t, "expected error, got nil")
}

func TestNoError(t *testing.T) {
	atb := &assertTB{TB: t}
	NoError(atb, nil)
	atb.pass(t)

	atb = &assertTB{TB: t}
	NoError(atb, fmt.Errorf("0"))
	atb.fail(t, "unexpected error: 0")
}

func TestErrorContains(t *testing.T) {
	err := fmt.Errorf(
		"closed socket: %w %w",
		io.EOF,
		&fs.PathError{Op: "read", Path: "socket", Err: io.ErrClosedPipe},
	)

	tests := []struct {
		name   string
		err    error
		target any
		fail   string
	}{
		{
			name:   "string match",
			err:    err,
			target: "closed socket:",
		},
		{
			name:   "string regex match",
			err:    err,
			target: "closed socket: .*",
		},
		{
			name:   "error match",
			err:    err,
			target: io.EOF,
		},
		{
			name:   "error wrap",
			err:    err,
			target: io.ErrClosedPipe,
		},
		{
			name: "fs.PathError match",
			err:  err,
			target: func() **fs.PathError {
				var pathError *fs.PathError
				return &pathError
			}(),
		},
		{
			name:   "nil error",
			err:    nil,
			target: "",
			fail:   "error is nil",
		},
		{
			name:   "string not found",
			err:    err,
			target: "open socket",
			fail:   "unexpected error:",
		},
		{
			name:   "error not found",
			err:    err,
			target: io.ErrNoProgress,
			fail:   "unexpected error:",
		},
		{
			name: "json.SyntaxError not found",
			err:  err,
			target: func() **json.SyntaxError {
				var jsonError *json.SyntaxError
				return &jsonError
			}(),
			fail: "unexpected error:",
		},
		{
			name:   "string not found in error",
			err:    err,
			target: "[",
			fail:   "does not contain \"[\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atb := &assertTB{TB: t}
			ErrorContains(atb, tt.err, tt.target)
			atb.check(t, tt.fail)
		})
	}
}

func TestErrorWant(t *testing.T) {
	atb := &assertTB{TB: t}
	ErrorWant(atb, true, fmt.Errorf("0"))
	atb.pass(t)

	atb = &assertTB{TB: t}
	ErrorWant(atb, false, nil)
	atb.pass(t)

	atb = &assertTB{TB: t}
	ErrorWant(atb, true, nil)
	atb.fail(t, "expected error: got nil")

	atb = &assertTB{TB: t}
	ErrorWant(atb, false, fmt.Errorf("0"))
	atb.fail(t, "unexpected error: 0")
}

func TestNil(t *testing.T) {
	atb := &assertTB{TB: t}
	Nil(atb, nil)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Nil(atb, map[string]int(nil))
	atb.pass(t)

	atb = &assertTB{TB: t}
	Nil(atb, chan int(nil))
	atb.pass(t)

	atb = &assertTB{TB: t}
	Nil(atb, 0)
	atb.fail(t, "expected nil, got 0")

	Panic(t, func() {
		Nil(t, fmt.Errorf("0"))
	})
}

func TestNotNil(t *testing.T) {
	atb := &assertTB{TB: t}
	NotNil(atb, 0)
	atb.pass(t)

	atb = &assertTB{TB: t}
	NotNil(atb, nil)
	atb.fail(t, "expected not nil, got nil")

	Panic(t, func() {
		NotNil(t, fmt.Errorf("0"))
	})
}

func TestZero(t *testing.T) {
	tests := []struct {
		value any
		fail  string
	}{
		{value: nil, fail: ""},
		{value: time.Time{}, fail: ""},
		{value: time.Time{}.In(time.Local), fail: ""},
		{value: 0, fail: ""},
		{value: .0, fail: ""},
		{value: make(chan int), fail: "expected zero, got 0x"},
		{value: map[string]string(nil), fail: ""},
		{value: make(map[string]string), fail: "expected zero, got map[]"},
		{value: []int(nil), fail: ""},
		{value: []int{}, fail: "expected zero, got []"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.value), func(t *testing.T) {
			atb := &assertTB{TB: t}
			Zero(atb, tt.value)
			atb.check(t, tt.fail)
		})
	}
}

func TestNotZero(t *testing.T) {
	tests := []struct {
		value any
		fail  string
	}{
		{value: nil, fail: "expected not zero, got <nil>"},
		{value: time.Time{}, fail: "expected not zero, got 0001-01-01 00:00:00 +0000 UTC"},
		{value: time.Time{}.In(time.Local), fail: "expected not zero, got 0001-01"},
		{value: 0, fail: "expected not zero, got 0"},
		{value: .0, fail: "expected not zero, got 0"},
		{value: make(chan int), fail: ""},
		{value: map[string]string(nil), fail: "expected not zero, got map[]"},
		{value: make(map[string]string), fail: ""},
		{value: []int(nil), fail: "expected not zero, got []"},
		{value: []int{}, fail: ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.value), func(t *testing.T) {
			atb := &assertTB{TB: t}
			NotZero(atb, tt.value)
			atb.check(t, tt.fail)
		})
	}
}

func TestLen(t *testing.T) {
	atb := &assertTB{TB: t}
	Len(atb, []int{1, 2, 3}, 3)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Len(atb, [3]int{1, 2, 3}, 3)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Len(atb, map[int]int{1: 1}, 1)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Len(atb, "hello", 5)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Len(atb, make(chan int), 0)
	atb.pass(t)

	atb = &assertTB{TB: t}
	Len(atb, []int{1, 2, 3}, 2)
	atb.fail(t, "expected length 2, got 3")
}

func TestTrue(t *testing.T) {
	atb := &assertTB{TB: t}
	True(atb, true)
	atb.pass(t)

	atb = &assertTB{TB: t}
	True(atb, false)
	atb.fail(t, "expected true, got false")
}

func TestFalse(t *testing.T) {
	atb := &assertTB{TB: t}
	False(atb, false)
	atb.pass(t)

	atb = &assertTB{TB: t}
	False(atb, true)
	atb.fail(t, "expected false, got true")
}

func TestPanic(t *testing.T) {
	atb := &assertTB{TB: t}
	Panic(atb, func() { panic(0) })
	atb.pass(t)

	atb = &assertTB{TB: t}
	Panic(atb, func() {})
	atb.fail(t, "expected panic, got nothing")
}

func TestNotPanic(t *testing.T) {
	atb := &assertTB{TB: t}
	NotPanic(atb, func() {})
	atb.pass(t)

	atb = &assertTB{TB: t}
	NotPanic(atb, func() { panic(0) })
	atb.fail(t, "unexpected panic: 0")
}

func TestDefer(t *testing.T) {
	atb := &assertTB{TB: t}
	func() {
		defer Defer(atb, func() error { return nil })
	}()
	atb.pass(t)

	atb = &assertTB{TB: t}
	func() {
		fn := func() error { return fmt.Errorf("0") }
		defer Defer(atb, fn)()
	}()
	atb.fail(t, "unexpected defer error: 0")
}

func TestTypeAssert(t *testing.T) {
	atb := &assertTB{TB: t}
	TypeAssert[int](atb, 0)
	atb.pass(t)

	atb = &assertTB{TB: t}
	TypeAssert[io.Reader](atb, &bytes.Buffer{})
	atb.pass(t)

	atb = &assertTB{TB: t}
	TypeAssert[string](atb, 0)
	atb.fail(t, "assertion string.(int) failed")
}

func TestMust(t *testing.T) {
	Panic(t, func() {
		Must(0, fmt.Errorf("err"))
	})
}

func TestMust2(t *testing.T) {
	Panic(t, func() {
		Must2("", 0, fmt.Errorf("err"))
	})
}

func TestMust3(t *testing.T) {
	Panic(t, func() {
		Must3(true, "", 0, fmt.Errorf("err"))
	})
}

type assertTB struct {
	testing.TB

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

func (atb *assertTB) check(t testing.TB, fail string) {
	t.Helper()
	if fail != "" {
		atb.fail(t, fail)
	} else {
		atb.pass(t)
	}
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
