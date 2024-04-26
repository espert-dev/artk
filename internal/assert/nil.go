package assert

import "reflect"

func Nil(t T, got any) bool {
	t.Helper()

	if !isNil(got) {
		report1(t, "expected nil", got)
		return false
	}
	return true
}

func NotNil(t T, got any) bool {
	t.Helper()

	if isNil(got) {
		report1(t, "expected not nil", got)
		return false
	}

	return true
}

func NilPointer[V any](t T, got *V) bool {
	t.Helper()

	if got != nil {
		report1(t, "expected nil pointer", got)
		return false
	}

	return true
}

func NotNilPointer[V any](t T, got *V) bool {
	t.Helper()

	if got == nil {
		report1(t, `expected not nil pointer`, got)
		return false
	}

	return true
}

func NilSlice[V any](t T, got []V) bool {
	t.Helper()

	if got != nil {
		report1(t, "expected nil slice", got)
		return false
	}

	return true
}

func NilMap[K comparable, V any](t T, got map[K]V) bool {
	t.Helper()

	if got != nil {
		report1(t, "expected nil map", got)
		return false
	}

	return true
}

func NilChan[V any](t T, got chan V) bool {
	t.Helper()

	if got != nil {
		report1(t, "expected nil chan", got)
		return false
	}

	return true
}

func NilInputChan[V any](t T, got chan<- V) bool {
	t.Helper()

	if got != nil {
		report1(t, "expected nil chan<-", got)
		return false
	}

	return true
}

func NilOutputChan[V any](t T, got chan<- V) bool {
	t.Helper()

	if got != nil {
		report1(t, "expected nil <-chan", got)
		return false
	}

	return true
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
	defer func() {
		// v.IsNil will panic if the type is not nil-able. This is
		// inconvenient, because there is no way to define a constraint
		// for nil-able. So, rather than duplicate the switch in IsNil,
		// we have opted to install a panic handler instead.
		if r := recover(); r != nil {
			ok = false
		}
	}()

	return v.IsNil()
}
