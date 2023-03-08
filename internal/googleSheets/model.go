package googleSheets

import (
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"time"
)

type namedHttpClient struct {
	Name   string
	Client *http.Client
}

func newNamedHttpClient(name string, client *http.Client) namedHttpClient {
	return namedHttpClient{
		Name:   name,
		Client: client,
	}
}

type namedGoogleDrive struct {
	Name  string
	Drive *drive.Service
}

func newNamedGoogleDrive(name string, drive *drive.Service) namedGoogleDrive {
	return namedGoogleDrive{
		Name:  name,
		Drive: drive,
	}
}

type namedGoogleSheet struct {
	Name  string
	Sheet *sheets.Service
}

func newNamedGoogleSheet(name string, sheet *sheets.Service) namedGoogleSheet {
	return namedGoogleSheet{
		Name:  name,
		Sheet: sheet,
	}
}

type clientDetails struct {
	Name   string
	Client *http.Client
	Drive  *drive.Service
	Sheet  *sheets.Service
}

func newClientDetails(name string, client *http.Client, drive *drive.Service, sheet *sheets.Service) clientDetails {
	return clientDetails{
		Name:   name,
		Client: client,
		Drive:  drive,
		Sheet:  sheet,
	}
}

type publishData struct {
	publishTime time.Time
}
