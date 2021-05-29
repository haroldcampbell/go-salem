package main

import (
	"fmt"
	"go-salem"
	"go_utils/utils"
)

type II struct {
	I int
}
type Counter struct {
	// Tag *string
	Age  *int
	Age2 int
	// Ages     []*int
	Money    *float32
	Index    *II
	Indices1 []II
	Indices2 []*II
}

// Points Example
// Using Salem with points on public fields

func main() {
	factory := salem.Mock(Counter{})
	// factory.EnsureCounter("Index").
	// WithExactItems(4)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Salem mocks:\n%v\n\n", str)
}
