// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem

import (
	"fmt"
	"reflect"
)

type processorType func(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value

func (p *Plan) initKindProcessors() {

	p.kindProcessors = make(map[reflect.Kind]processorType)

	p.kindProcessors[reflect.Map] = onMap(p)
	p.kindProcessors[reflect.Ptr] = onPtr(p)
	p.kindProcessors[reflect.Slice] = onSlice(p)
	p.kindProcessors[reflect.Struct] = onStruct(p)
	p.kindProcessors[reflect.Interface] = onInterface(p)
}

func onMap(p *Plan) processorType {
	return func(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {
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
		keyGenerator := createMapKeyGenerator(p.GetKindGenerator(mapKeyType.Kind()), fieldSequenceKeyAction)

		// Dynamically create the valueGenerator
		valueGenerator := p.createMapValueGenerator(fieldType.Elem(), qualifiedName)

		for mapItemIndex := 0; mapItemIndex < mapItemCount; mapItemIndex++ {
			key := keyGenerator(mapItemIndex)
			val := valueGenerator(mapItemIndex)

			newMap.SetMapIndex(reflect.ValueOf(key), val)
		}

		return newMap
	}
}

func onSlice(p *Plan) processorType {
	return func(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {
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
}

func onStruct(p *Plan) processorType {
	return func(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {
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
}

func onPtr(p *Plan) processorType {
	return func(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {

		ptrType := fieldType.Elem() // The pointer's type
		generator = p.getValueGenerator(ptrType, itemIndex, qualifiedName)

		newMockPtr := reflect.New(ptrType) // Make an instance based on the pointer type
		newElm := newMockPtr.Elem()

		val := p.generateFieldValue(generator, newElm.Type(), itemIndex, qualifiedName)

		vp := reflect.New(val.Type())
		vp.Elem().Set(reflect.ValueOf(val.Interface()))

		return vp
	}
}

func onInterface(p *Plan) processorType {
	return func(generator GenType, fieldType reflect.Type, itemIndex int, qualifiedName string) reflect.Value {
		val := generator()
		return reflect.ValueOf(val)
	}
}

// createMapKeyGenerator is used to dynamically create the keyGenerator so we
// don't need to call the 'if' inside of the for loop
func createMapKeyGenerator(generator func() interface{}, fieldSequenceKeyAction SequenceActionType) func(param int) interface{} {
	if fieldSequenceKeyAction != nil {
		return func(index int) interface{} {
			return fieldSequenceKeyAction(index)()
		}
	}
	// If we get here we are guaranteed that the generator is a
	// isPrimitiveKind(...) and comes from p.GetKindGenerator(...)
	return func(_ int) interface{} {
		return generator()
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
