package assert

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type equaler struct {
	// unexported ignores unexported fields of structs.
	unexported bool

	// skipEmptyFields ignores struct fields that are empty.
	skipEmptyFields bool
}

// EqualOption configures the equality check behavior.
type EqualOption func(o *equaler)

func (o *equaler) apply(opts ...EqualOption) cmp.Options {
	for _, opt := range opts {
		opt(o)
	}

	out := []cmp.Option{}
	if o.unexported {
		out = append(out, ignoreUnexported())
	} else {
		out = append(out, compareExported())
	}

	if o.skipEmptyFields {
		out = append(out, ignoreEmptyFields())
	}

	return out
}

// IgnoreUnexported returns an EqualOption that ignores unexported fields of structs.
func IgnoreUnexported() EqualOption {
	return func(o *equaler) {
		o.unexported = true
	}
}

// SkipEmptyFields returns an EqualOption that ignores struct fields that are empty.
func SkipEmptyFields() EqualOption {
	return func(o *equaler) {
		o.skipEmptyFields = true
	}
}

// Equal checks if two values are equal with the given options.
//
// This functions uses [go-cmp](https://pkg.go.dev/github.com/google/go-cmp) to determine equality.
func Equal[V any](t testing.TB, got V, want V, opts ...EqualOption) {
	if _, ok := any(got).(error); ok {
		panic("use assert.Error() for errors")
	}

	t.Helper()
	if !equal(got, want, opts...) {
		t.Fatalf("expected equal\n%s", diffValue(got, want, opts...))
	}
}

// NotEqual checks if two values are not equal.
// See [Equal] for rules used to determine equality.
func NotEqual[T any](t testing.TB, got T, want T, opts ...EqualOption) {
	if _, ok := any(got).(error); ok {
		panic("use assert.Error() for errors")
	}

	t.Helper()
	if equal(got, want, opts...) {
		t.Fatalf("expected not equal, but got equal")
	}
}

// Error checks if an error is not nil.
func Error(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// NoError checks if an error is nil.
func NoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ErrorContains checks if an error is not nil and contains the target.
//
// Target can be:
//
// 1. string
//
// The string is compiled as a regexp, and the error is matched against it.
// If it is not a valid regexp, it is used as a string to check if the error contains it.
//
// 2. error
//
// The error is checked if it is equal to the target using errors.Is.
//
// 3. type
//
// The error is checked if it can be converted to the target type using errors.As.
func ErrorContains(t testing.TB, err error, target any) {
	t.Helper()
	if err == nil {
		t.Fatalf("error is nil")
		return
	}

	// catch any errors.Is/As panics
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("error.Is/As panic %s", r)
		}
	}()

	switch e := target.(type) {
	case string:
		// if this is a valid regexp, compile it and use it
		// otherwise, just use it as a string

		// first check the string itself
		if !strings.Contains(err.Error(), e) {
			if re, err1 := regexp.Compile(e); err1 == nil {
				if !re.MatchString(err.Error()) {
					t.Fatalf("unexpected error: %q does not match %q", err, e)
				}
			} else {
				t.Fatalf("unexpected error: %q does not contain %q", err, e)
			}
		}

	case error:
		if !errors.Is(err, e) {
			t.Fatalf("unexpected error: %q is not %T", err, e)
		}

	default:
		if !errors.As(err, e) {
			t.Fatalf("unexpected error: %q is not %T", err, e)
		}
	}
}

