package salem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type money struct {
	Value float32
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
	type basicChild struct {
		ParentName string
		Name       string
		Age        int

		Attends school
	}

	type basicParent struct {
		Child basicChild
	}

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
	type vault struct {
		Notes []*money
	}

	tap := Tap().WithExactItems(5)

	f := Mock(vault{})
	f.Ensure("Notes", tap)

	results := f.Execute()
	actualMocks := results[0].(vault).Notes

	assert.Equal(t, 5, len(actualMocks), "should set nested slice of pointer")
}

func test_slice_ensure(t *testing.T) {
	type wallet struct {
		Notes []money
	}

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
