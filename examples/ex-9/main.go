// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"go-salem"
	"go-salem/examples"

	"github.com/haroldcampbell/go_utils/utils"
)

// Casting into target array Example
// Quick example showing how to cast results into a target type

func main() {
	factory := salem.Mock(examples.Person{}).WithExactItems(5)

	target := factory.ExecuteToType().([]examples.Person)

	str := utils.PrettyMongoString(target)
	fmt.Printf("Salem ExecuteToType mocks:\n%v\n\n", str)
}
