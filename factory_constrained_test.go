package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConstraint struct {
	didCallIsValid bool
	actualField    interface{}
}

func (t *testConstraint) IsValid(field interface{}) bool {
	t.didCallIsValid = true
	t.actualField = field

	return true
}

func Test_FactoryConstrainedEnsure(t *testing.T) {
	test_constrained_ensure(t)
	test_string_constrained_ensure(t)
	test_constraint_clash_with_ensure(t)
}

func test_constrained_ensure(t *testing.T) {
	type human struct {
		Name string // Field we are checking
	}

	tc := &testConstraint{}
	assert.False(t, tc.didCallIsValid)
	assert.Empty(t, tc.actualField)

	f := salem.Mock(human{})
	f.EnsureConstraint("Name", tc)
	f.Execute()

	assert.True(t, tc.didCallIsValid, "should call IsValid(...) on constraint handler")
	assert.NotEmpty(t, tc.actualField, "should pass field to constraint")
}

func test_constraint_clash_with_ensure(t *testing.T) {
	minBound := 4
	maxBound := 20

	type human struct {
		Name string
	}
	f := salem.Mock(human{})
	f.Ensure("Name", "xxxxx xxxxx xxxxx xxxxx xxxxx yyyy") // This is clashes with the maxBound and should cause a panic
	f.EnsureConstraint("Name", salem.ConstrainStringLength(minBound, maxBound))

	assert.Panics(t, func() {
		f.Execute()
	}, "should panic when an Ensure(...) clashes with an EnsureContraint(...)")
}

func test_string_constrained_ensure(t *testing.T) {
	type human struct {
		Name string // We want this to be 4 - 20 charaters
		Age  int
	}
	minBound := 4
	maxBound := 20
	f := salem.Mock(human{})
	f.EnsureConstraint("Name", salem.ConstrainStringLength(minBound, maxBound))

	results := f.Execute()
	actualMock := results[0].(human)

	assert.True(t, len(actualMock.Name) >= minBound, "should constraint field to lower-bound")
	assert.True(t, len(actualMock.Name) <= maxBound, "should constraint field to upper-bound")
}
