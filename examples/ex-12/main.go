// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"go-salem"

	"github.com/haroldcampbell/go_utils/utils"
)

type item struct {
	SKU  string
	Name string
}

type transaction struct {
	Item      item
	Qty       int
	UnitPrice float32
}

type database struct {
	Lookup map[string]string
}

// Maps Examples

func main() {
	basic_map()
	map_with_slice()
	map_with_exact_items()
	map_with_min_items()
	map_with_max_items()
	map_ensure_keys()
	map_ensure_values()
	map_ensure_keys_and_values()
}

func basic_map() {
	factory := salem.Mock(database{})
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "basic map", str)
}

func map_with_slice() {
	type staff struct {
		Sales map[string][]transaction // staff ID -> transactions
	}

	factory := salem.Mock(staff{})
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "map with slice", str)
}

func map_with_exact_items() {
	factory := salem.Mock(database{}).
		WithExactMapItems("Lookup", 5) // Generate exactly 5 items

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "map with exact item", str)
}

func map_with_min_items() {
	factory := salem.Mock(database{}).
		WithMinMapItems("Lookup", 5) // Generate at least 5 items

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "map with min item", str)
}

func map_with_max_items() {
	factory := salem.Mock(database{}).
		WithMaxMapItems("Lookup", 5) // Generate between 0 - 5 items

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "map with max item", str)
}

var keys = []interface{}{"2050391", "1705598", "22892120", "30716354", "33119748"}
var values = []interface{}{
	"How to check if a map contains a key in Go?",
	"VS2008 : Start an external program on Debug",
	"How to generate a random string of a fixed length in Go?",
	"How do I do a literal *int64 in Go?",
	"Convert time.Time to string",
}

func map_ensure_keys() {
	results := salem.Mock(database{}).
		EnsureMapKeySequence("Lookup", keys...). // Sets the keys and mock the values
		WithExactMapItems("Lookup", 5).
		Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "Ensuring map keys. Values are mocked", str)
}

func map_ensure_values() {
	results := salem.Mock(database{}).
		EnsureMapValueSequence("Lookup", values...). // Sets the values and mock the keys
		WithExactMapItems("Lookup", 5).
		Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "Ensuring map values. Keys are mocked", str)
}

func map_ensure_keys_and_values() {
	results := salem.Mock(database{}).
		EnsureMapValueSequence("Lookup", values...). // Sets the values
		EnsureMapKeySequence("Lookup", keys...).     // Sets the keys
		WithExactMapItems("Lookup", 5).
		Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "Ensuring map keys and values", str)
}
