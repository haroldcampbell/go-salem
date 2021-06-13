// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

type money struct {
	Value float32
}

func Test_FactoryEnsure(t *testing.T) {
	test_simple_ensure(t)
	test_nested_ensure(t)
	test_slice_ensure(t)
	test_slice_pointer_ensure(t)
	test_sequence_across_ensure(t)
	test_sequence_ensure(t)
	test_sequence_across_splice_ensure(t)
}

func test_simple_ensure(t *testing.T) {
	requiredValue := "happiness"

	f := salem.Mock(simple{})
	f.Ensure("num", 101)
	f.Ensure("PublicField", requiredValue)

	results := f.Execute()
	actualMock := results[0].(simple)

	assert.Equal(t, 1, len(results), "expect mock to be created")
	assert.Equal(t, 0, actualMock.num, "should not set private fields")
	assert.Equal(t, requiredValue, actualMock.PublicField, "should set public fields required value")
}

func test_nested_ensure(t *testing.T) {
	type basicChild struct {
		ParentName string
		Name       string
		Age        int

		Attends school
	}

	type basicParent struct {
		Child basicChild
	}

	requiredValue := "Goose BrightSpark"

	f := salem.Mock(basicParent{})
	f.Ensure("Child.ParentName", requiredValue)
	f.Ensure("Child.Attends.NextOfKin", requiredValue)

	results := f.Execute()
	actualMock := results[0].(basicParent)

	assert.Equal(t, requiredValue, actualMock.Child.ParentName, "should set nested public fields required value")
	assert.Equal(t, requiredValue, actualMock.Child.Attends.NextOfKin, "should set nested public fields required value")
}

func test_slice_pointer_ensure(t *testing.T) {
	type vault struct {
		Notes []*money
	}

	tap := salem.Tap().WithExactItems(5)

	f := salem.Mock(vault{})
	f.Ensure("Notes", tap)

	results := f.Execute()
	actualMocks := results[0].(vault).Notes

	assert.Equal(t, 5, len(actualMocks), "should set nested slice of pointer")
}

func test_slice_ensure(t *testing.T) {
	type wallet struct {
		Notes []money
	}

	tap := salem.Tap().WithExactItems(5)

	f := salem.Mock(wallet{})
	f.Ensure("Notes", tap)

	results := f.Execute()
	actualMocks := results[0].(wallet).Notes

	assert.Equal(t, 5, len(actualMocks), "should set nested slice")
}

func test_sequence_across_ensure(t *testing.T) {
	requiredValue := "happiness"

	f := salem.Mock(simple{})
	f.EnsureSequenceAcross("PublicField", "a", requiredValue, "c")
	f.WithExactItems(5)

	results := f.Execute()
	publicFields := make([]string, len(results))
	for i, v := range results {
		publicFields[i] = v.(simple).PublicField
	}

	expected := []string{"a", requiredValue, "c", "", ""}
	actual := []string{publicFields[0], publicFields[1], publicFields[2], publicFields[3], publicFields[4]}
	assert.Equal(t, requiredValue, publicFields[1], "should be set to sequence value")
	assert.Equal(t, expected, actual, "should set sequence values")
}

func test_sequence_ensure(t *testing.T) {
	type people struct {
		Name string
		Cash money
	}
	type market struct {
		Venue   string
		Farmers []people
	}

	seq := []string{"Bill", "Ruth", "Mary", "Sally"}
	tap := salem.Tap().
		EnsureSequence("Farmers.Name", seq[0], seq[1], seq[2], seq[3]).
		WithExactItems(2)

	f := salem.Mock(market{}).
		EnsureSequence("Venue", "Park Town", "Rose Hill", "Pretoria").
		Ensure("Farmers", tap).
		WithExactItems(2)

	results := f.Execute()

	m1 := results[0].(market)
	m2 := results[1].(market)

	assert.Equal(t, "Park Town", m1.Venue, "expect correct value from item index")
	assert.Equal(t, "Rose Hill", m2.Venue, "expect correct value from item index")

	s1 := m1.Farmers[0]
	s2 := m1.Farmers[1]
	s3 := m2.Farmers[0]
	s4 := m2.Farmers[1]

	assert.Equal(t, seq[0], s1.Name, "expect correct value from sequence using nested item index")
	assert.Equal(t, seq[1], s2.Name, "expect correct value from sequence using nested item index")

	// Sequence should restart for different list item index
	assert.Equal(t, seq[0], s3.Name, "expect sequence to restart")
	assert.Equal(t, seq[1], s4.Name, "expect sequence to restart")
}

func test_sequence_across_splice_ensure(t *testing.T) {
	type people struct {
		Name string
		Cash money
	}
	type market struct {
		Venue   string
		Farmers []people
	}

	seq := []string{"Bill", "Ruth", "Mary", "Sally"}
	tap := salem.Tap().
		EnsureSequenceAcross("Farmers.Name", seq[0], seq[1], seq[2], seq[3]).
		WithExactItems(2)

	f := salem.Mock(market{}).
		Ensure("Farmers", tap).
		WithExactItems(2)

	results := f.Execute()

	m1 := results[0].(market)
	m2 := results[1].(market)

	s1 := m1.Farmers[0]
	s2 := m1.Farmers[1]

	s3 := m2.Farmers[0]
	s4 := m2.Farmers[1]

	assert.Equal(t, seq[0], s1.Name)
	assert.Equal(t, seq[1], s2.Name)

	// Sequence should continue based on sequence items' index
	assert.Equal(t, seq[2], s3.Name, "expect sequence to continue based on sequence index")
	assert.Equal(t, seq[3], s4.Name)
}
