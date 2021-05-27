package salem

import (
	"math/rand"
)

type factory struct {
	rootType interface{}
	plan     *Plan
}

// Ensure sets the value of the fields we don't want to randomly generate
func (f *factory) Ensure(fieldName string, sharedValue interface{}) *factory {
	f.plan.RequireFieldValue(fieldName, sharedValue)

	return f
}

// WithMinItems generates at least n items
func (f *factory) WithMinItems(n int) *factory {
	f.plan.SetRunCount(MinRun, rand.Intn(n+10))

	return f
}

// WithMaxItems generates up to n items
func (f *factory) WithMaxItems(n int) *factory {
	f.plan.SetRunCount(MaxRun, rand.Intn(1+n))

	return f
}

// WithExactItems generates exactly n items
func (f *factory) WithExactItems(n int) *factory {
	f.plan.SetRunCount(ExactRun, n)

	return f
}
