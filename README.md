A simple go library for generating mock data based on structs.

## Example

`
type Engine struct {
Cylinders int
HP int
SerialNumber string
}

type Car struct {
TransactionGUID string

    Name      string
    Make      string
    Engine    Engine
    IsTwoDoor bool

}

type Transaction struct {
GUID string

    Car       Car
    OwnerName string
    Prices    float32

    privateField int // This should be ignored

}

factory := salem.Mock(Transaction{}).
Ensure("GUID", "GUID-153").
Ensure("Car.TransactionGUID", "GUID-153").
WithExactItems(3)

    results := factory.Execute()

`

Output:

` [ { "GUID": "GUID-153", "Car": { "TransactionGUID": "GUID-153", "Name": "QPHZBKVQBQGATXWYNUSMYWAQHWFVRHDZXTKRMTHKDSHB", "Make": "ZENRFCSTTTKNDKRYRXZRRMYW", "Engine": { "Cylinders": 57, "HP": 53, "SerialNumber": "GWMVVVKXDYWJBYHWRTTHZRCW" }, "IsTwoDoor": false }, "OwnerName": "FZSEGCQSWWMQYGXATWHGAJKTEDSUTTGSVHAHNRDHJZ", "Prices": 0.6369821 }, ... ]`

## Why

This library helps to generate structs with mock data when creating unit tests with libraries like `"github.com/stretchr/testify/suite"`
