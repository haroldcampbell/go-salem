package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FactoryMap(t *testing.T) {
	test_primative_types_map(t)
	test_map_of_structs(t)
	test_map_of_slices(t)
}

func test_primative_types_map(t *testing.T) {
	type human struct {
		Genes        map[string]string // gene -> DNA sequence
		VacationDays map[int]int       // month -> days off
	}

	result := salem.Mock(human{}).
		ExecuteToType().([]human)[0]

	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 1, len(result.Genes), "should map with 1 item")

	assert.NotEmpty(t, result.VacationDays, "should create numeric map")
	assert.Equal(t, 1, len(result.VacationDays), "should map with 1 item")
}

func test_map_of_structs(t *testing.T) {
	type dna struct {
		Sequence          string
		PatentHolder      string
		IsManMadeSequence *bool
	}
	type human struct {
		Genes map[string]dna // gene -> DNA sequence
	}

	result := salem.Mock(human{}).
		ExecuteToType().([]human)[0]

	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 1, len(result.Genes), "should map with 1 item")

	keys := []string{}
	for key, _ := range result.Genes {
		keys = append(keys, key)
	}
	item := result.Genes[keys[0]]
	assert.NotEmpty(t, item.Sequence, "should create nested Sequence")
	assert.NotEmpty(t, item.PatentHolder, "should create nested PatentHolder")
	assert.NotNil(t, item.IsManMadeSequence, "should create nested IsManMadeSequence")
}

func test_map_of_slices(t *testing.T) {
	type dna struct {
		Sequence          string
		PatentHolder      string
		IsManMadeSequence *bool
	}
	type human struct {
		Genes map[string][]dna // gene -> DNA sequence
	}

	result := salem.Mock(human{}).
		ExecuteToType().([]human)[0]

	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 1, len(result.Genes), "should map with 1 item")

	keys := []string{}
	for key, _ := range result.Genes {
		keys = append(keys, key)
	}
	item := result.Genes[keys[0]]
	assert.Equal(t, 1, len(item), "should create slice")

	assert.NotEmpty(t, item[0].Sequence, "should create nested Sequence")
	assert.NotEmpty(t, item[0].PatentHolder, "should create nested PatentHolder")
	assert.NotNil(t, item[0].IsManMadeSequence, "should create nested IsManMadeSequence")
}
