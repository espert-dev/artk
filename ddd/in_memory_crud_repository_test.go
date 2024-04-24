package ddd_test

import (
	"artk.dev/apperror"
	"artk.dev/ddd"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var _ ddd.AggregateRoot[int64, EntitySerialization] = &Entity{}

type Entity struct {
	id   int64
	name string
}

func (e *Entity) ID() int64 {
	return e.id
}

func (e *Entity) Name() string {
	return e.name
}

func (e *Entity) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return apperror.Validation("names cannot be whitespace")
	}

	e.name = name
	return nil
}

func (e *Entity) Serialize() EntitySerialization {
	return EntitySerialization{
		ID:   e.id,
		Name: e.name,
	}
}

var _ ddd.Serialization[*Entity] = EntitySerialization{}

type EntitySerialization struct {
	ID   int64
	Name string
}

func (s EntitySerialization) Deserialize() *Entity {
	return &Entity{
		id:   s.ID,
		name: s.Name,
	}
}

type EntityRepository interface {
	ddd.CrudRepository[*Entity, int64, EntitySerialization]
}

var _ EntityRepository = &InMemoryEntityRepository{}

type InMemoryEntityRepository struct {
	ddd.InMemoryCrudRepository[*Entity, int64, EntitySerialization]
}

func NewInMemoryEntityRepository() *InMemoryEntityRepository {
	r := &InMemoryEntityRepository{}
	r.Reset()
	return r
}

func TestInMemoryCrudRepository_Get_not_found(t *testing.T) {
	r := NewInMemoryEntityRepository()

	a, err := r.Get(context.TODO(), example().ID())
	if !apperror.IsNotFound(err) {
		t.Error("found unexpected item")
	}
	if a != nil {
		t.Error("expected aggregate to be nil")
	}
}

