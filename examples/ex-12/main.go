package main

import (
	"fmt"
	"go-salem"
	"go_utils/utils"
)

type item struct {
	SKU  string
	Name string
}

type transaction struct {
	UnitPrice float32
	Item      item
	Qty       int
}

// Maps Examples

func main() {
	basic_map()
	map_with_slice()
}

func basic_map() {
	type basic struct {
		Lookup map[string]string
	}

	factory := salem.Mock(basic{})

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "basic map", str)
}

func map_with_slice() {
	type staff struct {
		StaffSales map[string][]transaction // staff ID -> transactions
	}

	factory := salem.Mock(staff{})

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "map with slice", str)
}
