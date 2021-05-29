package salem

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type RunType uint
type FactoryActionType func(reflect.Value) GenType

const (
	Invalid RunType = iota

	MinRun
	MaxRun
	ExactRun
)

type PlanRun struct {
	RunType RunType
	Count   int
}

type fieldSetter struct {
	fieldAction   GenType
	factoryAction FactoryActionType
}
type Plan struct {
	ensuredFields map[string]fieldSetter // fields set via ensure
	generators    map[reflect.Kind]func() interface{}

	run *PlanRun

	evalItemCountAction func()
	parentName          string
}

func NewPlan() *Plan {
	p := &Plan{}

	p.ensuredFields = make(map[string]fieldSetter)
	p.generators = make(map[reflect.Kind]func() interface{})

	p.initDefaultGenerators()

	return p
}

func (p *Plan) SetItemCountHandler(handler func()) {
	p.evalItemCountAction = handler
}

func (p *Plan) EnsuredFieldValue(fieldName string, sharedValue interface{}) {
	setter := fieldSetter{
		fieldAction: func() interface{} {
			return sharedValue
		},
	}

	p.ensuredFields[fieldName] = setter
}

func (p *Plan) EnsuredFactoryFieldValue(fieldName string, sharedValue interface{}) {
	setter := fieldSetter{
		factoryAction: makeFactoryAction(sharedValue.(*Factory)),
	}

	p.ensuredFields[fieldName] = setter
}

func (p *Plan) SetRunCount(runType RunType, n int) {
	p.run = &PlanRun{
		RunType: runType,
		Count:   n,
	}
}

func (p *Plan) CopyParentRequiredFields(pp *Plan) {
	for k, v := range pp.ensuredFields {
		p.ensuredFields[k] = v
	}
}

func (p *Plan) Run(f *Factory) []interface{} {
	rand.Seed(time.Now().UnixNano())
	p.evalItemCountAction()

	items := make([]interface{}, 0)

	for i := 0; i < p.run.Count; i++ {
		items = append(items, p.generateRandomMock(f))
	}

	return items
}

func (p *Plan) generateRandomMock(f *Factory) interface{} {
	v := reflect.ValueOf(f.rootType)
	typeOfT := v.Type()

	// Create an mock instance of the struct
	newMockPtr := reflect.New(typeOfT)
	newElm := newMockPtr.Elem()

	rand.Seed(time.Now().UnixNano())

	if v.Kind() == reflect.Ptr {
		ptrType := v.Type().Elem()

		newMockPtr = reflect.New(ptrType) // Make an instance based on the pointer type
		newElm = newMockPtr.Elem()

		v = newElm
		typeOfT = v.Type()
	}

	if isPrimativeKind(v.Kind()) {
		generator := p.generators[v.Kind()]
		val := generator()
		return val
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i) // Get field in the struct
		k := field.Kind()
		fieldName := typeOfT.Field(i).Name

		iField := newElm.Field(i) // Get related instance field in the mock instance
		if !iField.CanSet() {
			continue // Skip private instance fields
		}

		qualifiedName := distinctFileName(p.parentName, fieldName)
		generator := p.getValueGenerator(k, iField, qualifiedName)
		val := p.generateFieldValue(k, generator, iField, qualifiedName)

		iField.Set(val)
	}
	return newMockPtr.Elem().Interface()
}

func (p *Plan) getValueGenerator(k reflect.Kind, iField reflect.Value, qualifiedName string) GenType {
	if p.ensuredFields[qualifiedName].factoryAction != nil {
		return p.ensuredFields[qualifiedName].factoryAction(iField)
	} else if p.ensuredFields[qualifiedName].fieldAction != nil {
		return p.ensuredFields[qualifiedName].fieldAction
	}

	return p.generators[k]
}

func (p *Plan) generateFieldValue(k reflect.Kind, generator GenType, iField reflect.Value, qualifiedName string) reflect.Value {
	if isPrimativeKind(k) {
		val := generator()

		return reflect.ValueOf(val)
	}

	// Complex Types
	switch k {
	case reflect.Slice:
		return p.updateSliceFieldValue(generator, iField)

	case reflect.Struct:
		return p.updateStructFieldValue(generator, iField, qualifiedName)

	case reflect.Ptr:
		ptrType := iField.Type().Elem() // The pointer's type
		ptrK := ptrType.Kind()
		generator := p.getValueGenerator(ptrK, iField, qualifiedName)

		newMockPtr := reflect.New(ptrType) // Make an instance based on the pointer type
		newElm := newMockPtr.Elem()

		val := p.generateFieldValue(ptrK, generator, newElm, qualifiedName)

		vp := reflect.New(val.Type())
		vp.Elem().Set(reflect.ValueOf(val.Interface()))

		return vp

	default:
		fmt.Printf("[updateFieldValue] (Unknow type) %v \n", iField.Type().Name())
	}
	panic("[updateFieldValue] Unsupported type")
}

func (p *Plan) updateSliceFieldValue(generator GenType, iField reflect.Value) reflect.Value {
	if generator == nil {
		factoryAction := makeFactoryAction(Tap())
		generator = factoryAction(iField)
	}

	factorySlice := generator()
	num := len(factorySlice.([]interface{}))

	newSlice := reflect.MakeSlice(iField.Type(), num, num)
	unboxedSlized := reflect.ValueOf(factorySlice)

	if iField.Type().Elem().Kind() == reflect.Ptr {
		for i := 0; i < unboxedSlized.Len(); i++ {
			obj := unboxedSlized.Index(i).Elem()
			newSlice.Index(i).Set(toPtr(obj))
		}
		return newSlice
	}

	for i := 0; i < unboxedSlized.Len(); i++ {
		newSlice.Index(i).Set(unboxedSlized.Index(i).Elem())
	}

	return newSlice
}

// toPtr converts obj to *obj
func toPtr(obj reflect.Value) reflect.Value {
	vp := reflect.New(reflect.TypeOf(obj.Interface()))
	vp.Elem().Set(reflect.ValueOf(obj.Interface()))

	return vp
}

func (p *Plan) updateStructFieldValue(generator GenType, iField reflect.Value, qualifiedName string) reflect.Value {
	if generator != nil {
		val := generator()
		return reflect.ValueOf(val)
	}

	iField.Interface()
	mock := Mock(iField.Interface())

	mock.plan.parentName = qualifiedName
	mock.plan.CopyParentRequiredFields(p)

	results := mock.Execute()

	return reflect.ValueOf(results[0])
}

func makeFactoryAction(fac *Factory) FactoryActionType {
	return func(iField reflect.Value) func() interface{} {
		// Extract the type from the slice and assign an instance to the rootType.
		// I want to go from []examples.Person to example.Person
		fac.rootType = reflect.New(reflect.TypeOf(iField.Interface()).Elem()).Elem().Interface() // Overwrite the factory with the public field type

		return func() interface{} { // The actual generator
			result := fac.Execute()
			return result
		}
	}
}

func distinctFileName(parentFieldName string, fieldName string) string {
	if parentFieldName == "" {
		return fieldName
	}

	return fmt.Sprintf("%s.%s", parentFieldName, fieldName)
}

func isPrimativeKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		return true
	}

	return false
}
