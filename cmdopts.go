package assert

import (
	"go/token"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

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

// ignoreEmptyFields returns an [cmp.Option] that only ignores the empty values.
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

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	}

	// v.Equal(zero)
	if method := v.MethodByName("Equal"); method.IsValid() {
		methodType := method.Type()
		if methodType.NumIn() == 1 && methodType.NumOut() == 1 && methodType.Out(0).Kind() == reflect.Bool {
			zero := reflect.Zero(v.Type())
			return method.Call([]reflect.Value{zero})[0].Bool()
		}
	}

	return false
}
