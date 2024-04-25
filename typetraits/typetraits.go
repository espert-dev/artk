// Package typetraits provides traits that can be used to restrict behaviour.
package typetraits

// NoCopy can be embedded into a struct to protect against copying.
// Any copies will be flagged by the copylocks analysis in go vet.
//
// Example:
//
//	type T struct {
//		typetraits.NoCompare
//	}
//
//	x := T{} // Vetting error!
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

// NoCompare can be embedded into a struct to prevent comparison.
//
// Example:
//
//	type T struct {
//		typetraits.NoCompare
//	}
//
//	T{} := T{} // Compilation error!
type NoCompare struct {
	_ [0]func()
}
