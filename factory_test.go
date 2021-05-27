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

	Attends School
}
type School struct {
	NextOfKin string
	Name      string
}

func Test_FactoryEnsure(t *testing.T) {
	test_simple_ensure(t)
	test_nested_ensure(t)
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

func Test_FactoryWithItems(t *testing.T) {
	f := Mock(empty{})
	test_with_items(t, f)
}

func test_with_items(t *testing.T, f *factory) {
	f.WithExactItems(2)
	assert.Equal(t, ExactRun, f.plan.run.RunType, "expect ExactRun from WithExactItems")
	assert.Equal(t, 2, f.plan.run.Count, "expect correct run count from WithExactItems")

	f.WithMinItems(2) // No need to test the run.Count as it will be based on a random value of n
	assert.Equal(t, MinRun, f.plan.run.RunType, "expect MinRun from WithExactItems")

	f.WithMaxItems(2) // No need to test the run.Count as it will be based on a random value of n
	assert.Equal(t, MaxRun, f.plan.run.RunType, "expect MaxRun from WithExactItems")
}
