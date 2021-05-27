package salem

import (
	"math/rand"
	"time"
)

func Mock(t interface{}) *factory {
	rand.Seed(time.Now().UnixNano())

	f := factory{rootType: t}
	f.plan = NewPlan()
	f.WithExactItems(1) // Default to 1 item

	return &f
}

func (f *factory) Execute() []interface{} {
	return f.plan.Run(f)
}
