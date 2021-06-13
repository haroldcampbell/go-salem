// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem_test

import (
	"go-salem"

	"github.com/stretchr/testify/assert"
)

func (s *mapSuite) test_map_exact_item_count() {
	type human struct {
		Genes        map[string]string // gene[name] -> (DNA string sequence)
		VacationDays map[int]int       // month -> days off
	}

	f := salem.Mock(human{}).
		WithExactMapItems("Genes", 5)

	result := f.ExecuteToType().([]human)[0]

	t := s.T()
	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 5, len(result.Genes), "should item count equal to WithExactMapItems(...)")
	assert.Equal(t, 1, len(result.VacationDays), "should not be affected bu WithExactMapItems(...)")

	plan := f.GetPlan()
	planRun := plan.GetMapPlanRun("Genes")
	assert.Equal(t, salem.ExactRun, planRun.RunType, "should set current RunType to salem.ExactRun")
}

func (s *mapSuite) test_map_max_item_count() {
	type human struct {
		Genes        map[string]string // gene[name] -> (DNA string sequence)
		VacationDays map[int]int       // month -> days off
	}

	f := salem.Mock(human{}).
		WithMaxMapItems("Genes", 5)

	result := f.ExecuteToType().([]human)[0]

	t := s.T()

	// Sanity checks
	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 1, len(result.VacationDays), "should not be affected bu WithMaxMapItems(...)")

	plan := f.GetPlan()
	planRun := plan.GetMapPlanRun("Genes")

	assert.NotNil(t, planRun, "should have plan for field set with WithMaxMapItems")
	assert.Equal(t, salem.MaxRun, planRun.RunType, "should set current RunType to salem.MaxRun")
}

func (s *mapSuite) test_map_min_item_count() {
	type human struct {
		Genes        map[string]string // gene[name] -> (DNA string sequence)
		VacationDays map[int]int       // month -> days off
	}

	lowerBounds := 5

	f := salem.Mock(human{}).
		WithMinMapItems("Genes", lowerBounds).
		WithMinMapItems("VacationDays", lowerBounds, 5)

	result := f.ExecuteToType().([]human)[0]

	t := s.T()

	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.GreaterOrEqual(t, len(result.Genes), lowerBounds, "expect WithMinMapItems(...) to generate at least n items")
	assert.GreaterOrEqual(t, len(result.VacationDays), lowerBounds, "expect WithMinMapItems(...) to generate value greater or equal lower lowerBounds")

	plan := f.GetPlan()
	planRun := plan.GetMapPlanRun("Genes")

	assert.NotNil(t, planRun, "should have plan for field set with WithMaxMapItems")
	assert.Equal(t, salem.MinRun, planRun.RunType, "should set current RunType to salem.MinRun")
}
