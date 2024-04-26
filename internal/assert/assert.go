// Package assert is a re-imagination of github.com/stretchr/testify/assert
// using generics.
//
// # Pros and cons of using generics
//
// The first advantage of using generics is that we can turn what would be
// runtime errors in testify into compile time errors. This is especially
// significant for negated assertions (e.g., NotEqual), which could become
// misleading. Think, for example:
//
//	x := int64(42)
//	assert.NotEqual(t, 42, x)
//
// In testify, the above code would pass. In this library, it would fail at
// compile time.
//
// The second advantage of using generics is a more intuitive experience when
// using typeless literals. Consider the below example:
//
//	x := int64(42)
//	assert.Equal(t, 42, x)
//
// In testify, the above code would fail, because the value 42 would be cast
// to int, instead of the int64 type of x. In this library, the literal would
// be cast to the type of x, and the test would pass.
//
// A third advantage of generics is a variation of the previous case. By using
// generics, we can replace some guaranteed runtime errors with compilation
// errors. For example:
//
//	x := int(42)
//	y := int64(42)
//	assert.Equal(t, 42, 42)
//
// Since the types are different, neither testify nor this library will accept
// this code. However, testify will fail at runtime while this library will
// instead fail to compile, leading to a faster feedback loop.
//
// On the other hand, using generics imposes some limitations that are not
// present in testify. The main one is that Go does not allow generic methods.
// This means that we would not be able to define, for example, testify's
// Assertion, or the assertion methods in testify's suite.Suite. As a result,
// every assertion in this package will instead be a free function.
//
// # Additional differences with other libraries
//
// This package has no external dependencies. By comparison, as of the time of
// writing, testify depends on two libraries that have not been maintained in
// multiple years. It also depends on a YAML library, which makes for a harder
// security posture.
//
// Testify tries to remain compatible with Go versions beyond their end of life.
// While laudable in some ways, it can lead to a less intuitive experience. For
// example, testify's assert.TestingT does not declare the Helper function as
// a requirement, relying on a sneaky dynamic cast to check support. While it
// is possible to document this kind of behavior, comments will never be quite
// as discoverable as code. For this reason, we have decided to limit our
// support window to track the official support window for Go versions, and
// keep functionality as simple as possible.
//
// Finally, we have opted to make our functions more closely tied to the
// details of the Go type system. For example, the function Equal on this
// package only supports types that satisfy the comparable constraint, and
// anything else will lead to a compilation error. If you need to compare
// slices, you can use Same to check if they share memory, or EqualSlices to
// recursively compare their contents.
//
// # Removed functionality
//
// This library has no analogue of testify's http package, which is deprecated.
// Use the standard package httptest instead.
//
// It also does not have an analogue of testify's mock package. For most cases,
// we instead recommend using configurable test doubles. If you want to mock
// an interface, you can define a struct:
//
//	type Fooer interface {
//		Foo(x int) int
//	}
//
//	type FooTestDouble struct {
//		FooFn func(int) int
//	}
//
//	func (d FooerTestDouble) Foo(x int) int {
//		return d.FooFn(x)
//	}
//
// which then can be configured with anonymous functions within tests:
//
//	func TestFoo(t *testing.T) {
//		// ...
//
//		var fooer FooerTestDouble
//		fooer.FooFn = func(x int) int {
//			// Your mocked functionality here.
//		}
//
//		// ...
//
//		DoSomething(fooer)
//
//		// ...
//	}
//
// For more complex cases, such as testing complex interactions between
// services, we suggest that you consider an integration test instead.
package assert

import (
	"fmt"
	"strings"
)

// T abstracts the functionality of testing.T used by this package.
type T interface {
	Helper()
	Error(args ...any)
}

func report2(t T, message string, expected any, actual any) {
	t.Helper()

	var output strings.Builder
	output.WriteString("\nError:     ")
	output.WriteString(message)
	output.WriteString("\nExpected:  ")
	output.WriteString(fmt.Sprintf("%+v", expected))
	output.WriteString("\nActual:    ")
	output.WriteString(fmt.Sprintf("%+v", actual))
	output.WriteRune('\n')

	t.Error(output.String())
}

func report1(t T, message string, value any) {
	t.Helper()

	var output strings.Builder
	output.WriteString("\nError:     ")
	output.WriteString(message)
	output.WriteString("\nValue:     ")
	output.WriteString(fmt.Sprintf("%+v", value))
	output.WriteRune('\n')

	t.Error(output.String())
}

func report0(t T, message string) {
	t.Helper()

	var output strings.Builder
	output.WriteString("\nError:     ")
	output.WriteString(message)
	output.WriteRune('\n')

	t.Error(output.String())
}
