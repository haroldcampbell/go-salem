package salem

import (
	"math/rand"
	"time"
)

// By default Mock is configured to generate 1 mock.
// This can be changed by using the factory.WithXXXItems(...) functions.
func Mock(t interface{}) *Factory {
	rand.Seed(time.Now().UnixNano())

	f := Factory{rootType: t}
	f.plan = NewPlan()
	f.WithExactItems(1) // Default to 1 item

	return &f
}
