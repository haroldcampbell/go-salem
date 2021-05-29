package salem

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Plan(t *testing.T) {
	test_default_generators(t)
	test_ensured_field_value(t)
}

func test_default_generators(t *testing.T) {
	p := NewPlan()

	assert.Equal(t, len(p.generators), 9, "expect all default generators to created")

	assert.NotEmpty(t, p.GetKindGenerator(reflect.Bool), "expect generator for reflect.Bool")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Int), "expect generator for reflect.Int")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Int8), "expect generator for reflect.Int8")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Int16), "expect generator for reflect.Int16")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Int32), "expect generator for reflect.Int32")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Int64), "expect generator for reflect.Int64")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Float32), "expect generator for reflect.Float32")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.Float64), "expect generator for reflect.Float64")
	assert.NotEmpty(t, p.GetKindGenerator(reflect.String), "expect generator for reflect.String")
}

func test_ensured_field_value(t *testing.T) {
	p := NewPlan()
	fieldName := "add"
	expected := 10

	p.EnsuredFieldValue(fieldName, expected)
	actual := p.ensuredFields[fieldName].fieldAction()

	assert.Equal(t, expected, actual, "expect EnsuredFieldValue(...) to set fuction to return required field value")
}
