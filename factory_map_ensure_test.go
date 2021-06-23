// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

type human struct {
	Genes map[string]string // gene[name] -> (DNA string sequence)
}

// func (s *mapSuite) Test_ensure_map_pointer() {
func Test_ensure_map_pointer(t *testing.T) {
	type human struct {
		Genes map[string]*string
	}

	result := salem.Mock(human{}).
		Ensure("Genes", map[string]*string{}). // <- The pointer makes it fail for some reason
		ExecuteToType().([]human)[0]

	// t := s.T()
	assert.NotEmpty(t, result.Genes, "should create string map with Ensure")
	assert.Equal(t, 1, len(result.Genes))
}

func (s *mapSuite) Test_ensure_map_keys() {
	expectedKeys := []interface{}{"A2M", "ABL1", "ADCY5", "AGPAT2", "AGTR1"}
	result := salem.Mock(human{}).
		EnsureMapKeySequence("Genes", expectedKeys...).
		WithExactMapItems("Genes", 5).
		WithExactMapItems("VacationDays", 5).
		ExecuteToType().([]human)[0]

	t := s.T()
	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 5, len(result.Genes))

	actual := []string{}
	for key := range result.Genes {
		actual = append(actual, key)
	}
	assert.ElementsMatch(t, expectedKeys, actual, "Should create keys from EnsureMayKeySequence(...)")
}

func (s *mapSuite) Test_ensure_map_values() {
	expectedValues := []interface{}{"alpha-2-macroglobulin", "ABL proto-oncogene 1", "adenylate cyclase 5", "1-acylglycerol-3-phosphate O-acyltransferase 2", "angiotensin II receptor, type 1"}
	result := salem.Mock(human{}).
		EnsureMapValueSequence("Genes", expectedValues...).
		WithExactMapItems("Genes", 5).
		ExecuteToType().([]human)[0]

	t := s.T()
	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 5, len(result.Genes))

	actual := []string{}
	for _, val := range result.Genes {
		actual = append(actual, val)
	}
	assert.ElementsMatch(t, expectedValues, actual, "Should create keys from EnsureMayKeySequence(...)")
}

func (s *mapSuite) Test_ensure_key_value() {
	expectedKeys := []interface{}{"A2M", "ABL1", "ADCY5", "AGPAT2", "AGTR1"}
	expectedValues := []interface{}{"alpha-2-macroglobulin", "ABL proto-oncogene 1", "adenylate cyclase 5", "1-acylglycerol-3-phosphate O-acyltransferase 2", "angiotensin II receptor, type 1"}

	result := salem.Mock(human{}).
		EnsureMapValueSequence("Genes", expectedValues...).
		EnsureMapKeySequence("Genes", expectedKeys...).
		WithExactMapItems("Genes", 5).
		ExecuteToType().([]human)[0]

	t := s.T()
	assert.Equal(t, 5, len(result.Genes))

	actualK := []string{}
	actualV := []string{}
	for key, val := range result.Genes {
		actualK = append(actualK, key)
		actualV = append(actualV, val)
	}

	assert.ElementsMatch(t, expectedKeys, actualK)
	assert.ElementsMatch(t, expectedValues, actualV)
}
