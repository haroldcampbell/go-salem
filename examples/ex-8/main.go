package main

import (
	"fmt"
	"go-salem"
	"go_utils/utils"
)

type basic struct {
	SKU string
}

// Sequences Example
// An example where you can specify the actual values for the fields

func main() {
	factory := salem.Mock(basic{})
	factory.EnsureSequence("SKU", "a", "b", "c").
		WithExactItems(5) // Additional items are set to an empty string

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Fields from sequences mocks:\n%v\n\n", str)
}
