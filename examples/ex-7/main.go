package main

import (
	"fmt"
	"go-salem"

	"github.com/haroldcampbell/go_utils/utils"
)

type sku struct {
	GUID   string
	Region string
}

type basic struct {
	Tag  *string
	Cost *float32
	SKU  sku
}

// Omitting fields example
// Use the factory.Omit(...) function to exclude generating values for the specified public fields.
func main() {
	factory := salem.Mock(basic{})
	factory.Omit("Tag"). // Top level field directly in the basic struct
				Omit("SKU.GUID") // Nested fields

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("Omitting fields mocks:\n%v\n\n", str)
}
