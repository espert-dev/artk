// Package assume provides assertions that panic on violation.
//
// This serves two purposes:
//  1. Failing fast.
//  2. Remove unnecessary branches and the temptation to test them.
//
// # Nil
//
// Detecting nil is a bit of a mess in go because of nil interfaces.
// We have decided to avoid generic predicates such as NotNil in favour
// of the much less error-prone NotZero, NotNilSlice, NotNilMap, etc.
// These functions use generics to convey the narrow type instead of any,
// which is essential to avoid some important pitfalls.
package assume
