package service

import (
	"context"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
	"sort"
	"time"
)

type top5Data struct {
	data    *marketDataStream.PublishTop5
	//touched bool
	plate   [][]interface{}
}

type TimerMessage struct {
}

type MarketPricesLatestService struct {
	Name               string
	DriveService       *drive.Service
	SheetService       *sheets.Service
	Id                 string
	PubSub             *pubsub.PubSub
	cancelContext      context.Context
	cancelFunc         context.CancelFunc
	MessageRouter      *messageRouter.MessageRouter
	Top5Map            map[string]*top5Data
	namedRanges        map[string]bool
	availableRangeName map[string]bool
}

func (self *MarketPricesLatestService) Start(ctx context.Context) error {
	self.availableRangeName["TimeStamp"] = true
	channel := self.PubSub.Sub("Top5Data")
	go self.start(channel)
	return nil
}

func (self *MarketPricesLatestService) Stop(ctx context.Context) error {
	self.cancelFunc()
	return nil
}

func SendContext(ctx context.Context, ch chan<- interface{}, i interface{}) bool {
	select {
	case <-ctx.Done():
		return false
	case ch <- i:
		return true
	}
}

func (self *MarketPricesLatestService) start(channel chan interface{}) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	reloadSpreadSheetTicker := time.NewTicker(time.Second * 15)
	defer reloadSpreadSheetTicker.Stop()

	reloadSpreadSheetFn := func() {
		spreadsheet, err := self.SheetService.Spreadsheets.Get(self.Id).Do()
		if err != nil {
			return
		}
		SendContext(self.cancelContext, channel, spreadsheet)
	}

	go func() {
		spreadsheet, err := self.initialLoadOfSpreadSheet()
		if err != nil {
			return
		}
		SendContext(self.cancelContext, channel, spreadsheet)
	}()
loop:
	for true {
		select {
		case _, ok := <-reloadSpreadSheetTicker.C:
			if ok {
				go reloadSpreadSheetFn()
			}
		case _, ok := <-ticker.C:
			if ok {
				_, _ = self.MessageRouter.Route(&TimerMessage{})
			}
			continue loop
		case <-self.cancelContext.Done():
			break loop
		case v, ok := <-channel:
			if ok {
				_, _ = self.MessageRouter.Route(v)
				if len(channel) == 0 {
					_, _ = self.MessageRouter.Route(&rxgo.EmptyQueue{})
				}
				continue loop
			}
			break loop
		}
	}
	// unsubscribe
	self.PubSub.Unsub(channel)
	// flush to avoid deadlocks
	for range channel {
	}
}

func (self *MarketPricesLatestService) handlePublishTop(incomingMessage *marketDataStream.PublishTop5) error {
	if data, ok := self.Top5Map[incomingMessage.Instrument]; ok {
		data.data = incomingMessage
		//data.touched = true
	} else {
		self.availableRangeName[incomingMessage.Instrument] = true
		plate := make([][]interface{}, 5)
		plate[0] = make([]interface{}, 4)
		plate[1] = make([]interface{}, 4)
		plate[2] = make([]interface{}, 4)
		plate[3] = make([]interface{}, 4)
		plate[4] = make([]interface{}, 4)
		self.Top5Map[incomingMessage.Instrument] = &top5Data{
			data:    incomingMessage,
			//touched: true,
			plate:   plate,
		}
	}
	return nil
}

func (self *MarketPricesLatestService) handleEmptyQueue(incomingMessage *rxgo.EmptyQueue) error {
	return nil
}

func (self *MarketPricesLatestService) getAvailableData() *sheets.ValueRange {
	var resetFlag []string
	var ss []string
	if len(self.availableRangeName) > 0 {
		for k, v := range self.availableRangeName {
			ss = append(ss, k)
			if v {
				resetFlag = append(resetFlag, k)
			}
		}
		for _, s := range resetFlag {
			self.availableRangeName[s] = false
		}
		sort.Strings(ss)

		values := make([][]interface{}, len(ss))
		for i := 0; i < len(values); i++ {
			values[i] = make([]interface{}, 1)
			values[i][0] = ss[i]
		}

		valueRange := &sheets.ValueRange{
			MajorDimension:  "ROWS",
			Range:           "AvailableData",
			Values:          values,
			ServerResponse:  googleapi.ServerResponse{},
			ForceSendFields: nil,
			NullFields:      nil,
		}
		return valueRange
	}
	return nil
}

func (self *MarketPricesLatestService) getTimeStamp() *sheets.ValueRange {
	if _, ok := self.namedRanges["TimeStamp"]; ok {
		values := make([][]interface{}, 1)
		values[0] = make([]interface{}, 1)
		values[0][0] = time.Now().Format(time.RFC3339)

		valueRange := &sheets.ValueRange{
			MajorDimension:  "ROWS",
			Range:           "TimeStamp",
			Values:          values,
			ServerResponse:  googleapi.ServerResponse{},
			ForceSendFields: nil,
			NullFields:      nil,
		}
		return valueRange
	}
	return nil
}

