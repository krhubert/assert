package assert

import "testing"

func ExampleEqual(t *testing.T) {
	type Foo struct {
		Bar string
		bar []int
	}

	Equal(
		t,
		Foo{Bar: "Bar1", bar: []int{2, 2, 3}},
		Foo{Bar: "Bar", bar: []int{1, 2, 3}},
	)
	// Output:
}
