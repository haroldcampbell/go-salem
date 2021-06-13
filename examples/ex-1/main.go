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

// Getting Started Example
// Using Salem in the most basic setup
func main() {
	factory := salem.Mock(examples.Person{})
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}
