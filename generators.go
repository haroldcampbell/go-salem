package salem

import (
	"hmc/utils"
	"math"
	"math/rand"
	"reflect"
)

type GenType = func() interface{}

func (p *Plan) initDefaultGenerators() {
	p.generators[reflect.Bool] = randBool

	p.generators[reflect.Int] = randInt
	p.generators[reflect.Int8] = randInt8
	p.generators[reflect.Int16] = randInt16
	p.generators[reflect.Int32] = randInt32
	p.generators[reflect.Int64] = randInt64

	p.generators[reflect.Float32] = randFloat32
	p.generators[reflect.Float64] = randFloat64

	p.generators[reflect.String] = randString
}

func (p *Plan) GetKindGenerator(k reflect.Kind) GenType {
	return p.generators[k]
}

func randBool() interface{} {
	if rand.Intn(1000000)%2 == 0 {
		return true
	}

	return false
}
func randInt() interface{} {
	return rand.Intn(math.MaxInt8)
}
func randInt8() interface{} {
	return rand.Intn(math.MaxInt8)
}
func randInt16() interface{} {
	return rand.Int31n(math.MaxInt16)
}
func randInt32() interface{} {
	return rand.Int31n(math.MaxInt32)
}
func randInt64() interface{} {
	return rand.Int63n(math.MinInt64)
}
func randFloat32() interface{} {
	return rand.Float32()
}
func randFloat64() interface{} {
	return rand.Float64()
}
func randString() interface{} {
	len := rand.Intn(50)
	return utils.RandCharacters(3 + len) // Ensure we always have at least 3 chars
}