// ErrorWant checks if an error is expected for the test.
// A common usage in tests is:
//
//	type tests struct {
//		name    string
//		// other fields
//		wantErr bool
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			err := fn()
//			assert.ErrorWant(t, tt.wantErr, err)
//		})
//	}
func ErrorWant(t testing.TB, want bool, err error) {
	t.Helper()
	if want && err == nil {
		t.Fatalf("expected error: got nil")
	} else if !want && err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Zero checks if got is zero value.
func Zero[T comparable](t testing.TB, got T) {
	t.Helper()
	if got != *new(T) {
		t.Fatalf("expected zero, got %v", got)
	}
}

// NotZero checks if got is not zero value.
func NotZero[V comparable](t testing.TB, got V) {
	t.Helper()
	if got == *new(V) {
		t.Fatalf("expected not zero, got %v", got)
	}
}

// Nil checks if got is nil.
func Nil(t testing.TB, got any) {
	if _, ok := got.(error); ok {
		panic("use assert.NoError() for errors")
	}

	t.Helper()
	if !isNil(got) {
		t.Fatalf("expected nil, got %v", got)
	}
}

// NotNil checks if got is not nil.
func NotNil(t testing.TB, got any) {
	if _, ok := got.(error); ok {
		panic("use assert.Error() for errors")
	}

	t.Helper()
	if isNil(got) {
		t.Fatalf("expected not nil, got nil")
	}
}

// Len checks if the length of got is l.
// got can be any go type accepted by builtin len function.
func Len[V any](t testing.TB, got V, want int) {
	t.Helper()

	l := reflect.ValueOf(got).Len()
	if l != want {
		t.Fatalf("expected length %d, got %d", want, l)
	}
}

// True checks if got is true.
func True(t testing.TB, got bool) {
	t.Helper()
	if !got {
		t.Fatalf("expected true, got false")
	}
}

// False checks if got is false.
func False(t testing.TB, got bool) {
	t.Helper()
	if got {
		t.Fatalf("expected false, got true")
	}
}

// Panic checks if f panics.
func Panic(t testing.TB, f func()) {
	t.Helper()

	defer func() {
		t.Helper()
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got nothing")
		}
	}()
	f()
}

// NotPanic checks if f does not panic.
func NotPanic(t testing.TB, f func()) {
	t.Helper()

	defer func() {
		t.Helper()
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	f()
}

// Defer returns a function that will call fn and check if an error is returned.
func Defer(t testing.TB, fn func() error) func() {
	t.Helper()
	return func() {
		if err := fn(); err != nil {
			t.Fatalf("unexpected defer error: %v", err)
		}
	}
}

// TypeAssert checks if got is of type V and returns it.
func TypeAssert[V any](t testing.TB, got any) V {
	t.Helper()
	v, ok := got.(V)
	if !ok {
		t.Fatalf("assertion %T.(%T) failed", v, got)
	}
	return v
}

// Must is a helper function to handle a single return value from a function.
func Must[P1 any](p1 P1, err error) P1 {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return p1
}

// Must2 is a helper function to handle two return values from a function.
func Must2[P1 any, P2 any](p1 P1, p2 P2, err error) (P1, P2) {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return p1, p2
}

// Must3 is a helper function to handle three return values from a function.
func Must3[P1 any, P2 any, P3 any](p1 P1, p2 P2, p3 P3, err error) (P1, P2, P3) {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return p1, p2, p3
}

func equal[V any](got V, want V, opts ...EqualOption) bool {
	eq := equaler{}
	cmpOpts := eq.apply(opts...)
	return cmp.Equal(got, want, cmpOpts...)
}

func isNil(obj any) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Pointer, reflect.UnsafePointer, reflect.Interface,
		reflect.Slice:
		return v.IsNil()
	}
	return false
}

func diffValue[V any](a V, b V, opts ...EqualOption) string {
	// first let GoStringer format the values if they implement it
	out := ""
	if _, ok := any(a).(fmt.GoStringer); ok {
		out += diffGoStringer(any(a).(fmt.GoStringer), any(b).(fmt.GoStringer))
		out += "\n"
	}

	eq := equaler{}
	cmpOpts := eq.apply(opts...)
	out += "diff:\n"
	out += cmp.Diff(a, b, cmpOpts...)
	return out
}

func diffGoStringer(a, b fmt.GoStringer) string {
	got := "nil"
	if !isNil(a) {
		got = a.GoString()
	}

	want := "nil"
	if !isNil(b) {
		want = b.GoString()
	}
	return fmt.Sprintf(" got: %s\nwant: %s\n", got, want)
}
