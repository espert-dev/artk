package typetraits_test

import (
	"artk.dev/typetraits"
)

func ExampleNoCompare() {
	type T struct {
		_ typetraits.NoCompare
	}

	x := T{}
	_ = x

	// Uncomment to try:
	// x == T{} // Error!
}

func ExampleNoCopy() {
	type T struct {
		_ typetraits.NoCopy
	}

	// Uncomment to try:
	// T{} := T{}  // Error!

	// This is just here to keep the compiler happy.
	_ = T{}

	// In general, you want to use an anonymous field to avoid exposing
	// the dummy methods in typetraits.NoCopy. Otherwise, the below code
	// becomes possible, which is harmless but confusing.
	type S struct {
		typetraits.NoCopy // Embedding instead of anonymous field!
	}

	z := new(S)
	z.Lock()         // Generally, this method should not be exposed.
	defer z.Unlock() // Generally, this method should not be exposed.
}
