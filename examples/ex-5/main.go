package main

import (
	"fmt"
	"go-salem"
	"go-salem/examples"

	"github.com/haroldcampbell/go_utils/utils"
)

type AddressBook struct {
	Contact []examples.Person
}

// Slices Example
// Using Salem to control the length of slices
func main() {
	simple()
	// nested()
}

func simple() {
	factory := salem.Mock(AddressBook{})
	factory.Ensure("Contact", salem.Tap().
		Ensure("FName", "Ted").
		WithMaxItems(5))
	factory.WithExactItems(2)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}

func nested() {
	factory := salem.Mock(examples.Transaction{})
	factory.Ensure("Car.HeadLights", salem.Tap().
		WithExactItems(2))

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}
