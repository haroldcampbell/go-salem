package main

import (
	"fmt"
	"go-salem"
	"go-salem/examples"

	"github.com/haroldcampbell/go_utils/utils"
)

// Nested Properties Example
// Setting the properties of nested fields
//
// Nested public fields can be access by using `.` as a path inidicator.
//
//	Using the example structs below.

//		type Transaction struct {
//			Car       Car // <--
//			//... fields ignore for the sake of the example
//		}
//
//		type Car struct {
//			TransactionGUID string // <--
//			//... fields ignore for the sake of the example
//		}
//
//	Using this call to factory.Ensure("Car.TransactionGUID", "GUID-153").
// 	We can access the `TransactionGUID` field on the nestes `Car` struct field.
//

func main() {
	factory := salem.Mock(examples.Transaction{}).
		Ensure("Car.TransactionGUID", "GUID-153").
		WithExactItems(3)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}
