package salem

import (
	"math/rand"
	"reflect"
)

type nilStruct struct {
}

type Factory struct {
	rootType interface{}
	plan     *Plan
}

// Tap creates a factory based on the public fields type.
func Tap() *Factory {
	nestedFactory := Mock(nilStruct{})

	return nestedFactory
}

// GetPlan return a point to the current plan
func (f *Factory) GetPlan() *Plan {
	return f.plan
}

// Execute execute the factory instructions to generate the mocks
func (f *Factory) Execute() []interface{} {
	return f.plan.Run(f)
}

// ExecuteToType returns a slice that tis the same type as the Mock's parameter.
//
// This allows easy typecasting into the underlying mocks type.
// Example:
// 		factory := salem.Mock(examples.Person{}).WithExactItems(5)
// 		target := factory.ExecuteToType().([]examples.Person) //<- we can do this
func (f *Factory) ExecuteToType() interface{} {
	results := f.plan.Run(f)

	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(f.rootType)), len(results), len(results))

	for i, s := range results {
		c := reflect.ValueOf(s)

		slice.Index(i).Set(c)
	}

	return slice.Interface()
}

// Omit fields with the specifed name
func (f *Factory) Omit(fieldName string) *Factory {
	f.plan.OmitField(fieldName)

	return f
}

// Ensure sets the value of the fields we don't want to randomly generate
func (f *Factory) Ensure(fieldName string, sharedValue interface{}) *Factory {

	switch sharedValue.(type) {
	case *Factory:
		f.plan.EnsuredFactoryFieldValue(fieldName, sharedValue)

	default:
		f.plan.EnsuredFieldValue(fieldName, sharedValue)
	}

	return f
}

// EnsureConstraint set a constraint that limits the generated value.
//
// The constraint panics if there is an f.Ensure(...) which generates a value resulting in a false constraint
//
// Alternatively, the method also fails after trying to generate a constraint after
// several attempts.
//
// The default attempts is defined by SuggestedConstraintRetryAttempts.
// Use f.GetPlan().SetMaxConstraintsRetryAttempts(...) the change the number of retry attempts.
func (f *Factory) EnsureConstraint(fieldName string, constraint FieldConstraint) *Factory {
	f.plan.EnsuredFieldValueConstraint(fieldName, constraint)

	return f
}

// EnsureSequence is used to specify the actual values for the fields.
// The values default to their empty value if the items exceed the number of squence items.
func (f *Factory) EnsureSequence(fieldName string, seq ...interface{}) *Factory {
	f.plan.EnsureSequence(fieldName, seq)

	return f
}

// WithMinItems generates at least n items
// By default WithMinItems will generated [n, n+10) items. To change the upperBounds
// specifiy the span.
//
// Only the first value of the span slice is used and it must be > 0.
func (f *Factory) WithMinItems(n int, span ...int) *Factory {
	f.plan.SetItemCountHandler(func() {
		var upperBounds int = 10

		if len(span) > 0 && span[0] > 0 {
			upperBounds = span[0]
		}

		f.plan.SetRunCount(MinRun, n+rand.Intn(upperBounds))
	})

	return f
}

// WithMaxItems generates up to [0, n] items
func (f *Factory) WithMaxItems(n int) *Factory {
	f.plan.SetItemCountHandler(func() {
		f.plan.SetRunCount(MaxRun, rand.Intn(1+n))
	})

	return f
}

// WithExactItems generates exactly n items
func (f *Factory) WithExactItems(n int) *Factory {
	f.plan.SetItemCountHandler(func() {
		f.plan.SetRunCount(ExactRun, n)
	})

	return f
}
