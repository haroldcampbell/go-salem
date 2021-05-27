package main

import (
	"fmt"
	"go-salem"

	"github.com/haroldcampbell/go_utils/utils"
)

type Engine struct {
	Cylinders    int
	HP           int
	SerialNumber string
}

type Car struct {
	TransactionGUID string

	Name      string
	Make      string
	Engine    Engine
	IsTwoDoor bool
}

type Transaction struct {
	GUID string

	Car       Car
	OwnerName string
	Prices    float32

	privateField int // This should be ignored
}

func main() {
	factory := salem.Mock(Transaction{}).
		Ensure("GUID", "GUID-153").
		Ensure("Car.TransactionGUID", "GUID-153").
		WithExactItems(3)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}
