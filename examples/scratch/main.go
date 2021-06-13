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

type friend struct {
	Age    int
	Money  float32
	Region string
}

type person struct {
	Age     int
	Name    string
	Surname string
	Friends []friend
}

func main() {
	factory := salem.Mock(person{}).
		EnsureSequence("Surname", "Campbell", "Wu", "Barret").
		Ensure("Friends", salem.Tap().
					EnsureSequence("Friends.Region", "China", "Jamaica", "Mauritius").
					WithExactItems(3)).
		WithExactItems(2) // Additional items are set to an empty string

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Fields from sequences mocks:\n%v\n\n", str)
}
