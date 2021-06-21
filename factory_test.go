// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

type empty struct{}
type simple struct {
	num         int
	PublicField string
}
type school struct {
	NextOfKin string
	Name      string
}

func Test_Factory(t *testing.T) {
	f := salem.Mock(simple{})
	results := f.Execute()

	assert.Equal(t, 1, len(results), "expect Execute() to create array")
	assert.IsType(t, []interface{}{}, results, "expect Execute() to return []interface{}")

	targetReuslt := f.ExecuteToType().([]simple)
	assert.Equal(t, 1, len(targetReuslt), "expect ExecuteToType() to create array")
	assert.IsType(t, []simple{}, targetReuslt, "expect ExecuteToType() to return interface{} that is a slice of base type")
}

func Test_InterfaceEdgeCase(t *testing.T) {
	type canvasElementModel struct {
		Data interface{}
	}

	target := salem.Mock(canvasElementModel{}).
		ExecuteToType().([]canvasElementModel)[0]

	assert.Equal(t, nil, target.Data, "expect interface to be nil by default")
}

func Test_InterfaceEdgeCase_WithEnsure(t *testing.T) {
	type canvasElementModel struct {
		Data interface{}
	}

	data1 := "data attributes 101"
	target1 := salem.Mock(canvasElementModel{}).
		Ensure("Data", data1).
		ExecuteToType().([]canvasElementModel)[0]

	assert.Equal(t, data1, target1.Data, "expect interface types to be set with Ensure(string)")

	data2 := 10
	target2 := salem.Mock(canvasElementModel{}).
		Ensure("Data", data2).
		ExecuteToType().([]canvasElementModel)[0]

	assert.Equal(t, data2, target2.Data, "expect interface types to be set with Ensure(int)")
}
func Test_FactoryOmit(t *testing.T) {
	f := salem.Mock(school{})

	results := f.Execute()
	actualMock := results[0].(school)
	assert.NotEmpty(t, actualMock.NextOfKin, "expect field to be set")

	f.Omit("NextOfKin")
	results = f.Execute()
	actualMock = results[0].(school)

	assert.Empty(t, actualMock.NextOfKin, "expect field to be empty")
}

func Test_FactoryWithItems(t *testing.T) {
	test_with__items(t)
	test_with_exact_items(t)
	test_with_min_items(t)
	test_with_max_items(t)
}

func test_with__items(t *testing.T) {
	f := salem.Mock(empty{})
	assert.Nil(t, f.GetPlan().GetPlanRun(), "expect plan.run to be nil before Execute")
}
func test_with_exact_items(t *testing.T) {
	f := salem.Mock(empty{})

	f.WithExactItems(2)
	assert.Nil(t, f.GetPlan().GetPlanRun(), "expect WithExactItems to not change plan.run")

	f.Execute()
	assert.Equal(t, salem.ExactRun, f.GetPlan().GetPlanRun().RunType, "expect ExactRun from WithExactItems")
	assert.Equal(t, 2, f.GetPlan().GetPlanRun().Count, "expect correct run count from WithExactItems")
}

func test_with_min_items(t *testing.T) {
	f := salem.Mock(empty{})

	f.WithMinItems(2)
	assert.Nil(t, f.GetPlan().GetPlanRun(), "expect WithMinItems to not change plan.run")

	f.Execute()
	assert.Equal(t, salem.MinRun, f.GetPlan().GetPlanRun().RunType, "expect MinRun from WithMinItems")
	assert.GreaterOrEqual(t, f.GetPlan().GetPlanRun().Count, 2, "expect Count to be 2 or more")
}

func test_with_max_items(t *testing.T) {
	f := salem.Mock(empty{})

	f.WithMaxItems(2) // No need to test the run.Count as it will be based on a random value of n
	assert.Nil(t, f.GetPlan().GetPlanRun(), "expect WithMaxItems to not change plan.run")

	f.Execute()
	assert.Equal(t, salem.MaxRun, f.GetPlan().GetPlanRun().RunType, "expect MaxRun from WithExactItems")
}
