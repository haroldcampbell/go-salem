A simple go library for generating mock data based on structs.

## Why

A minimal library that helps to generate structs with mock data when creating unit tests with libraries like `"github.com/stretchr/testify/suite"`.

Creating mocked data is difficult and time consuming. This leads to a lot of boiler plate code resulting in brittle tests.

I wanted a faster way to define mocks based on constrains and structure instead of the exact content of all of the fields.

## Example

Here's a simple example.

You have a `Person` strcut to mock.

```
type Person struct {
	FName   string
	Surname string
	Age     int

	privateField int // This should be ignored
}
```

Simply pass the struct to `salem.Mock` and then run `Execute` on the factor that is created for your struct.

```
	factory := salem.Mock(examples.Person{})
	results := factory.Execute()
```

Sample Output for the example:

```
Salem mocks:
[
  {
    "FName": "BDMHKCTVZMER",
    "Surname": "BGGCEPUWFSKWHGH",
    "Age": 26
  }
]
```

See the [examples](./examples/README.md) folder for more information.

## Features

-   Control the value of public fields that are mocked with `Ensure(...)`
-   Control the number of mocks generated with `WithMinItems()`, `WithMaxItems()` and `WithExactItems()`
-   Mock nested structs automatically
-   Control nested public fields with path name e.g. `ChildField.NestedChild.OtherNestedChild`
