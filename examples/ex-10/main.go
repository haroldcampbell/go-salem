package main

import (
	"fmt"
	"go-salem"
	"go-salem/examples"

	"github.com/haroldcampbell/go_utils/utils"
)

// ConstrainStringLength example
// Simple example showing how to set the constraint for a string field.
// Salem will try to generate a value that falls within the bounds.
//
// The default attempts is defined by salem.SuggestedConstraintRetryAttempts.
// Use factory.GetPlan().SetMaxConstraintsRetryAttempts(...) the change the number of retry attempts.

func main() {
	factory := salem.Mock(examples.Person{}).
		EnsureConstraint("FName", salem.ConstrainStringLength(4, 10)).
		EnsureConstraint("Surname", salem.ConstrainStringLength(4, 10)).
		WithExactItems(5)

	target := factory.ExecuteToType().([]examples.Person)

	str := utils.PrettyMongoString((target))
	fmt.Printf("Salem ExecuteToType mocks:\n%v\n\n", str)
}
