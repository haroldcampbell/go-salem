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

	rt := reflect.TypeOf(f.rootType)

	if rt.Kind() == reflect.Ptr { // For the when Mock is given a pointer. Eg salem.Mock(&groupByColumn{})
		for i, s := range results {
			// Convert the value to a pointer
			ps := reflect.New(reflect.TypeOf(s))
			c := reflect.ValueOf(s)
			ps.Elem().Set(c)
			// Since we are dealing with a pointer slice
			slice.Index(i).Set(ps)
		}

		return slice.Interface()
	}

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
// The sequence items are based their item index in the overall item list.
// The values default to their empty value if the items exceed the number of squence items.
func (f *Factory) EnsureSequence(fieldName string, seq ...interface{}) *Factory {
	f.plan.EnsureSequence(fieldName, seq)

	return f
}

// EnsureSequenceAcross is used to specify the actual values for the fields.
// The sequence items are based on their sequence index as they are generated.
// The values default to their empty value if the items exceed the number of squence items.
func (f *Factory) EnsureSequenceAcross(fieldName string, seq ...interface{}) *Factory {
	f.plan.EnsureSequenceAcross(fieldName, seq)

	return f
}

// EnsureMapKeySequence sets one or more keys for a map field
func (f *Factory) EnsureMapKeySequence(fieldName string, seq ...interface{}) *Factory {
	f.plan.EnsureMapKeySequence(fieldName, seq)

	return f
}

// EnsureMapValueSequence set one or more values for the map field
func (f *Factory) EnsureMapValueSequence(fieldName string, seq ...interface{}) *Factory {
	f.plan.EnsureMapValueSequence(fieldName, seq)

	return f
}

// WithMinItems generates at least n items
// By default WithMinItems will generated [n, n+10) items.
//
// To change the upperBounds from n+10 specifiy the span which will change the range to [n, n+span)
// Only the first value of the span slice is used and it must be > 0.
// In other words, WithMinItems(n, span, ignored, ignored, ...)
func (f *Factory) WithMinItems(n int, span ...int) *Factory {
	f.plan.SetItemCountHandler(func() {
		f.plan.SetRunCount(MinRun, minItem(n, span))
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

// WithExactMapItems generates exactly n items for a field that is a map
func (f *Factory) WithExactMapItems(fieldName string, n int) *Factory {
	f.plan.SetMapItemCountHandler(fieldName, func() {
		f.plan.SetMapRunCount(fieldName, ExactRun, n)
	})
	return f
}

// WithMaxMapItems generates up to [0, n] items for a field that is a map
func (f *Factory) WithMaxMapItems(fieldName string, n int) *Factory {
	f.plan.SetMapItemCountHandler(fieldName, func() {
		f.plan.SetMapRunCount(fieldName, MaxRun, rand.Intn(1+n))
	})
	return f
}

// WithMinMapItems generates at least n items for a field that is a map
// By default WithMinMapItems will generated [n, n+10) items.
//
// See @WithMinItems for more discussion
func (f *Factory) WithMinMapItems(fieldName string, n int, span ...int) *Factory {
	f.plan.SetMapItemCountHandler(fieldName, func() {
		f.plan.SetMapRunCount(fieldName, MinRun, minItem(n, span))
	})
	return f
}

// minItem returns a an in generated [n, n+10) items or [n, n+span[0])
func minItem(n int, span []int) int {
	var upperBounds int = 10

	if len(span) > 0 && span[0] > 0 {
		upperBounds = span[0]
	}

	return n + rand.Intn(upperBounds)
}
