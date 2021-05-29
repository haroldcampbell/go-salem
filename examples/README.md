# Examples

A list of examples showcasing the different API options of `Salem`

## Getting Started

Here's how you use Salem in the most basic setup.

```
	factory := salem.Mock(examples.Person{})
	results := factory.Execute()
```

Simply pass an instance of a type to `salem.Mock(...)` can call `factory.Execute()` to get an `[]interface{}`.

See [example 1](./ex-1/main.go).

## Properties

It is easy to control the generated properties the mocks.

```
    factory := salem.Mock(examples.Person{}).
		Ensure("FName", "Sammy").  // Constrain the FName field to Sammy
		Ensure("Surname", "Smith") // Constrain the Surname field to Smith

	results := factory.Execute()
```

Use `factory.Ensure(...)` to explicitly set the value for the fields.
In the code snippet `FName` and `Surname` the fields on the mocks will have "Sammy" and "Smith" respectively.

See [example 2](./ex-2/main.go) for more.

## Number of mocks

You can alo control the number of mock generate.

```
	factory := salem.Mock(examples.Person{}).
		Ensure("FName", "Sammy") // Constrain the FName field to Sammy

	factory.WithExactItems(3) // Generates exactly 3 mock Person structs

	results := factory.Execute()
```

You have 3 `factory.WithXXXItems(...)` functions that allow you to directly manipulate the item count.

See [example 3](./ex-3/main.go) for more.

## Nested Properties

A dot('.') notation is also provided to configure nested fields

```
    factory := salem.Mock(examples.Transaction{}).
		Ensure("Car.TransactionGUID", "GUID-153").
		WithExactItems(3)

	results := factory.Execute()
```

See [example 4](./ex-4/main.go).

## Slices

With the `salem.Tap` function you can easily control the length for fields that are slices.

```
    type AddressBook struct {
	    Contact []examples.Person // We want to specific the number of these nested mocks
    }
```

```
    factory := salem.Mock(AddressBook{})
	factory.Ensure("Contact", salem.Tap(). // We tap the field
		Ensure("FName", "Ted"). // Constrain the value
		WithMaxItems(5)) // Congiure the count
	factory.WithExactItems(2)

	results := factory.Execute()
```

See [example 5](./ex-5/main.go) for more.

## Pointers

Salem also provides support for mocking pointers.

```
    type basic struct {
		Tag   *string
		Age   *int
		Money *float32
	}

	factory := salem.Mock(basic{})
	results := factory.Execute()
```

...and pointer slices.

```
    type basic struct {
		Tag   []*string
		Age   []*int
		Money []*float32
	}

	factory := salem.Mock(basic{})
	factory.Ensure("Age", salem.Tap().WithExactItems(3)) // Control the number of elements of the nested slice
	results := factory.Execute()
```

See [example 6](./ex-6/main.go).

## Omit fields

You can easily omit fields with the `factory.Omit(...)` function.

```
    factory := salem.Mock(basic{})
	factory.Omit("Tag"). // Top level field directly in the basic struct
				Omit("SKU.GUID") // Nested fields
```

See [example 7](./ex-7/main.go).
