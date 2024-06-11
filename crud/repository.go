package crud

import (
	"artk.dev/ddd"
	"context"
)

// Repository defines the basic interface of a CRUD repository.
type Repository[
	A ddd.AggregateRoot[I, S],
	I comparable,
	S ddd.Serialization[A],
] interface {
	Getter[A, I, S]
	Inserter[A, I, S]
	Updater[A, I, S]
	Upserter[A, I, S]
	Deleter[A, I, S]
}

// Getter abstracts getting an entity by ID from a repository.
type Getter[
	A ddd.AggregateRoot[I, S],
	I comparable,
	S ddd.Serialization[A],
] interface {
	// Get returns the entity with the specified ID.
	// If none is found, it returns an apperror.NotFound error.
	Get(ctx context.Context, id I) (A, error)
}

// Inserter abstracts inserting a new entity into a repository.
type Inserter[
	A ddd.AggregateRoot[I, S],
	I comparable,
	S ddd.Serialization[A],
] interface {
	// Insert a new entity into the repository.
	// If there is already a value with the same ID,
	// it returns an apperror.Conflict error.
	Insert(ctx context.Context, x A) error
}

// Updater abstracts updating an entity already present in a repository.
type Updater[
	A ddd.AggregateRoot[I, S],
	I comparable,
	S ddd.Serialization[A],
] interface {
	// Update an entity already present in the repository.
	// If none is found, it returns an apperror.NotFound error.
	Update(ctx context.Context, id I, update func(x A) error) error
}

// Upserter abstracts inserting a value (if it is not already present in the
// repository) or updating it (if it is already present).
type Upserter[
	A ddd.AggregateRoot[I, S],
	I comparable,
	S ddd.Serialization[A],
] interface {
	// Upsert does one of two things:
	//   - If the entity with the specified ID does not exist,
	//     it inserts a new entity using the provided insert function.
	//   - If the entity with the specified ID exists,
	//     it updates the entity using the provided update function.
	Upsert(ctx context.Context,
		id I,
		insert func() (A, error),
		update func(x A) error,
	) error
}

// Deleter abstracts deleting an entity from the repository.
type Deleter[
	A ddd.AggregateRoot[I, S],
	I comparable,
	S ddd.Serialization[A],
] interface {
	// Delete the entity with the given ID from the repository.
	// If none is found, it returns an apperror.NotFound error.
	Delete(ctx context.Context, id I) error
}
