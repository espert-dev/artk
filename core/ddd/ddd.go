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
type AggregateRoot[ID comparable, SerializationType any] interface {
	// ID uniquely identifies the entity.
	ID() ID

	// Serialize returns a serialization for persistence implementations.
	Serialize() SerializationType
}

// Serialization is the serialized representation of an aggregate root.
//
// A Serialization should contain all the data necessary to persist an
// aggregate. They are expected to only have exported fields and to
// contain no business logic. They are commonly used by persistence
// implementations to deserialize persisted data back into entities.
type Serialization[AggregateRoot any] interface {
	Deserialize() AggregateRoot
}
