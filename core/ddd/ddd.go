// Package ddd provides markers for structs that indicate their DDD role.
//
// Some of these types use type traits to restrict or warn against dangerous
// behaviours, while others are purely informative.
package ddd

import (
	"artk.dev/core/typetraits"
)

// Entity is embedded into structs that represent DDD entities.
type Entity struct {
	_ typetraits.NoCopy
	_ typetraits.NoCompare
}

// Value is embedded into structs that represent DDD value objects.
type Value struct{}

// AggregateRoot is embedded into structs that represent DDD aggregate roots.
type AggregateRoot struct{}
