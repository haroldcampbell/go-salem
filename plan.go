package salem

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type RunType uint
type FactoryActionType func(reflect.Value, string) GenType
type SequenceActionType func(int) GenType // func(itemIndex int) GenType

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
	fieldAction         GenType
	factoryAction       FactoryActionType
	fieldSequenceAction SequenceActionType
}

type Plan struct {
	omittedFields     map[string]bool            // ignore these fields
	ensuredFields     map[string]fieldSetter     // fields set via ensure
	constrainedFields map[string]FieldConstraint // fields constraints
	generators        map[reflect.Kind]func() interface{}

	run *PlanRun

	evalItemCountAction func()
	parentName          string

	maxConstraintRetryAttempts int
}

func NewPlan() *Plan {
	p := &Plan{}

	p.omittedFields = make(map[string]bool)
	p.ensuredFields = make(map[string]fieldSetter)
	p.constrainedFields = make(map[string]FieldConstraint)
	p.generators = make(map[reflect.Kind]func() interface{})
	p.maxConstraintRetryAttempts = SuggestedConstraintRetryAttempts

	p.initDefaultGenerators()

	return p
}
func (p *Plan) GetPlanRun() *PlanRun {
	return p.run
}
func (p *Plan) SetMaxConstraintsRetryAttempts(maxRetries int) {
	p.maxConstraintRetryAttempts = maxRetries
}

func (p *Plan) SetItemCountHandler(handler func()) {
	p.evalItemCountAction = handler
}

func (p *Plan) OmitField(fieldName string) {
	p.omittedFields[fieldName] = true

}

func (p *Plan) AddFieldAction(fieldName string, fieldActionHanlder GenType) {
	p.ensuredFields[fieldName] = fieldSetter{fieldAction: fieldActionHanlder}
}

func (p *Plan) EnsuredFieldValue(fieldName string, sharedValue interface{}) {
	p.AddFieldAction(fieldName, func() interface{} {
		return sharedValue
	})
}

func (p *Plan) EnsuredFieldValueConstraint(fieldName string, constraint FieldConstraint) {
	p.constrainedFields[fieldName] = constraint
}

func (p *Plan) EnsuredFactoryFieldValue(fieldName string, sharedValue interface{}) {
	setter := fieldSetter{
		factoryAction: makeFactoryAction(sharedValue.(*Factory), p),
	}

	p.ensuredFields[fieldName] = setter
}

// EnsureSequence returns a seq item based on the item index
func (p *Plan) EnsureSequence(fieldName string, seq []interface{}) {
	seqHandler := func(itemIndex int) GenType {
		action := func() interface{} {
			var val interface{}

			if itemIndex < len(seq) {
				val = seq[itemIndex]
			}

			result := val

			return result
		}

		return action
	}

	p.ensuredFields[fieldName] = fieldSetter{fieldSequenceAction: seqHandler}
}

// EnsureSequenceAcross returns the seq[] items based on the sequence index.
// The  sequenceIndex is iindependent of the item index.
func (p *Plan) EnsureSequenceAcross(fieldName string, seq []interface{}) {
	var sequenceIndex int
	seqHandler := func(itemIndex int) GenType {

		action := func() interface{} {
			var val interface{}

			if sequenceIndex < len(seq) {
				val = seq[sequenceIndex]
			}

			result := val
			sequenceIndex += 1

			return result
		}

		return action
	}

	p.ensuredFields[fieldName] = fieldSetter{fieldSequenceAction: seqHandler}
}
func (p *Plan) SetRunCount(runType RunType, n int) {
	p.run = &PlanRun{
		RunType: runType,
		Count:   n,
	}
}

func (p *Plan) CopyParentConstraints(pp *Plan) {
	for k, v := range pp.ensuredFields {
		p.ensuredFields[k] = v
	}

	for k, v := range pp.omittedFields {
		p.omittedFields[k] = v
	}
}

