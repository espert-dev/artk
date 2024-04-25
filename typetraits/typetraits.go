// Package typetraits provides traits that can be used to restrict behaviour.
//
// The size of the types defined in this package is guaranteed to be zero
// and have no behaviour, which eliminates runtime overhead or memory
// misalignment concerns.
package typetraits

// NoCopy can be embedded or nested into a struct to protect against copying.
// Any copy will be flagged by the copylocks analysis in go vet.
//
// The size of this struct is guaranteed to be zero.
//
// Example:
//
//	type T struct {
//		typetraits.NoCopy
//	}
//
//	x := T{} // The initial assignment is OK.
//	y := x   // Copy: vetting error!
type NoCopy struct {
	// Use an anonymous field to prevent leaking methods.
	_ noCopy
}

type noCopy struct{}

func (n *noCopy) Lock() {
	// Do nothing.
}

func (n *noCopy) Unlock() {
	// Do nothing.
}

// NoCompare can be embedded or nested into a struct to prevent comparison.
// Any comparison will cause a compilation error.
//
// The size of this struct is guaranteed to be zero.
//
// Example:
//
//	type T struct {
//		typetraits.NoCompare
//	}
//
//	T{} == T{} // Comparison: compilation error!
type NoCompare struct {
	_ [0]func()
}
