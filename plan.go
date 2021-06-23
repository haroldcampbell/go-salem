// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type RunType uint
type FactoryActionType func(reflect.Type, string) GenType
type SequenceActionType func(int) GenType // func(itemIndex int) GenType
type fieldHandlerType func(int) interface{}

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

// mapSetter used to hold the generators for keys of values for a map field
type mapSetter struct {
	fieldSequenceKeyAction   SequenceActionType
	fieldSequenceValueAction SequenceActionType
}

type Plan struct {
	omittedFields     map[string]bool            // ignore these fields
	ensuredFields     map[string]fieldSetter     // fields set via ensure
	constrainedFields map[string]FieldConstraint // fields constraints
	generators        map[reflect.Kind]func() interface{}

	ensuredMapFields map[string]mapSetter // map key fields set via ensure

	run                 *PlanRun
	evalItemCountAction func()

	mapPlanRun             map[string]*PlanRun // plan runs for different Maps
	evalMapItemCountAction map[string]func()

	fieldHandlers map[string]fieldHandlerType // FieldName -> Handler

	parentName string

	maxConstraintRetryAttempts int
}

func NewPlan() *Plan {
	p := &Plan{}

	p.omittedFields = make(map[string]bool)
	p.ensuredFields = make(map[string]fieldSetter)
	p.constrainedFields = make(map[string]FieldConstraint)
	p.generators = make(map[reflect.Kind]func() interface{})
	p.maxConstraintRetryAttempts = SuggestedConstraintRetryAttempts

	p.ensuredMapFields = make(map[string]mapSetter)
	p.mapPlanRun = make(map[string]*PlanRun)
	p.evalMapItemCountAction = make(map[string]func())

	p.fieldHandlers = make(map[string]fieldHandlerType)

	p.initDefaultGenerators()

	return p
}
func (p *Plan) GetPlanRun() *PlanRun {
	return p.run
}
func (p *Plan) GetMapPlanRun(fieldName string) *PlanRun {
	return p.mapPlanRun[fieldName]
}

func (p *Plan) SetMaxConstraintsRetryAttempts(maxRetries int) {
	p.maxConstraintRetryAttempts = maxRetries
}

func (p *Plan) SetItemCountHandler(handler func()) {
	p.evalItemCountAction = handler
}
func (p *Plan) SetMapItemCountHandler(fieldName string, handler func()) {
	p.evalMapItemCountAction[fieldName] = handler
}
func (p *Plan) OmitField(fieldName string) {
	p.omittedFields[fieldName] = true

}

func (p *Plan) EnsuredFieldValue(fieldName string, sharedValue interface{}) {
	setter := fieldSetter{
		fieldAction: func() interface{} {
			return sharedValue
		},
	}

	p.ensuredFields[fieldName] = setter
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
	p.ensuredFields[fieldName] = fieldSetter{fieldSequenceAction: sequenceCallbackCreator(seq)}
}

// EnsureSequenceAcross returns the seq[] items based on the sequence index.
// The  sequenceIndex is independent of the item index.
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

func (p *Plan) EnsureMapKeySequence(fieldName string, seq []interface{}) {
	var setter mapSetter
	if _, ok := p.ensuredMapFields[fieldName]; ok {
		setter = p.ensuredMapFields[fieldName]
		setter.fieldSequenceKeyAction = sequenceCallbackCreator(seq)
	} else {
		setter = mapSetter{fieldSequenceKeyAction: sequenceCallbackCreator(seq)}
	}
	p.ensuredMapFields[fieldName] = setter
}

func (p *Plan) EnsureMapValueSequence(fieldName string, seq []interface{}) {
	var setter mapSetter
	if _, ok := p.ensuredMapFields[fieldName]; ok {
		setter = p.ensuredMapFields[fieldName]
		setter.fieldSequenceValueAction = sequenceCallbackCreator(seq)
	} else {
		setter = mapSetter{fieldSequenceValueAction: sequenceCallbackCreator(seq)}
	}
	p.ensuredMapFields[fieldName] = setter
}

