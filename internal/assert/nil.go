package assert

import "reflect"

func Nil(t T, got any) bool {
	t.Helper()

	if !isNil(got) {
		t.Errorf(`expected nil, got %v`, got)
		return false
	}
	return true
}

func NotNil(t T, got any) bool {
	t.Helper()

	if isNil(got) {
		t.Errorf(`expected not nil`)
		return false
	}

	return true
}

func NilPointer[V any](t T, got *V) bool {
	t.Helper()

	if got != nil {
		t.Errorf(`expected nil pointer, got %v`, got)
		return false
	}

	return true
}

func NotNilPointer[V any](t T, got *V) bool {
	t.Helper()

	if got == nil {
		t.Errorf(`expected not nil pointer, got %v`, got)
		return false
	}

	return true
}

func NilSlice[V any](t T, got []V) bool {
	t.Helper()

	if got != nil {
		t.Errorf(`expected nil slice, got %v`, got)
		return false
	}

	return true
}

func NilMap[K comparable, V any](t T, got map[K]V) bool {
	t.Helper()

	if got != nil {
		t.Errorf("expected nil map, got %v", got)
		return false
	}

	return true
}

func NilChan[V any](t T, got chan V) {

}

func NilInputChan[V any](t T, got chan<- V) {

}

func NilOutputChan[V any](t T, got chan<- V) {

}

// isNil checks for nil-ability, considering complex situations such as nil
// interfaces.
func isNil(x any) (ok bool) {
	// Handle nil any.
	if x == nil {
		return true
	}

	// Handle nil values with an associated type.
	v := reflect.ValueOf(x)

	// Mimics reflect.Value.IsNil.
	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Pointer,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
