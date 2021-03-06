// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"go-salem"

	"github.com/haroldcampbell/go_utils/utils"
)

type basic struct {
	SKU string
}

// Sequences Example
// An example where you can specify the actual values for the fields
// The sequence items are chosen based on the item index

func main() {
	factory := salem.Mock(basic{})
	factory.EnsureSequence("SKU", "a", "b", "c").
		WithExactItems(5) // Additional items are set to an empty string

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Fields from sequences mocks:\n%v\n\n", str)
}