func (p *Plan) Run(f *Factory) []interface{} {
	rand.Seed(time.Now().UnixNano())
	p.evalItemCountAction()

	items := make([]interface{}, 0, p.run.Count)

	for itemIndex := 0; itemIndex < p.run.Count; itemIndex++ {
		mockType := reflect.TypeOf(f.rootType)

		items = append(items, p.generateRandomMock(mockType, itemIndex))
	}

	return items
}

func (p *Plan) generateRandomMock(mockType reflect.Type, itemIndex int) interface{} {
	rand.Seed(time.Now().UnixNano())

	// Create an mock instance of the struct
	newMockPtr := reflect.New(mockType)
	newElm := newMockPtr.Elem()

	if mockType.Kind() == reflect.Ptr {
		ptrType := mockType.Elem()

		newMockPtr = reflect.New(ptrType) // Make an instance based on the pointer type
		newElm = newMockPtr.Elem()

		mockType = newElm.Type()
	}

	if isPrimativeKind(mockType.Kind()) {
		generator := p.generators[mockType.Kind()]
		val := generator()
		return val
	}

	for i := 0; i < mockType.NumField(); i++ {
		structField := mockType.Field(i) // Get field in the struct
		k := structField.Type.Kind()
		fieldName := mockType.Field(i).Name

		iField := newElm.Field(i) // Get related instance field in the mock instance
		if !iField.CanSet() {
			continue // Skip private instance fields
		}

		qualifiedName := distinctFileName(p.parentName, fieldName)
		if p.omittedFields[qualifiedName] == true {
			continue // Skip omitted fields
		}

		val := p.generateValue(k, iField, itemIndex, qualifiedName)

		if !val.IsValid() {
			// Uncomment to make EnsureSequence() set the values that fall outside of the sequence
			// if isPrimativeKind(k) {
			// gen := p.generators[k]
			// val = reflect.ValueOf(gen())
			// } else {
			continue
			// }
		}
		iField.Set(val)
	}
	return newMockPtr.Elem().Interface()
}
func (p *Plan) generateValue(k reflect.Kind, iField reflect.Value, itemIndex int, qualifiedName string) reflect.Value {
	constraint := p.constrainedFields[qualifiedName]
	generator := p.getValueGenerator(k, iField, itemIndex, qualifiedName)

	if constraint == nil {
		return p.generateFieldValue(k, generator, iField, itemIndex, qualifiedName)
	}

	isValueFromEnsureAction := p.ensuredFields[qualifiedName].factoryAction != nil || p.ensuredFields[qualifiedName].fieldAction != nil

	if isValueFromEnsureAction {
		val := p.generateFieldValue(k, generator, iField, itemIndex, qualifiedName)
		if !constraint.IsValid(val.Interface()) {
			panic(fmt.Sprintf("Constraint clashes with one of your Ensure methods. Invalid FieldConstraint for field '%v'. Constraint: %#v.", qualifiedName, constraint))
		}
		return val
	}

	var val reflect.Value
	var attempt = 0

	for {
		val = p.generateFieldValue(k, generator, iField, itemIndex, qualifiedName)
		attempt += 1

		if constraint.IsValid(val.Interface()) {
			break
		}

		if attempt > p.maxConstraintRetryAttempts {
			panic(fmt.Sprintf("Unable to meet constraint %v after '%v' tries", constraint, p.maxConstraintRetryAttempts))
		}
	}

	return val
}

func (p *Plan) getValueGenerator(k reflect.Kind, iField reflect.Value, itemIndex int, qualifiedName string) GenType {
	if p.ensuredFields[qualifiedName].factoryAction != nil {
		return p.ensuredFields[qualifiedName].factoryAction(iField, qualifiedName)
	} else if p.ensuredFields[qualifiedName].fieldSequenceAction != nil {
		return p.ensuredFields[qualifiedName].fieldSequenceAction(itemIndex)
	} else if p.ensuredFields[qualifiedName].fieldAction != nil {
		return p.ensuredFields[qualifiedName].fieldAction
	}

	return p.generators[k]
}