func sequenceCallbackCreator(seq []interface{}) func(itemIndex int) GenType {
	return func(itemIndex int) GenType {
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
}

// EnsureMayKeySequence returns a seq item based on the item index
func (p *Plan) EnsureMayKeySequence(fieldName string, seq []interface{}) {
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

func (p *Plan) SetRunCount(runType RunType, n int) {
	p.run = &PlanRun{
		RunType: runType,
		Count:   n,
	}
}

func (p *Plan) SetMapRunCount(fieldName string, runType RunType, n int) {
	p.mapPlanRun[fieldName] = &PlanRun{
		RunType: runType,
		Count:   n,
	}
}

func (p *Plan) AddFieldHandler(fieldName string, handler fieldHandlerType) {
	if p.fieldHandlers[fieldName] != nil {
		panic(fmt.Sprintf("There is already a fieldHandler assigned to to the `%v` field. Remove the fieldhandler first.", fieldName))
	}

	p.fieldHandlers[fieldName] = handler
}

func (p *Plan) RemoveFieldHandler(fieldName string) {
	p.fieldHandlers[fieldName] = nil
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

	if isPrimitiveKind(mockType) {
		generator := p.GetKindGenerator(mockType.Kind())
		val := generator()
		return val
	}

	for i := 0; i < mockType.NumField(); i++ {
		fieldName := mockType.Field(i).Name

		iField := newElm.Field(i) // Get related instance field in the mock instance
		if !iField.CanSet() {
			continue // Skip private instance fields
		}

		qualifiedName := distinctFileName(p.parentName, fieldName)
		if p.omittedFields[qualifiedName] == true {
			continue // Skip omitted fields
		}

		val := p.generateValue(iField.Type(), itemIndex, qualifiedName)

		if !val.IsValid() {
			continue
		}
		iField.Set(val)
	}
	return newMockPtr.Elem().Interface()
}

func (p *Plan) generateValue(fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {
	if p.fieldHandlers[qualifiedName] != nil {
		generator := p.fieldHandlers[qualifiedName]
		val := generator(itemIndex)
		return reflect.ValueOf(val)
	}

	generator := p.getValueGenerator(fieldType, itemIndex, qualifiedName)

	constraint := p.constrainedFields[qualifiedName]
	if constraint == nil { // Generate field value and exit since no constraint
		return p.generateFieldValue(generator, fieldType, itemIndex, qualifiedName)
	}

	isValueFromEnsureAction := p.ensuredFields[qualifiedName].factoryAction != nil || p.ensuredFields[qualifiedName].fieldAction != nil
	if isValueFromEnsureAction { // Ensure Ensured field meets constraint
		val := p.generateFieldValue(generator, fieldType, itemIndex, qualifiedName)
		if !constraint.IsValid(val.Interface()) {
			panic(fmt.Sprintf("Constraint clashes with one of your Ensure methods. Invalid FieldConstraint for field '%v'. Constraint: %#v.", qualifiedName, constraint))
		}
		return val
	}

	var attempt = 0
	var val reflect.Value
	for { // Try till constraint is met or give up after maxConstraintRetryAttempts attemps
		val = p.generateFieldValue(generator, fieldType, itemIndex, qualifiedName)
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

func (p *Plan) getValueGenerator(fieldType reflect.Type, itemIndex int, qualifiedName string) GenType {
	if p.ensuredFields[qualifiedName].factoryAction != nil {
		return p.ensuredFields[qualifiedName].factoryAction(fieldType, qualifiedName)

	} else if p.ensuredFields[qualifiedName].fieldSequenceAction != nil {
		return p.ensuredFields[qualifiedName].fieldSequenceAction(itemIndex)

	} else if p.ensuredFields[qualifiedName].fieldAction != nil {
		return p.ensuredFields[qualifiedName].fieldAction

	}

	return p.GetKindGenerator(fieldType.Kind())
}

func (p *Plan) generateFieldValue(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {
	if isPrimitiveKind(fieldType) {
		val := generator()
		return reflect.ValueOf(val)
	}

	k := fieldType.Kind()

	// Complex Types
	switch k {
	case reflect.Map:
		return p.updateMapFieldValue(fieldType, qualifiedName)

	case reflect.Slice:
		return p.updateSliceFieldValue(generator, fieldType, qualifiedName)

	case reflect.Struct:
		return p.updateStructFieldValue(generator, fieldType, qualifiedName)

	case reflect.Ptr:
		ptrType := fieldType.Elem() // The pointer's type
		generator := p.getValueGenerator(ptrType, itemIndex, qualifiedName)

		newMockPtr := reflect.New(ptrType) // Make an instance based on the pointer type
		newElm := newMockPtr.Elem()

		val := p.generateFieldValue(generator, newElm.Type(), itemIndex, qualifiedName)

		vp := reflect.New(val.Type())
		vp.Elem().Set(reflect.ValueOf(val.Interface()))

		return vp

	case reflect.Interface:
		val := generator()
		return reflect.ValueOf(val)

	default:
		fmt.Printf("[updateFieldValue] (Unknow type) %v \n", fieldType.Name())
	}
	panic(fmt.Sprintf("[updateFieldValue] Unsupported type: %#v kind:%v", fieldType, k))
}

func (p *Plan) updateMapFieldValue(fieldType reflect.Type, qualifiedName string) reflect.Value {
	newMap := reflect.MakeMap(fieldType)
	mapKeyType := fieldType.Key()

	fieldSequenceKeyAction := p.ensuredMapFields[qualifiedName].fieldSequenceKeyAction
	if fieldSequenceKeyAction == nil && !isPrimitiveKind(mapKeyType) {
		// Can't be generate the field by fieldSequenceAction(...) or p.GetKindGenerator(...)
		panic(fmt.Sprintf("Don't know how to make the key-generator. Field: %v", qualifiedName))
	}

	var mapItemCount = 1
	if p.evalMapItemCountAction[qualifiedName] != nil {
		p.evalMapItemCountAction[qualifiedName]()
		mapItemCount = p.mapPlanRun[qualifiedName].Count
	}

	// Dynamically create the keyGenerator
	keyGenerator := p.createMapKeyGenerator(fieldSequenceKeyAction, mapKeyType)

	// Dynamically create the valueGenerator
	valueGenerator := p.createMapValueGenerator(fieldType.Elem(), qualifiedName)

	for mapItemIndex := 0; mapItemIndex < mapItemCount; mapItemIndex++ {
		key := keyGenerator(mapItemIndex)
		val := valueGenerator(mapItemIndex)

		newMap.SetMapIndex(reflect.ValueOf(key), val)
	}

	return newMap
}

// createMapKeyGenerator is used to dynamically create the keyGenerator so we
// don't need to call the 'if' inside of the for loop
func (p *Plan) createMapKeyGenerator(fieldSequenceKeyAction SequenceActionType, mapKeyType reflect.Type) func(param int) interface{} {
	if fieldSequenceKeyAction != nil {
		return func(index int) interface{} {
			return fieldSequenceKeyAction(index)()
		}
	}
	// If we get here we are guaranteed that the generator is a
	// isPrimitiveKind(...) and comes from p.GetKindGenerator(...)
	return func(_ int) interface{} {
		return p.GetKindGenerator(mapKeyType.Kind())()
	}
}

// createMapValueGenerator is used to dynamically create the valueGenerator for a map's value
func (p *Plan) createMapValueGenerator(mapValueType reflect.Type, qualifiedName string) func(param int) reflect.Value {
	fieldSequenceValueAction := p.ensuredMapFields[qualifiedName].fieldSequenceValueAction

	if fieldSequenceValueAction != nil {
		return func(index int) reflect.Value {
			result := fieldSequenceValueAction(index)()
			return reflect.ValueOf(result)
		}
	} else if isPrimitiveKind(mapValueType) {
		return func(_ int) reflect.Value {
			result := p.GetKindGenerator(mapValueType.Kind())()
			return reflect.ValueOf(result)
		}
	} else if isPrtPrimitiveKind(mapValueType) {
		return func(_ int) reflect.Value {
			// Get the primitive type the pointer points to then generate the primitive value
			result := p.GetKindGenerator(mapValueType.Elem().Kind())()

			// Convert value to a pointer
			vp := reflect.New(mapValueType.Elem())
			vp.Elem().Set(reflect.ValueOf(result))

			return vp // Return the pointer
		}
	}

	return func(_ int) reflect.Value {
		return p.generateFieldValue(nil, mapValueType, 0, qualifiedName)
	}
}

func (p *Plan) updateSliceFieldValue(generator GenType, fieldType reflect.Type, qualifiedName string) reflect.Value {
	if generator == nil {
		factoryAction := makeFactoryAction(Tap(), p)
		generator = factoryAction(fieldType, qualifiedName)
	}

	factorySlice := generator()
	unboxedSlized := reflect.ValueOf(factorySlice)
	num := unboxedSlized.Len()
	newSlice := reflect.MakeSlice(fieldType, num, num)

	if fieldType.Elem().Kind() == reflect.Ptr {
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

func (p *Plan) updateStructFieldValue(generator GenType, fieldType reflect.Type, qualifiedName string) reflect.Value {
	if generator != nil {
		val := generator()
		return reflect.ValueOf(val)
	}

	fieldPtr := reflect.New(fieldType) // Make an instance based on the pointer type
	fieldVal := fieldPtr.Elem()

	mock := Mock(fieldVal.Interface())

	mock.plan.parentName = qualifiedName
	mock.plan.CopyParentConstraints(p)

	results := mock.Execute()

	return reflect.ValueOf(results[0])
}

func makeFactoryAction(fac *Factory, currentPlan *Plan) FactoryActionType {
	return func(fieldType reflect.Type, qualifiedName string) func() interface{} {
		fieldPtr := reflect.New(fieldType) // Make an instance based on the pointer type
		fieldVal := fieldPtr.Elem()
		// Extract the type from the slice and assign an instance to the rootType.
		// I want to go from []examples.Person to example.Person
		fac.rootType = reflect.New(reflect.TypeOf(fieldVal.Interface()).Elem()).Elem().Interface() // Overwrite the factory with the public field type

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

func isPrimitiveKind(t reflect.Type) bool {
	k := t.Kind()

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

func isPrtPrimitiveKind(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		return false
	}

	return isPrimitiveKind(t.Elem())
}
