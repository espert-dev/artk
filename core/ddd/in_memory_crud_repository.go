package ddd

import (
	"artk.dev/core/apperror"
	"artk.dev/core/clone"
	"context"
	"sync"
)

// InMemoryCrudRepository provides a generic CRUD repository implementation
// for any aggregate root. Mainly meant to be used in in-memory databases for
// tests and prototyping.
//
// Implements CrudRepository.
type InMemoryCrudRepository[
	A AggregateRoot[I, S],
	I comparable,
	S Serialization[A],
] struct {
	Mutex          sync.RWMutex
	Serializations map[I]S
}

// Reset (re-)initializes the repository.
// It must be called before other methods.
func (r *InMemoryCrudRepository[A, I, S]) Reset() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Serializations = make(map[I]S)
}

// Get returns the entity with the specified ID.
// If none is found, it returns an apperror.NotFound error.
//
// Provides Getter.
func (r *InMemoryCrudRepository[A, I, S]) Get(
	_ context.Context,
	id I,
) (A, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	serialization, ok := r.Serializations[id]
	if !ok {
		var zero A
		return zero, apperror.NotFound("not found: %v", id)
	}

	serialization = clone.Of(serialization)
	item := serialization.Deserialize()
	return item, nil
}

// Insert a new entity into the repository.
// If there is already a value with the same ID,
// it returns an apperror.Conflict error.
//
// Provides Inserter.
func (r *InMemoryCrudRepository[A, I, S]) Insert(
	_ context.Context,
	item A,
) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	id := item.ID()
	if _, ok := r.Serializations[id]; ok {
		return apperror.Conflict("already exists: %v", id)
	}

	serialization := item.Serialize()
	serialization = clone.Of(serialization)
	r.Serializations[id] = serialization
	return nil
}

// Update an entity already present in the repository.
// If none is found, it returns an apperror.NotFound error.
//
// Provides Updater.
func (r *InMemoryCrudRepository[A, I, S]) Update(
	_ context.Context,
	id I,
	update func(x A) error,
) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	serialization, ok := r.Serializations[id]
	if !ok {
		return apperror.NotFound("not found: %v", id)
	}

	serialization = clone.Of(serialization)
	item := serialization.Deserialize()
	if err := update(item); err != nil {
		return err
	}

	serialization = item.Serialize()
	serialization = clone.Of(serialization)
	r.Serializations[id] = serialization
	return nil
}

// Upsert does one of two things:
//   - If the entity with the specified ID does not exist,
//     it inserts a new entity using the provided insert function.
//   - If the entity with the specified ID exists,
//     it updates the entity using the provided update function.
//
// Provides Upserter.
func (r *InMemoryCrudRepository[A, I, S]) Upsert(
	_ context.Context,
	id I,
	insert func() (A, error),
	update func(x A) error,
) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	var err error
	var item A
	if serialization, ok := r.Serializations[id]; ok {
		c := clone.Of(serialization)
		item = c.Deserialize()
		err = update(item)
	} else {
		item, err = insert()
	}
	if err != nil {
		return err
	}

	serialization := item.Serialize()
	serialization = clone.Of(serialization)
	r.Serializations[id] = serialization
	return nil
}

// Delete the entity with the given ID from the repository.
// If none is found, it returns an apperror.NotFound error.
//
// Provides Deleter.
func (r *InMemoryCrudRepository[A, I, S]) Delete(
	_ context.Context,
	id I,
) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if _, ok := r.Serializations[id]; !ok {
		return apperror.NotFound("not found: %v", id)
	}

	delete(r.Serializations, id)
	return nil
}