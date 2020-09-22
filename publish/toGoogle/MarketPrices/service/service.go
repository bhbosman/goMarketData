package service

import (
	"context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

type Service struct {
	Name         string
	DriveService *drive.Service
	SheetService *sheets.Service
	Id           string
}

func (self *Service) Start(ctx context.Context) error {
	fileList, err := self.DriveService.Files.List().Do()
	if err != nil {
		return err
	}
	for _, fileInfo := range fileList.Files {
		if fileInfo.Name == "MarketPrices-Latest" {
			self.Id = fileInfo.Id
			break
		}
	}
	if self.Id == "" {
		createCall := self.SheetService.Spreadsheets.Create(
			&sheets.Spreadsheet{
				Properties: &sheets.SpreadsheetProperties{
					Title: "MarketPrices-Latest",
				},
			})
		spreadSheet, err := createCall.Do()
		if err != nil {
			return nil
		}
		self.Id = spreadSheet.SpreadsheetId
		println(spreadSheet.SpreadsheetUrl)
	}
	return nil
}

func (self *Service) Stop(ctx context.Context) error {
	return nil
}

func NewService(name string, driveService *drive.Service, sheetService *sheets.Service) *Service {
	return &Service{Name: name, DriveService: driveService, SheetService: sheetService}
}
