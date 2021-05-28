package examples

type Person struct {
	FName   string
	Surname string
	Age     int

	privateField int // This should be ignored
}

type Engine struct {
	Cylinders    int
	HP           int
	SerialNumber string
}

type HeadLight struct {
	Watt    int
	Voltage int
	Color   string
}
type Car struct {
	TransactionGUID string

	Name       string
	Make       string
	Engine     Engine
	HeadLights []HeadLight
	IsTwoDoor  bool
}

type Transaction struct {
	GUID string

	Car       Car
	OwnerName string
	Prices    float32

	privateField int // This should be ignored
}
