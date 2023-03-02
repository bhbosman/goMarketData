package googleSheets

import (
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

type availableInstrumentSheetInformation struct {
	data *sheets.Sheet
}

func (self *availableInstrumentSheetInformation) ReadFromSheet(sheetMap map[string]*sheets.Sheet) {
	if v, ok := sheetMap[self.data.Properties.Title]; ok {
		self.data = v
	}
}

func (self *availableInstrumentSheetInformation) Title() string {
	return self.data.Properties.Title
}

func (self *availableInstrumentSheetInformation) CreateSheet() *sheets.Sheet {
	return self.data
}

func (self *availableInstrumentSheetInformation) SheetId() int64 {
	return self.data.Properties.SheetId
}

func newAvailableInstrumentSheetInformation(
	title string,
	sheetId int64,
) *availableInstrumentSheetInformation {
	return &availableInstrumentSheetInformation{
		data: &sheets.Sheet{
			Properties: &sheets.SheetProperties{
				DataSourceSheetProperties: &sheets.DataSourceSheetProperties{
					Columns: nil,
					DataExecutionStatus: &sheets.DataExecutionStatus{
						ErrorCode:       "",
						ErrorMessage:    "",
						LastRefreshTime: "",
						State:           "",
						ForceSendFields: nil,
						NullFields:      nil,
					},
					DataSourceId:    "",
					ForceSendFields: nil,
					NullFields:      nil,
				},
				GridProperties: &sheets.GridProperties{
					ColumnCount:             0,
					ColumnGroupControlAfter: false,
					FrozenColumnCount:       0,
					FrozenRowCount:          0,
					HideGridlines:           false,
					RowCount:                0,
					RowGroupControlAfter:    false,
					ForceSendFields:         nil,
					NullFields:              nil,
				},
				Hidden:      false,
				Index:       0,
				RightToLeft: false,
				SheetId:     sheetId,
				SheetType:   "",
				TabColor: &sheets.Color{
					Alpha:           0,
					Blue:            0,
					Green:           0,
					Red:             0,
					ForceSendFields: nil,
					NullFields:      nil,
				},
				TabColorStyle: &sheets.ColorStyle{
					RgbColor: &sheets.Color{
						Alpha:           0,
						Blue:            0,
						Green:           0,
						Red:             0,
						ForceSendFields: nil,
						NullFields:      nil,
					},
					ThemeColor:      "",
					ForceSendFields: nil,
					NullFields:      nil,
				},
				Title: title,
				ServerResponse: googleapi.ServerResponse{
					HTTPStatusCode: 0,
					Header:         nil,
				},
				ForceSendFields: nil,
				NullFields:      nil,
			},
		},
	}
}
