package main

import (
	"fmt"
	"go-salem"
	"go-salem/examples"

	"github.com/haroldcampbell/go_utils/utils"
)

// Number of mocks Example
// Demonstrating API that control the number of items created.
//
// There are several methods that can be used to control the number of mocks generated.
//
// WithExactItems - exact n itmes
// WithMinItems - range of items between [n, n+upperBounds]
// WithMaxItems - range of items between [0, 1+n)

func main() {
	factory := salem.Mock(examples.Person{}).
		Ensure("FName", "Sammy") // Constrain the FName field to Sammy
	// By default salem will generate 1 mock.
	// This can be changed by using the WithXXXItems functions.

	// Uncomment the different example functions to try different WithXXX options
	// minExample(factory)
	// maxExample(factory)
	// exactExample(factory)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}

func exactExample(factory *salem.Factory) {
	factory.WithExactItems(3) // Generates exactly 3 mock Person structs
}
func maxExample(factory *salem.Factory) {
	factory.WithMaxItems(10) // Generates [0, 10] mock Person structs
}
func minExample(factory *salem.Factory) {
	// Generates [3, 13] mock Person structs
	// By defualt WithMinItems generates between n and n + 10 mocks
	factory.WithMinItems(3)

	// To change the upperBounds specify it in the span.
	// In the case below, WithMinItems will generate [n, n+20) mocks
	//factory.WithMinItems(3, 20)
}
