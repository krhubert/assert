package assert

import (
	"go/token"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// compareExported returns an [cmp.Option] that compares all exported fields of a struct,
func compareExported() cmp.Option {
	return cmp.Exporter(func(reflect.Type) bool { return true })
}

// ignoreUnexported returns an [cmp.Option] that only ignores the immediate unexported
// fields of a struct, including anonymous fields of unexported types.
func ignoreUnexported() cmp.Option {
	return cmp.FilterPath(
		func(p cmp.Path) bool {
			sf, ok := p.Index(-1).(cmp.StructField)
			if !ok {
				return false
			}

			return !token.IsExported(sf.Name())
		},
		cmp.Ignore(),
	)
}

// ignoreEmptyFields returns an [cmp.Option]
// ignores fields that are empty in the expected value.
func ignoreEmptyFields() cmp.Option {
	return cmp.FilterPath(
		func(p cmp.Path) bool {
			sf, ok := p.Index(-1).(cmp.StructField)
			if !ok {
				return false
			}

			_, wantv := sf.Values()
			return isEmptyValue(wantv)
		},
		cmp.Ignore(),
	)
}

// ignoreZeroFields returns an [cmp.Option] that
// ignores fields that have a zero value.
func ignoreZeroFields() cmp.Option {
	return cmp.FilterPath(
		func(p cmp.Path) bool {
			sf, ok := p.Index(-1).(cmp.StructField)
			if !ok {
				return false
			}

			_, wantv := sf.Values()
			return isZeroValue(wantv)
		},
		cmp.Ignore(),
	)
}

// ignoreFieldNames returns an [cmp.Option] that ignores fields of the
// given names on a single struct type.
//
// It respects the names of exported fields that are forwarded due to struct embedding.
// The struct type is specified by passing in a value of that type.
//
// The name may be a dot-delimited string (e.g., "Foo.Bar") to ignore a
// specific sub-field that is embedded or nested within the parent struct.
func ignoreFieldNames(typ any, names ...string) cmp.Option {
	sf := newStructFilter(typ, names...)
	return cmp.FilterPath(sf.filter, cmp.Ignore())
}
