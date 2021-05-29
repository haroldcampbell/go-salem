package main

import (
	"fmt"
	"go-salem"
	"go_utils/utils"
)

// Pointers Example
// Using Salem with points on public fields
func main() {
	primitivePointer()
	objectPointer()

	primitivePointerArray()
	objectPointerArray()

	sideEffects()
}

func primitivePointer() {
	type basic struct {
		Tag   *string
		Age   *int
		Money *float32
	}

	factory := salem.Mock(basic{})
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("primitivePointer mocks:\n%v\n\n", str)
}

func primitivePointerArray() {
	type basic struct {
		Tag   []*string
		Age   []*int
		Money []*float32
	}

	factory := salem.Mock(basic{})
	factory.Ensure("Age", salem.Tap().WithExactItems(3)) // Control the number of elements of the nested slice
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("primitivePointerArray mocks:\n%v\n\n", str)
}

func objectPointer() {
	type note struct {
		Name  string
		Value float32
	}

	type cash struct {
		Note *note
	}

	factory := salem.Mock(cash{})
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("objectPointer mocks:\n%v\n\n", str)
}

func objectPointerArray() {
	type note struct {
		Name  string
		Value float32
	}

	type cash struct {
		Note *note
	}

	type wallet struct {
		Cash []*cash
	}

	factory := salem.Mock(wallet{})
	factory.Ensure("Cash", salem.Tap().WithMaxItems(10)) // Control the number of elements of the nested slice
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("objectPointerArray mocks:\n%v\n\n", str)
}

func sideEffects() {
	factory := salem.Mock("")
	factory.WithExactItems(10) // Make 10 copies based on the primitive
	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("sideEffects mocks:\n%v\n\n", str)
}
