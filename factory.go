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
// By default WithMinItems will generated [n, n+10) items. To change the upperBounds
// specifiy the span.
//
// Only the first value of the span slice is used and it must be > 0.
func (f *factory) WithMinItems(n int, span ...int) *factory {
	var upperBounds int = 10

	if len(span) > 0 && span[0] > 0 {
		upperBounds = span[0]
	}

	f.plan.SetRunCount(MinRun, n+rand.Intn(upperBounds))

	return f
}

// WithMaxItems generates up to [0, n] items
func (f *factory) WithMaxItems(n int) *factory {
	f.plan.SetRunCount(MaxRun, rand.Intn(1+n))

	return f
}

// WithExactItems generates exactly n items
func (f *factory) WithExactItems(n int) *factory {
	f.plan.SetRunCount(ExactRun, n)

	return f
}
