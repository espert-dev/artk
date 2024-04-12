// Package typetraits provides traits that can be used to restrict behaviour.
package typetraits

// NoCopy can be embedded into a struct to warn against copy.
type NoCopy struct{}

func (n *NoCopy) Lock() {
	// Do nothing.
}

func (n *NoCopy) Unlock() {
	// Do nothing.
}

// NoCompare can be embedded into a struct to prevent comparison.
type NoCompare struct {
	_ [0]func()
}
