package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

type datastoreCellModel struct {
	GUID          string
	ProjectGUID   string
	DatastoreGUID string
	CellName      string // The name of the header
	RecordType    string // The data type for the record e.g. Date, Number, Text, VarChar, etc
	CellOrdIndex  int    // ordinal index of the cell
}

type columnAggregate struct {
	GUID               string
	Cell               datastoreCellModel
	DisplayTitle       string
	AggregateOperation string
}

type datastoreModel struct {
	GUID         string // Datastore guid
	ProjectGUID  string
	ResourceName string // The raw name of the file
}

type groupByColumn struct {
	GUID             string // Guid for this record
	Cell             datastoreCellModel
	DisplayTitle     string
	ColumnAggregates []columnAggregate
}

type filterColumn struct {
	GUID           string
	Cell           datastoreCellModel
	FilterOrdIndex int // Ordinality of the Filter
}
type filterItem struct {
	FilterValue      interface{}
	IsActive         bool
	FilterColumnGUID string
}

type dataGroupModel struct {
	GUID        string // represents the datagroup guid
	ProjectGUID string
	GroupName   string
	Datastore   *datastoreModel

	FilterColumns     []*filterColumn
	AggregatedColumns []*columnAggregate
	GroupByColumns    []*groupByColumn

	FilterItems []*filterItem // <- Not needed on the client interface
}

type placeHolder struct {
	Inside []*groupByColumn
}

func Test_FactoryNestedPointerSlice(t *testing.T) {
	mockProjectGUID := "PGUID-101"
	mockDataGroupGUID := "DGUID-202"
	mockDatastoreGUID := "DSUID-303"
	mockFilterColumnsGUID := "FilterValue101"

	dc1 := salem.Mock(datastoreCellModel{}).
		Ensure("ProjectGUID", mockProjectGUID).
		ExecuteToType().([]datastoreCellModel)[0]

	// test := salem.Mock(placeHolder{}). // <- [BUG] Shouldn't need the wrapper
	// 					Ensure("Inside", salem.Tap().
	// 						Ensure("Cell", dc1).
	// 						Ensure("ColumnAggregates.Cell", dc1)).
	// 					ExecuteToType().([]placeHolder)
	// mockGroupByColumns := test[0].Inside

	mockDataGroupModel := salem.Mock(dataGroupModel{}).
		Ensure("GUID", mockDataGroupGUID).
		EnsureConstraint("GroupName", salem.ConstrainStringLength(1, 39)). // <- will need an unhappy path to test this
		Ensure("ProjectGUID", mockProjectGUID).
		Ensure("Datastore.ProjectGUID", mockProjectGUID).
		Ensure("Datastore.GUID", mockDatastoreGUID).
		Ensure("Datastore.ResourceName", "R101").
		Ensure("FilterColumns.GUID", mockFilterColumnsGUID). // THIS ISN"T BEING RESPECTED
		Ensure("FilterColumns.Cell", dc1).
		Ensure("AggregatedColumns.GUID", "@@@@@"). // DELETE THIS LINE
		Ensure("AggregatedColumns.Cell", dc1).
		// Ensure("GroupByColumns", mockGroupByColumns). //<- [BUG] This shouldn't fail.
		// // Omit("FilterColumns").
		// Omit("AggregatedColumns").
		// Omit("GroupByColumns").
		// Omit("FilterItems").
		ExecuteToType().([]dataGroupModel)[0]

	// str := utils.PrettyMongoString(mockDataGroupModel)
	// fmt.Printf("mockDataGroupModel: %v\n", str)

	assert.Equal(t, mockDataGroupGUID, mockDataGroupModel.GUID)
	assert.Equal(t, mockProjectGUID, mockDataGroupModel.ProjectGUID)

	assert.Equal(t, mockDatastoreGUID, mockDataGroupModel.Datastore.GUID)
	assert.Equal(t, mockProjectGUID, mockDataGroupModel.Datastore.ProjectGUID)
	assert.Equal(t, "R101", mockDataGroupModel.Datastore.ResourceName)

	assert.Equal(t, mockFilterColumnsGUID, mockDataGroupModel.FilterColumns[0].GUID)
	assert.Equal(t, dc1, mockDataGroupModel.FilterColumns[0].Cell)

	assert.Equal(t, "@@@@@", mockDataGroupModel.AggregatedColumns[0].GUID)
	assert.Equal(t, dc1, mockDataGroupModel.AggregatedColumns[0].Cell)
}
