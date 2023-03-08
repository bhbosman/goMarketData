package googleSheets

import (
	"context"
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"google.golang.org/api/sheets/v4"
)

type IGoogleSheets interface {
	ISendMessage.ISendMessage
	ISendMessage.IMultiSendMessage
}

type IGoogleSheetsService interface {
	IGoogleSheets
	IFxService.IFxServices
}

type IGoogleSheetsData interface {
	IGoogleSheets
	IDataShutDown.IDataShutDown
	Start(ctx context.Context) error
}

type requestResult struct {
	clearRange    string
	updatedValues *sheets.Request
}

type iNamedRange interface {
	Name() string
	SetIsDirty(bool)
	IsDirty() bool
	CreateRange() *sheets.NamedRange
	ValueRange() (bool, *sheets.ValueRange)
	RangeUpdate() (doUpdate bool, requestResult requestResult)
	ReadFromSheet(namedRangeMap map[string]*sheets.NamedRange)
}

type iSheetInformation interface {
	CreateSheet() *sheets.Sheet
	Title() string
	ReadFromSheet(sheetMap map[string]*sheets.Sheet)
	SheetId() int64
}
