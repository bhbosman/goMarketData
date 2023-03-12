package googleSheets

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

type data struct {
	MessageRouter                       *messageRouter.MessageRouter
	fmdServiceHelper                    fullMarketDataHelper.IFullMarketDataHelper
	pubSub                              *pubsub.PubSub
	HttpClient                          *http.Client
	DriveService                        *drive.Service
	SheetService                        *sheets.Service
	SpreadsheetId                       string
	logger                              *zap.Logger
	namedRangeMap                       map[string]iNamedRange
	sheetInformationMap                 map[int64]iSheetInformation
	availableInstrumentsNamedRange      *availableInstrumentsNamedRange
	availableInstrumentSheetInformation *availableInstrumentSheetInformation
}

func (self *data) Start(ctx context.Context) error {
	SpreadsheetId, err := self.initialLoadOfSpreadSheet()
	if err != nil {
		return err
	}
	self.SpreadsheetId = SpreadsheetId
	return nil
}

func (self *data) MultiSend(messages ...interface{}) {
	self.MessageRouter.MultiRoute(messages...)
}

func (self *data) Send(message interface{}) error {
	self.MessageRouter.Route(message)
	return nil
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) handlePublishData(msg *publishData) {
	self.doPublish()
}

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_InstrumentList_Response(msg *stream.FullMarketData_InstrumentList_Response) {
	self.availableInstrumentsNamedRange.instrumentMap = make(map[string]fullMarketDataManagerService.InstrumentStatus)
	for _, instrument := range msg.Instruments {
		self.availableInstrumentsNamedRange.instrumentMap[instrument.Instrument] = fullMarketDataManagerService.InstrumentStatus{
			Instrument: instrument.Instrument,
			Status:     instrument.Status,
		}
	}
	self.availableInstrumentsNamedRange.isDirty = true
}

func (self *data) doPublish() {
	if self.SpreadsheetId != "" {
		var valueRangeData []*sheets.ValueRange
		var requests []*sheets.Request
		var clearRanges []string

		for _, namedRange := range self.namedRangeMap {
			update, request := namedRange.RangeUpdate()
			if update {
				requests = append(requests, request.updatedValues)
				if request.clearRange != "" {
					clearRanges = append(clearRanges, request.clearRange)
				}
			}
			if namedRange.IsDirty() {
				b, valueRange := namedRange.ValueRange()
				if b {
					valueRangeData = append(valueRangeData, valueRange)
				}
			}
		}

		if len(requests) > 0 {
			updateCall := self.SheetService.Spreadsheets.BatchUpdate(
				self.SpreadsheetId,
				&sheets.BatchUpdateSpreadsheetRequest{
					IncludeSpreadsheetInResponse: false,
					Requests:                     requests,
					ResponseIncludeGridData:      false,
					ResponseRanges:               nil,
					ForceSendFields:              nil,
					NullFields:                   nil,
				},
			)
			_, err := updateCall.Do()
			if err != nil {
				self.logger.Error("On publish", zap.Error(err))
			}
		}

		if len(clearRanges) > 0 {
			clearCall := self.SheetService.Spreadsheets.Values.BatchClear(
				self.SpreadsheetId,
				&sheets.BatchClearValuesRequest{
					Ranges:          clearRanges,
					ForceSendFields: nil,
					NullFields:      nil,
				},
			)
			_, err := clearCall.Do()
			if err != nil {
				self.logger.Error("On publish", zap.Error(err))
			}
		}

		if len(valueRangeData) > 0 {
			updateCall := self.SheetService.Spreadsheets.Values.BatchUpdate(
				self.SpreadsheetId,
				&sheets.BatchUpdateValuesRequest{
					ValueInputOption: "RAW",
					Data:             valueRangeData,
				},
			)
			_, err := updateCall.Do()
			if err != nil {
				self.logger.Error("On publish", zap.Error(err))
			}
		}
	}
}

func (self *data) initialLoadOfSpreadSheet() (string, error) {
	id := ""
	fileList, err := self.DriveService.Files.List().Do()
	if err != nil {
		return "", err
	}
	for _, fileInfo := range fileList.Files {
		if fileInfo.Name == "MarketPrices-Latest" {
			id = fileInfo.Id
			break
		}
	}

	if id == "" {
		return self.createSpreadSheet()
	}
	return self.loadAndReadSheet(id)
}

func (self *data) createSpreadSheet() (string, error) {
	var NamedRanges []*sheets.NamedRange
	for _, namedRange := range self.namedRangeMap {
		NamedRanges = append(NamedRanges, namedRange.CreateRange())
	}

	var Sheets []*sheets.Sheet
	for _, sheetInformation := range self.sheetInformationMap {
		sheet := sheetInformation.CreateSheet()
		Sheets = append(Sheets, sheet)
	}

	createCall := self.SheetService.Spreadsheets.Create(
		&sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: "MarketPrices-Latest",
			},
			Sheets:      Sheets,
			NamedRanges: NamedRanges,
		},
	)
	spreadSheet, err := createCall.Do()
	if err != nil {
		return "", err
	}
	return spreadSheet.SpreadsheetId, nil
}

func (self *data) loadAndReadSheet(id string) (string, error) {
	spreadSheet, err := self.SheetService.Spreadsheets.Get(id).Do()
	if err != nil {
		return "", err
	}
	namedRangeMap := make(map[string]*sheets.NamedRange)
	for _, namedRange := range spreadSheet.NamedRanges {
		namedRangeMap[namedRange.Name] = namedRange
	}
	for _, namedRange := range self.namedRangeMap {
		namedRange.ReadFromSheet(namedRangeMap)
	}

	sheetMap := make(map[string]*sheets.Sheet)
	for _, sheet := range spreadSheet.Sheets {
		sheetMap[sheet.Properties.Title] = sheet
	}
	for _, sheetInformation := range self.sheetInformationMap {
		sheetInformation.ReadFromSheet(sheetMap)
	}

	return spreadSheet.SpreadsheetId, nil
}

func newData(
	fmdServiceHelper fullMarketDataHelper.IFullMarketDataHelper,
	pubSub *pubsub.PubSub,
	HttpClient *http.Client,
	DriveService *drive.Service,
	SheetService *sheets.Service,
	logger *zap.Logger,
) (IGoogleSheetsData, error) {
	result := &data{
		MessageRouter:                       messageRouter.NewMessageRouter(),
		fmdServiceHelper:                    fmdServiceHelper,
		pubSub:                              pubSub,
		HttpClient:                          HttpClient,
		DriveService:                        DriveService,
		SheetService:                        SheetService,
		logger:                              logger,
		namedRangeMap:                       make(map[string]iNamedRange),
		availableInstrumentsNamedRange:      newStaticNamedRange("AvailableData"),
		availableInstrumentSheetInformation: newAvailableInstrumentSheetInformation("Available data", 1),
		sheetInformationMap:                 make(map[int64]iSheetInformation),
	}

	// handlers
	result.MessageRouter.Add(result.handleEmptyQueue)
	result.MessageRouter.Add(result.handlePublishData)
	result.MessageRouter.Add(result.handleFullMarketData_InstrumentList_Response)

	// map sheet objects to each other
	result.availableInstrumentsNamedRange.SetSheet(result.availableInstrumentSheetInformation)

	// add to maps
	result.namedRangeMap[result.availableInstrumentsNamedRange.Name()] = result.availableInstrumentsNamedRange
	result.sheetInformationMap[result.availableInstrumentSheetInformation.SheetId()] = result.availableInstrumentSheetInformation

	return result, nil
}
