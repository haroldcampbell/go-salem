package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mapSuite struct {
	suite.Suite
}

func Test_FactoryMap(t *testing.T) {
	suite.Run(t, new(mapSuite))
}

func (s *mapSuite) Test_primative_types_map() {
	type human struct {
		Genes        map[string]string // gene[name] -> (DNA string sequence)
		VacationDays map[int]int       // month -> days off
	}

	result := salem.Mock(human{}).
		ExecuteToType().([]human)[0]

	t := s.T()
	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 1, len(result.Genes), "should map with 1 item")

	assert.NotEmpty(t, result.VacationDays, "should create numeric map")
	assert.Equal(t, 1, len(result.VacationDays), "should map with 1 item")
}

func (s *mapSuite) Test_map_of_structs() {
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

	t := s.T()
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

func (s *mapSuite) Test_map_of_slices() {
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

	t := s.T()
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

func (s *mapSuite) test_map_of_maps() {
	type part struct {
		Name string // Exon | Intron | Exon
	}
	type human struct {
		Genes map[string]map[string]part // gene[name] -> (part[key]-> part)
	}

	result := salem.Mock(human{}).
		ExecuteToType().([]human)[0]

	t := s.T()
	assert.NotEmpty(t, result.Genes, "should create string map")
	assert.Equal(t, 1, len(result.Genes), "should map with 1 item")
}