func (self *MarketPricesLatestService) handleSpreadsheet(incomingMessage *sheets.Spreadsheet) error {
	self.Id = incomingMessage.SpreadsheetId
	println(incomingMessage.SpreadsheetUrl)
	self.namedRanges = make(map[string]bool)
	for _, namedRange := range incomingMessage.NamedRanges {
		self.namedRanges[namedRange.Name] = true
	}
	return nil
}

func (self *MarketPricesLatestService) handleTimerMessage(incomingMessage *TimerMessage) error {
	if self.Id != "" {
		var Data []*sheets.ValueRange
		valueRange := self.getAvailableData()
		if valueRange != nil {
			Data = append(Data, valueRange)
		}
		valueRange = self.getTimeStamp()
		if valueRange != nil {
			Data = append(Data, valueRange)
		}
		valueRanges := self.getTop5()
		for _, valueRange = range valueRanges{
			if valueRange != nil {
				Data = append(Data, valueRange)
			}
		}

		if len(Data) > 0 {
			BatchUpdateValuesRequest := &sheets.BatchUpdateValuesRequest{
				Data:             Data,
				ValueInputOption: "RAW",
			}
			updateCall := self.SheetService.Spreadsheets.Values.BatchUpdate(self.Id, BatchUpdateValuesRequest)
			_, err := updateCall.Do()
			if err != nil {
				println(err.Error())
			}
		}
	}
	return nil
}
func (self *MarketPricesLatestService) createSpreadSheet() (*sheets.Spreadsheet, error) {
	sheet01 := &sheets.Sheet{
		Properties: &sheets.SheetProperties{
			Title:   "Available Data",
			SheetId: 1,
		},
	}
	namedRange := &sheets.NamedRange{
		Name:         "AvailableData",
		NamedRangeId: "AvailableData",
		Range: &sheets.GridRange{
			EndColumnIndex:   1,
			EndRowIndex:      0,
			SheetId:          1,
			StartColumnIndex: 0,
			StartRowIndex:    0,
			ForceSendFields:  nil,
			NullFields:       nil,
		},
		ForceSendFields: nil,
		NullFields:      nil,
	}
	var NamedRanges []*sheets.NamedRange
	var Sheets []*sheets.Sheet

	Sheets = append(Sheets, sheet01)
	NamedRanges = append(NamedRanges, namedRange)
	createCall := self.SheetService.Spreadsheets.Create(
		&sheets.Spreadsheet{
			NamedRanges: NamedRanges,
			Properties: &sheets.SpreadsheetProperties{
				Title: "MarketPrices-Latest",
			},
			Sheets: Sheets,
		})
	spreadSheet, err := createCall.Do()
	if err != nil {
		return nil, err
	}
	return spreadSheet, nil
}

func (self *MarketPricesLatestService) initialLoadOfSpreadSheet() (*sheets.Spreadsheet, error) {
	id := ""
	fileList, err := self.DriveService.Files.List().Do()
	if err != nil {
		return nil, err
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
	return self.SheetService.Spreadsheets.Get(id).Do()
}

func (self *MarketPricesLatestService) getTop5() []*sheets.ValueRange {
	var result []*sheets.ValueRange
	for k, v := range self.Top5Map {
		if _, ok := self.namedRanges[k]; ok {
			//if v.touched {
			//	v.touched = false
				for i := 0; i < 5; i++ {
					if len(v.data.Bid) > i {
						v.plate[i][0] = v.data.Bid[i].Volume
						v.plate[i][1] = v.data.Bid[i].Price
					} else {
						v.plate[i][0] = 0.0
						v.plate[i][1] = 0.0
					}
					if len(v.data.Ask) > i {
						v.plate[i][2] = v.data.Ask[i].Price
						v.plate[i][3] = v.data.Ask[i].Volume
					} else {
						v.plate[i][2] = 0.0
						v.plate[i][3] = 0.0
					}
				}
				valueRange := &sheets.ValueRange{
					MajorDimension:  "ROWS",
					Range:           k,
					Values:          v.plate,
					ServerResponse:  googleapi.ServerResponse{},
					ForceSendFields: nil,
					NullFields:      nil,
				}
				result = append(result, valueRange)
			//}
		}
	}
	return result
}

func NewService(
	name string,
	driveService *drive.Service,
	sheetService *sheets.Service,
	PubSub *pubsub.PubSub) *MarketPricesLatestService {
	cancel, cancelFunc := context.WithCancel(context.Background())

	result := &MarketPricesLatestService{
		Name:               name,
		DriveService:       driveService,
		SheetService:       sheetService,
		Id:                 "",
		PubSub:             PubSub,
		cancelContext:      cancel,
		cancelFunc:         cancelFunc,
		MessageRouter:      messageRouter.NewMessageRouter(),
		Top5Map:            make(map[string]*top5Data),
		namedRanges:        make(map[string]bool),
		availableRangeName: make(map[string]bool),
	}
	_ = result.MessageRouter.Add(result.handlePublishTop)
	_ = result.MessageRouter.Add(result.handleEmptyQueue)
	_ = result.MessageRouter.Add(result.handleTimerMessage)
	_ = result.MessageRouter.Add(result.handleSpreadsheet)

	return result
}
