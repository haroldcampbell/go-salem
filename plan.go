package salem

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type RunType uint
type deferredExecuteType func(reflect.Value) GenType

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
	fptr            GenType
	deferredExecute deferredExecuteType
}
type Plan struct {
	ensuredFields map[string]fieldSetter // fields set via ensure
	generators    map[reflect.Kind]func() interface{}

	run PlanRun

	deferredRun func()
	parentName  string
}

func NewPlan() *Plan {
	p := &Plan{}

	p.ensuredFields = make(map[string]fieldSetter)
	p.generators = make(map[reflect.Kind]func() interface{})

	p.initDefaultGenerators()

	return p
}

func (p *Plan) EnsuredFieldValue(fieldName string, sharedValue interface{}) {
	setter := fieldSetter{
		fptr: func() interface{} {
			return sharedValue
		},
	}

	p.ensuredFields[fieldName] = setter
}

func (p *Plan) makeDeferredExecute(fac *Factory) deferredExecuteType {
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

func (p *Plan) EnsuredDeferredFieldValue(fieldName string, sharedValue interface{}) {
	setter := fieldSetter{
		deferredExecute: p.makeDeferredExecute(sharedValue.(*Factory)),
	}

	p.ensuredFields[fieldName] = setter
}

func (p *Plan) SetRunCount(runType RunType, n int) {
	p.run = PlanRun{
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
	p.deferredRun()

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

		p.updateFieldValue(k, generator, iField, qualifiedName)
	}

	return newMockPtr.Elem().Interface()
}

func (p *Plan) getValueGenerator(k reflect.Kind, iField reflect.Value, qualifiedName string) GenType {
	if p.ensuredFields[qualifiedName].deferredExecute != nil {
		return p.ensuredFields[qualifiedName].deferredExecute(iField)
	} else if p.ensuredFields[qualifiedName].fptr != nil {
		return p.ensuredFields[qualifiedName].fptr
	}

	return p.generators[k]
}

func (p *Plan) updateFieldValue(k reflect.Kind, generator GenType, iField reflect.Value, qualifiedName string) {
	if isPrimativeKind(k) {
		val := generator()
		iField.Set(reflect.ValueOf(val))

		return
	}

	// Complex Types
	switch k {
	case reflect.Slice:
		p.updateSliceFieldValue(generator, iField)

	case reflect.Struct:
		p.updateStructFieldValue(generator, iField, qualifiedName)

	default:
		fmt.Printf("[updateFieldValue] (Unknow type) %v \n", iField.Type().Name())
	}
}

func (p *Plan) updateSliceFieldValue(generator GenType, iField reflect.Value) {
	if generator == nil {
		deferredExecute := p.makeDeferredExecute(Tap())
		generator = deferredExecute(iField)
	}

	factorySlice := generator()
	num := len(factorySlice.([]interface{}))

	newSlice := reflect.MakeSlice(iField.Type(), num, num)
	unboxedSlized := reflect.ValueOf(factorySlice)

	for i := 0; i < unboxedSlized.Len(); i++ {
		newSlice.Index(i).Set(unboxedSlized.Index(i).Elem())
	}

	iField.Set(newSlice)
}

func (p *Plan) updateStructFieldValue(generator GenType, iField reflect.Value, qualifiedName string) {
	if generator != nil {
		val := generator()
		iField.Set(reflect.ValueOf(val))
	} else {
		iField.Interface()
		mock := Mock(iField.Interface())

		mock.plan.parentName = qualifiedName
		mock.plan.CopyParentRequiredFields(p)

		results := mock.Execute()

		iField.Set(reflect.ValueOf(results[0]))
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
