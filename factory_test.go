package salem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type empty struct{}
type simple struct {
	num         int
	PublicField string
}

type basicParent struct {
	Child basicChild
}
type basicChild struct {
	ParentName string
	Name       string
	Age        int

	Attends school
}
type school struct {
	NextOfKin string
	Name      string
}

type money struct {
	Value float32
}
type wallet struct {
	Notes []money
}

type vault struct {
	Notes []*money
}

func Test_FactoryOmit(t *testing.T) {
	f := Mock(school{})

	results := f.Execute()
	actualMock := results[0].(school)
	assert.NotEmpty(t, actualMock.NextOfKin, "expect field to be set")

	f.Omit("NextOfKin")
	results = f.Execute()
	actualMock = results[0].(school)

	assert.Empty(t, actualMock.NextOfKin, "expect field to be empty")
}

func Test_FactoryEnsure(t *testing.T) {
	test_simple_ensure(t)
	test_nested_ensure(t)
	test_slice_ensure(t)
	test_slice_pointer_ensure(t)
	test_sequence_ensure(t)
}

func test_simple_ensure(t *testing.T) {
	requiredValue := "happiness"

	f := Mock(simple{})
	f.Ensure("num", 101)
	f.Ensure("PublicField", requiredValue)

	results := f.Execute()
	actualMock := results[0].(simple)

	assert.Equal(t, 1, len(results), "expect mock to be created")
	assert.Equal(t, 0, actualMock.num, "should not set private fields")
	assert.Equal(t, requiredValue, actualMock.PublicField, "should set public fields required value")
}

func test_nested_ensure(t *testing.T) {
	requiredValue := "Goose BrightSpark"

	f := Mock(basicParent{})
	f.Ensure("Child.ParentName", requiredValue)
	f.Ensure("Child.Attends.NextOfKin", requiredValue)

	results := f.Execute()
	actualMock := results[0].(basicParent)

	assert.Equal(t, requiredValue, actualMock.Child.ParentName, "should set nested public fields required value")
	assert.Equal(t, requiredValue, actualMock.Child.Attends.NextOfKin, "should set nested public fields required value")
}

func test_slice_pointer_ensure(t *testing.T) {
	tap := Tap().WithExactItems(5)

	f := Mock(vault{})
	f.Ensure("Notes", tap)

	results := f.Execute()
	actualMocks := results[0].(vault).Notes

	assert.Equal(t, 5, len(actualMocks), "should set nested slice of pointer")
}

func test_slice_ensure(t *testing.T) {
	tap := Tap().WithExactItems(5)

	f := Mock(wallet{})
	f.Ensure("Notes", tap)

	results := f.Execute()
	actualMocks := results[0].(wallet).Notes

	assert.Equal(t, 5, len(actualMocks), "should set nested slice")
}

func test_sequence_ensure(t *testing.T) {
	requiredValue := "happiness"

	f := Mock(simple{})
	f.EnsureSequence("PublicField", "a", requiredValue, "c")
	f.WithExactItems(5)

	results := f.Execute()
	publicFields := make([]string, len(results))
	for i, v := range results {
		publicFields[i] = v.(simple).PublicField
	}

	expected := []string{"a", requiredValue, "c", "", ""}
	actual := []string{publicFields[0], publicFields[1], publicFields[2], publicFields[3], publicFields[4]}
	assert.Equal(t, requiredValue, publicFields[1], "should be set to sequence value")
	assert.Equal(t, expected, actual, "should set sequence values")
}

func Test_FactoryWithItems(t *testing.T) {
	test_with__items(t)
	test_with_exact_items(t)
	test_with_min_items(t)
	test_with_max_items(t)
}

func test_with__items(t *testing.T) {
	f := Mock(empty{})
	assert.Nil(t, f.plan.run, "expect plan.run to be nil before Execute")
}
func test_with_exact_items(t *testing.T) {
	f := Mock(empty{})

	f.WithExactItems(2)
	assert.Nil(t, f.plan.run, "expect WithExactItems to not change plan.run")

	f.Execute()
	assert.Equal(t, ExactRun, f.plan.run.RunType, "expect ExactRun from WithExactItems")
	assert.Equal(t, 2, f.plan.run.Count, "expect correct run count from WithExactItems")
}

func test_with_min_items(t *testing.T) {
	f := Mock(empty{})

	f.WithMinItems(2)
	assert.Nil(t, f.plan.run, "expect WithMinItems to not change plan.run")

	f.Execute()
	assert.Equal(t, MinRun, f.plan.run.RunType, "expect MinRun from WithMinItems")
	assert.GreaterOrEqual(t, f.plan.run.Count, 2, "expect Count to be 2 or more")
}

func test_with_max_items(t *testing.T) {
	f := Mock(empty{})

	f.WithMaxItems(2) // No need to test the run.Count as it will be based on a random value of n
	assert.Nil(t, f.plan.run, "expect WithMaxItems to not change plan.run")

	f.Execute()
	assert.Equal(t, MaxRun, f.plan.run.RunType, "expect MaxRun from WithExactItems")
}
