package assume

import "fmt"

// NotZero panics if the value is zero.
// Only applicable to comparable types.
//
// This function can be used to check comparable types against nil.
func NotZero[V comparable](v V) {
	NotZerof(v, "zero value")
}

// NotZerof panics if the value is zero.
// Only applicable to comparable types.
func NotZerof[V comparable](v V, format string, args ...any) {
	var zero V
	if v == zero {
		panic(fmt.Sprintf(format, args...))
	}
}

// NotNilSlice panics if the slice is nil.
func NotNilSlice[V any](s []V) {
	NotNilSlicef(s, "nil slice")
}

// NotNilSlicef panics if the slice is nil.
func NotNilSlicef[V any](s []V, format string, args ...any) {
	if s == nil {
		panic(fmt.Sprintf(format, args...))
	}
}

// NotNilMap panics if the map is nil.
func NotNilMap[K comparable, V any](m map[K]V) {
	NotNilMapf(m, "nil map")
}

// NotNilMapf panics if the map is nil.
func NotNilMapf[K comparable, V any](m map[K]V, format string, args ...any) {
	if m == nil {
		panic(fmt.Sprintf(format, args...))
	}
}
