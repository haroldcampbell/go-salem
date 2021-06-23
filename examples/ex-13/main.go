// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"go-salem"

	"github.com/haroldcampbell/go_utils/utils"
)

type person struct {
	Name    string
	Surname string
}

// Here we show how to use a function handler to create a field value.
// This is useful we we want to inject values dynamically from some other source
func main() {
	simpleHandler()
	advancedHandler()
}

func simpleHandler() {
	names := []string{"Mary", "lilo", "Frank"}
	fieldHandler := func(itemIndex int) interface{} {
		return names[itemIndex]
	}

	factory := salem.Mock(person{}).
		OnField("Name", fieldHandler).
		WithExactItems(3)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "simpleHandler with name from handler", str)
}

type resource struct {
	currReaderIndex int
	mockFileData    []string
}

func (r *resource) getNextItemName(itemIndex int) interface{} {
	if len(r.mockFileData) == 0 {
		panic("no files to load")
	}
	if r.currReaderIndex > len(r.mockFileData) {
		r.currReaderIndex = 0
	}

	nextFile := r.mockFileData[r.currReaderIndex]
	r.currReaderIndex++

	return nextFile
}

func advancedHandler() {
	res := resource{
		mockFileData: []string{"Mary", "lilo", "Frank"},
	}
	factory := salem.Mock(person{}).
		OnField("Name", res.getNextItemName).
		WithExactItems(3)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", "advancedHandler example with handler", str)
}
