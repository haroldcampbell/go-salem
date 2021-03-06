# Examples

A list of examples showcasing the different API options of `Salem`

**Note**: Before running the examples you'll need to get the dependency.

```
	go get github.com/haroldcampbell/go_utils/utils
```

This dependency is used for pretty-printting.

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

## Supplying sequences

There are times when you want to supply a secquence of values. This can be done with the `factory.EnsureSequence` function.

```
	factory := salem.Mock(basic{})
	factory.EnsureSequence("SKU", "a", "b", "c", "sss").
		WithExactItems(4)

	results := factory.Execute()
```

See [example 8](./ex-8/main.go).

## Casting to target slices

Cast results into a target types is also possible without the need to write additional copying code.

```
	factory := salem.Mock(examples.Person{}).WithExactItems(5)
	target := factory.ExecuteToType().([]examples.Person)
```

See [example 9](./ex-9/main.go).

## Constrained string fields

Salem allows you generate a values that falls within bounds set by `ConstrainStringLength(min, max)`.

```
	factory := salem.Mock(examples.Person{}).
		EnsureConstraint("FName", salem.ConstrainStringLength(4, 10)).
		EnsureConstraint("Surname", salem.ConstrainStringLength(4, 10)).
		WithExactItems(5)
```

See [example 10](./ex-10/main.go).

## Sequencing across items

There are two general types of sequencing:

-   Based on the items' index with `EnsureSequence`
-   Base on the sequences' index with `EnsureSequenceAcross`

See [example 11](./ex-11/main.go).

## Maps

Salem aslo has the ability to generate maps.

There is support for values that are primitive (e.g. int, string, etc.) types as well as structs and slices.

```
	type basic struct {
		Lookup map[string]string
	}

	factory := salem.Mock(basic{})
	results := factory.Execute()
```

...and slices of structs.

```
	type staff struct {
		StaffSales map[string][]transaction // staff ID -> transactions
	}

	factory := salem.Mock(staff{})
	results := factory.Execute()
```

You also have 3 `factory.WithXXXMapItems(...)` functions to directly manipulate the map item count.

-   `WithExactMapItems("Lookup", 5)`
-   `WithMinMapItems("Lookup", 5)`
-   `WithMaxMapItems("Lookup", 5)`

Example:

```
	type staff struct {
		Lookup map[string]string
	}

	factory := salem.Mock(staff{}).
		WithExactMapItems("Lookup", 5)

	results := factory.Execute()

```

To set the default key and values of public map fields use `EnsureMapKeySequence(...)` and `EnsureMapValueSequence(...)`.

```
	keys := []interface{}{"2050391", "1705598", "22892120", "30716354", "33119748"}
	values := []interface{}{
		"How to check if a map contains a key in Go?",
		"VS2008 : Start an external program on Debug",
		"How to generate a random string of a fixed length in Go?",
		"How do I do a literal *int64 in Go?",
		"Convert time.Time to string",
	}

	results := salem.Mock(database{}).
		EnsureMapValueSequence("Lookup", keys...). // Sets the values
		EnsureMapKeySequence("Lookup", values...). // Sets the keys
		WithExactMapItems("Lookup", 5).
		Execute()
```

See [example 12](./ex-12/main.go).

## Function Handlers

For those occasions that you need to function handler to create a field value, the `OnField(...)` function is perfect.

This is most useful for integrating data from external APIs, databases, files, etc.

```
	names := []string{"Mary", "lilo", "Frank"}
	fieldHandler := func(itemIndex int) interface{} {
		return names[itemIndex]
	}

	factory := salem.Mock(person{}).
		OnField("Name", fieldHandler).
		WithExactItems(3)

	results := factory.Execute()
```

See [example 13](./ex-13/main.go) for more examples
