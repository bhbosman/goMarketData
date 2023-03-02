package googleSheets

import (
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
	"math"
	"sort"
)

type availableInstrumentsNamedRange struct {
	data          *sheets.NamedRange
	instrumentMap map[string]fullMarketDataManagerService.InstrumentStatus
	isDirty       bool
	sheet         iSheetInformation
}

func (self *availableInstrumentsNamedRange) ReadFromSheet(namedRangeMap map[string]*sheets.NamedRange) {
	if value, ok := namedRangeMap[self.data.Name]; ok {
		self.data = value
	}
}

func (self *availableInstrumentsNamedRange) RangeUpdate() (bool, requestResult) {
	result := requestResult{
		clearRange:    "",
		updatedValues: nil,
	}
	EndRowIndex := int64(len(self.instrumentMap))

	delta := EndRowIndex - self.data.Range.EndRowIndex
	b := delta == 0
	if !b {
		if delta < 0 {
			startColumn, _ := indexToColumn(self.data.Range.StartColumnIndex + 1)
			endColumn, _ := indexToColumn(self.data.Range.EndColumnIndex)

			result.clearRange = fmt.Sprintf("%v!%v%v:%v%v",
				self.sheet.Title(),
				startColumn, EndRowIndex+1,
				endColumn, self.data.Range.EndRowIndex,
			)
		}
		self.data.Range.EndRowIndex = EndRowIndex
		result.updatedValues = &sheets.Request{
			UpdateDimensionProperties: nil,
			UpdateNamedRange: &sheets.UpdateNamedRangeRequest{
				Fields:          "*",
				NamedRange:      self.data,
				ForceSendFields: nil,
				NullFields:      nil,
			},
		}
		return true, result
	}
	return false, result
}

func (self *availableInstrumentsNamedRange) CreateRange() *sheets.NamedRange {
	return self.data
}

func (self *availableInstrumentsNamedRange) SetIsDirty(b bool) {
	self.isDirty = b
}

func (self *availableInstrumentsNamedRange) IsDirty() bool {
	return self.isDirty
}

func (self *availableInstrumentsNamedRange) Name() string {
	return self.data.Name
}

func (self *availableInstrumentsNamedRange) SetSheet(sheet iSheetInformation) {
	self.sheet = sheet
	self.data.Range.SheetId = sheet.SheetId()
}

func (self *availableInstrumentsNamedRange) ValueRange() (bool, *sheets.ValueRange) {
	if self.isDirty {
		ss := make(fullMarketDataManagerService.InstrumentStatusArray, len(self.instrumentMap))
		index := 0
		for _, instrumentStatus := range self.instrumentMap {
			ss[index] = instrumentStatus
			index++
		}
		sort.Sort(ss)
		values := make([][]interface{}, len(ss))
		for i := 0; i < len(ss); i++ {
			values[i] = make([]interface{}, 2)
			values[i][0] = ss[i].Instrument
			values[i][1] = ss[i].Status
		}

		valueRange := &sheets.ValueRange{
			MajorDimension:  "ROWS",
			Range:           "AvailableData",
			Values:          values,
			ServerResponse:  googleapi.ServerResponse{},
			ForceSendFields: nil,
			NullFields:      nil,
		}
		self.isDirty = false
		return true, valueRange
	}
	return false, nil
}

func newStaticNamedRange(name string) *availableInstrumentsNamedRange {
	return &availableInstrumentsNamedRange{
		data: &sheets.NamedRange{
			Name:         name,
			NamedRangeId: name,
			Range: &sheets.GridRange{
				EndColumnIndex:   2,
				EndRowIndex:      2,
				SheetId:          0,
				StartColumnIndex: 0,
				StartRowIndex:    0,
				ForceSendFields:  nil,
				NullFields:       nil,
			},
			ForceSendFields: nil,
			NullFields:      nil,
		},
		instrumentMap: make(map[string]fullMarketDataManagerService.InstrumentStatus),
		isDirty:       false,
	}
}

// indexToColumn takes in an index value & converts it to A1 Notation
// Index 1 is Column A
// E.g. 3 == C, 29 == AC, 731 == ABC
func indexToColumn(index int64) (string, error) {

	// Validate index size
	var maxIndex int64 = 18278
	if index > maxIndex {
		return "", fmt.Errorf("index cannot be greater than %v (column ZZZ)", maxIndex)
	}

	// Get column from index
	l := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if index > 26 {
		letterA, _ := indexToColumn(int64(math.Floor(float64(index-1) / 26)))
		letterB, _ := indexToColumn(index % 26)
		return letterA + letterB, nil
	} else {
		if index == 0 {
			index = 26
		}
		return string(l[index-1]), nil
	}

}

// columnToIndex takes in A1 Notation & converts it to an index value
// Column A is index 1
// E.g. C == 3, AC == 29, ABC == 731
func columnToIndex(column string) (int, error) {

	// Calculate index from column string
	var index int
	var a uint8 = "A"[0]
	var z uint8 = "Z"[0]
	var alphabet = z - a + 1
	i := 1
	for n := len(column) - 1; n >= 0; n-- {
		r := column[n]
		if r < a || r > z {
			return 0, fmt.Errorf("invalid character in column, expected A-Z but got [%c]", r)
		}
		runePos := int(r-a) + 1
		index += runePos * int(math.Pow(float64(alphabet), float64(i-1)))
		i++
	}

	// Return column index & success
	return index, nil

}
