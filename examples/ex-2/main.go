package main

import (
	"fmt"
	"go-salem"
	"go-salem/examples"

	"github.com/haroldcampbell/go_utils/utils"
)

// Properties Example
// Controling the generated properties the mocks.
//
// In the example factory.Ensure(...) is used to explicitly set the value for the
// FName and Surname public fields on the mock Person struct.
//
// Running the example multiple times produces different values for the public Age field.
// However, the FName and Surname fields will always return Sammy and Smith respectively.
// This is because they are constrained by factory.Ensure(...).

func main() {
	factory := salem.Mock(examples.Person{}).
		Ensure("FName", "Sammy").  // Constrain the FName field to Sammy
		Ensure("Surname", "Smith") // Constrain the Surname field to Smith

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}