func TestInMemoryCrudRepository_Get_found(t *testing.T) {
	r := NewInMemoryEntityRepository()

	original := example()
	givenItExists(t, r, original)

	retrieved, err := r.Get(context.TODO(), original.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if retrieved == nil {
		t.Fatal("expected aggregate to be not nil")
	}
	if same(retrieved, original) {
		t.Error("unexpected memory sharing")
	}
	if id := retrieved.ID(); id != original.ID() {
		t.Error("different ID:", id)
	}
	if name := retrieved.Name(); name != original.Name() {
		t.Error("different name:", name)
	}
}

func TestInMemoryCrudRepository_Insert_succeeds_for_new(t *testing.T) {
	r := NewInMemoryEntityRepository()

	entity := example()
	givenItDoesNotExist(t, r, entity.ID())

	err := r.Insert(context.TODO(), entity)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	thenItExists(t, r, entity.ID())
}

func TestInMemoryCrudRepository_Insert_conflict_if_duplicate(t *testing.T) {
	r := NewInMemoryEntityRepository()

	entity := example()
	givenItExists(t, r, entity)

	err := r.Insert(context.TODO(), entity)
	if !apperror.IsConflict(err) {
		t.Fatal("expected conflict error")
	}
}

func TestInMemoryCrudRepository_Update_succeeds_for_existing(t *testing.T) {
	r := NewInMemoryEntityRepository()

	original := example()
	givenItExists(t, r, original)

	const newName = "The Answer to Life, the Universe, and *"
	err := r.Update(context.TODO(), original.ID(), func(x *Entity) error {
		return x.Rename(newName)
	})
	if err != nil {
		t.Fatal("unexpected update failure:", err)
	}

	updated, err := r.Get(context.TODO(), original.ID())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if id := updated.ID(); id != original.ID() {
		t.Error("unexpected ID change:", id)
	}
	if name := updated.Name(); name != newName {
		t.Error("unexpected name:", name)
	}
}

func TestInMemoryCrudRepository_Update_not_found_if_missing(t *testing.T) {
	r := NewInMemoryEntityRepository()

	id := example().ID()
	givenItDoesNotExist(t, r, id)

	const newName = "The Answer to Life, the Universe, and *"
	err := r.Update(context.TODO(), id, func(x *Entity) error {
		return x.Rename(newName)
	})
	if !apperror.IsNotFound(err) {
		t.Fatal("missing expected not found error")
	}
}

func TestInMemoryCrudRepository_Update_propagates_errors(t *testing.T) {
	r := NewInMemoryEntityRepository()

	entity := example()
	givenItExists(t, r, entity)

	expectedError := errors.New("expected test error")
	err := r.Update(context.TODO(), entity.ID(), func(_ *Entity) error {
		return expectedError
	})
	if !errors.Is(err, expectedError) {
		t.Fatalf("missing expected error propagation, got %v", err)
	}
}

func TestInMemoryCrudRepository_Upsert_succeeds_for_new(t *testing.T) {
	r := NewInMemoryEntityRepository()

	entity := example()
	givenItDoesNotExist(t, r, entity.ID())

	err := r.Upsert(
		context.TODO(),
		entity.ID(),
		func() (*Entity, error) {
			return entity, nil
		},
		func(_ *Entity) error {
			t.Fatal("unexpected call to update handler")
			return nil
		},
	)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	thenItExists(t, r, entity.ID())
}

func TestInMemoryCrudRepository_Upsert_succeeds_for_existing(t *testing.T) {
	r := NewInMemoryEntityRepository()

	original := example()
	givenItExists(t, r, original)

	const newName = "The Answer to Life, the Universe, and *"
	err := r.Upsert(
		context.TODO(),
		original.ID(),
		func() (*Entity, error) {
			t.Fatal("unexpected call to insert handler")
			return nil, nil
		},
		func(x *Entity) error {
			return x.Rename(newName)
		},
	)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	updated, err := r.Get(context.TODO(), original.ID())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if id := updated.ID(); id != original.ID() {
		t.Error("unexpected ID change:", id)
	}
	if name := updated.Name(); name != newName {
		t.Error("unexpected name:", name)
	}
}

func TestInMemoryCrudRepository_Upsert_propagates_insert_errors(t *testing.T) {
	r := NewInMemoryEntityRepository()

	id := example().ID()
	givenItDoesNotExist(t, r, id)

	expectedError := errors.New("expected test error")
	err := r.Upsert(
		context.TODO(),
		id,
		func() (*Entity, error) {
			return nil, expectedError
		},
		func(_ *Entity) error {
			t.Fatal("unexpected call to update handler")
			return nil
		},
	)
	if !errors.Is(err, expectedError) {
		t.Fatalf("missing expected error propagation, got %v", err)
	}
}

func TestInMemoryCrudRepository_Upsert_propagates_update_errors(t *testing.T) {
	r := NewInMemoryEntityRepository()

	entity := example()
	givenItExists(t, r, entity)

	expectedError := errors.New("expected test error")
	err := r.Upsert(
		context.TODO(),
		entity.ID(),
		func() (*Entity, error) {
			t.Fatal("unexpected call to insert handler")
			return nil, nil
		},
		func(*Entity) error {
			return expectedError
		},
	)
	if !errors.Is(err, expectedError) {
		t.Fatalf("missing expected error propagation, got %v", err)
	}
}

func TestInMemoryCrudRepository_Delete_succeeds_if_found(t *testing.T) {
	r := NewInMemoryEntityRepository()

	entity := example()
	givenItExists(t, r, entity)

	err := r.Delete(context.TODO(), entity.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	thenItDoesNotExist(t, r, entity.ID())
}

func TestInMemoryCrudRepository_Delete_not_found_if_missing(t *testing.T) {
	r := NewInMemoryEntityRepository()

	id := example().ID()
	givenItDoesNotExist(t, r, id)

	err := r.Delete(context.TODO(), id)
	if !apperror.IsNotFound(err) {
		t.Fatal("unexpected error:", err)
	}

	thenItDoesNotExist(t, r, id)
}

func TestInMemoryCrudRepository_NotFound(t *testing.T) {
	r := NewInMemoryEntityRepository()

	const id int64 = 1234
	for _, tc := range []struct {
		Name        string
		constructor func(id int64) string
		msg         string
	}{
		{
			Name:        "default error message",
			constructor: nil,
			msg:         "not found: 1234",
		},
		{
			Name: "custom error message",
			constructor: func(id int64) string {
				return fmt.Sprintf("missing item %v", id)
			},
			msg: "missing item 1234",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			r.Errors.NotFound = tc.constructor
			err := r.NotFound(id)
			equalErrorMessage(t, tc.msg, err)
		})
	}
}

func TestInMemoryCrudRepository_AlreadyExists(t *testing.T) {
	r := NewInMemoryEntityRepository()

	const id int64 = 1234
	for _, tc := range []struct {
		Name        string
		constructor func(id int64) string
		msg         string
	}{
		{
			Name:        "default error message",
			constructor: nil,
			msg:         "already exists: 1234",
		},
		{
			Name: "custom error message",
			constructor: func(id int64) string {
				return fmt.Sprintf("duplicate item %v", id)
			},
			msg: "duplicate item 1234",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			r.Errors.AlreadyExists = tc.constructor
			err := r.AlreadyExists(id)
			equalErrorMessage(t, tc.msg, err)
		})
	}
}

func equalErrorMessage(t *testing.T, expected string, err error) {
	if err == nil {
		t.Fatal("missing expected error")
	}
	if msg := err.Error(); expected != msg {
		t.Errorf("expected %q, got %q", expected, msg)
	}
}

func givenItExists(
	t *testing.T,
	r EntityRepository,
	e *Entity,
) {
	t.Helper()

	err := r.Insert(context.TODO(), e)
	if err != nil {
		t.Fatal("unexpected error in pre-condition:", err)
	}
}

func givenItDoesNotExist(
	t *testing.T,
	r EntityRepository,
	id int64,
) {
	t.Helper()

	_, err := r.Get(context.TODO(), id)
	if !apperror.IsNotFound(err) {
		t.Fatal("found entity expected to be missing")
	}
}

func thenItExists(
	t *testing.T,
	r EntityRepository,
	id int64,
) {
	t.Helper()

	_, err := r.Get(context.TODO(), id)
	if err != nil {
		t.Fatal("expected entity was missing:", id)
	}
}

func thenItDoesNotExist(
	t *testing.T,
	r EntityRepository,
	id int64,
) {
	t.Helper()

	_, err := r.Get(context.TODO(), id)
	if apperror.IsNotFound(err) {
		// Expected situation.
		return
	}
	if err != nil {
		t.Fatal("unexpected error getting entity:", id)
	}

	t.Fatal("unexpectedly found entity:", id)
}

func example() *Entity {
	return EntitySerialization{
		ID:   42,
		Name: "The Answer",
	}.Deserialize()
}

func same[T any](x, y T) bool {
	return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
}