func (p *Plan) generateFieldValue(k reflect.Kind, generator GenType, iField reflect.Value, itemIndex int, qualifiedName string) reflect.Value {
	if isPrimativeKind(k) {
		val := generator()
		return reflect.ValueOf(val)
	}

	// Complex Types
	switch k {
	case reflect.Map:
		return p.updateMapFieldValue(generator, iField, qualifiedName)

	case reflect.Slice:
		return p.updateSliceFieldValue(generator, iField, qualifiedName)

	case reflect.Struct:
		return p.updateStructFieldValue(generator, iField, qualifiedName)

	case reflect.Ptr:
		ptrType := iField.Type().Elem() // The pointer's type
		ptrK := ptrType.Kind()
		generator := p.getValueGenerator(ptrK, iField, itemIndex, qualifiedName)

		newMockPtr := reflect.New(ptrType) // Make an instance based on the pointer type
		newElm := newMockPtr.Elem()

		val := p.generateFieldValue(ptrK, generator, newElm, itemIndex, qualifiedName)

		vp := reflect.New(val.Type())
		vp.Elem().Set(reflect.ValueOf(val.Interface()))

		return vp

	case reflect.Interface:
		return reflect.ValueOf(nil)

	default:
		fmt.Printf("[updateFieldValue] (Unknow type) %v \n", iField.Type().Name())
	}
	panic(fmt.Sprintf("[updateFieldValue] Unsupported type: %#v kind:%v", iField, k))
}

func (p *Plan) updateMapFieldValue(generator GenType, iField reflect.Value, qualifiedName string) reflect.Value {
	// TODO: Implement support for maps
	newMap := reflect.MakeMap(iField.Type())
	mapKeyType := iField.Type().Key()
	mapValueType := iField.Type().Elem()

	var keyGenerator GenType
	var valGenerator GenType
	if isPrimativeKind(mapKeyType.Kind()) {
		keyGenerator = p.GetKindGenerator(mapKeyType.Kind())
	} else {
		panic("Don't know how to generate key generator")
	}
	key := keyGenerator()

	var val interface{}
	if isPrimativeKind(mapValueType.Kind()) {
		valGenerator = p.GetKindGenerator(mapValueType.Kind())
		val = valGenerator()
	} else {
		// Create values that aren't primative types
		val = p.generateRandomMock(mapValueType, 0)
		if val == nil {
			panic(fmt.Sprintf("Didn't find a value generator for: %v type", mapValueType))
		}
	}

	newMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))

	return newMap
}

func (p *Plan) updateSliceFieldValue(generator GenType, iField reflect.Value, qualifiedName string) reflect.Value {
	if generator == nil {
		factoryAction := makeFactoryAction(Tap(), p)
		generator = factoryAction(iField, qualifiedName)
	}

	factorySlice := generator()
	unboxedSlized := reflect.ValueOf(factorySlice)
	num := unboxedSlized.Len()
	newSlice := reflect.MakeSlice(iField.Type(), num, num)

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
	mock.plan.CopyParentConstraints(p)

	results := mock.Execute()

	return reflect.ValueOf(results[0])
}

func makeFactoryAction(fac *Factory, currentPlan *Plan) FactoryActionType {
	return func(iField reflect.Value, qualifiedName string) func() interface{} {
		// Extract the type from the slice and assign an instance to the rootType.
		// I want to go from []examples.Person to example.Person
		fac.rootType = reflect.New(reflect.TypeOf(iField.Interface()).Elem()).Elem().Interface() // Overwrite the factory with the public field type

		return func() interface{} { // The actual generator
			fac.plan.parentName = qualifiedName
			fac.plan.CopyParentConstraints(currentPlan)

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
