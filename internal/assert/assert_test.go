package assert_test

import (
	"artk.dev/internal/assert"
	"fmt"
	"testing"
)

func TestEqual_values_are_equal_to_themselves(t *testing.T) {
	for _, v := range []any{
		0,
		1,
		"",
		"foo",
	} {
		t.Run(fmt.Sprintf("%v", v), func(t *testing.T) {
			assert.Equal(t, v, v)
		})
	}
}
