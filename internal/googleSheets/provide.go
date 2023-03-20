package googleSheets

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goConn"
	service2 "github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"os"
	"path/filepath"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						HttpClients   []namedHttpClient  `group:"MarketPricesClients"`
						GoogleDrives  []namedGoogleDrive `group:"MarketPricesClients"`
						GoogleSheets  []namedGoogleSheet `group:"MarketPricesClients"`
						ClientDetails []clientDetails    `group:"MarketPricesClients"`
					},
				) (map[string]namedHttpClient, map[string]namedGoogleDrive, map[string]namedGoogleSheet, map[string]clientDetails, error) {
					httpClientMap := make(map[string]namedHttpClient)
					googleDriveMap := make(map[string]namedGoogleDrive)
					googleSheetsMap := make(map[string]namedGoogleSheet)
					clientDetailsMap := make(map[string]clientDetails)

					for _, httpClient := range params.HttpClients {
						httpClientMap[httpClient.Name] = httpClient
					}
					for _, googleDrive := range params.GoogleDrives {
						googleDriveMap[googleDrive.Name] = googleDrive
					}
					for _, googleSheet := range params.GoogleSheets {
						googleSheetsMap[googleSheet.Name] = googleSheet
					}
					for _, clientDetail := range params.ClientDetails {
						clientDetailsMap[clientDetail.Name] = clientDetail
					}
					return httpClientMap, googleDriveMap, googleSheetsMap, clientDetailsMap, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Name: "MarketPricesApplication",
				Target: func(
					params struct {
						fx.In
						Logger *zap.Logger
					},
				) (*oauth2.Config, error) {
					homePath, _ := os.UserHomeDir()
					b, err := os.ReadFile(filepath.Join(homePath, ".MarketPrices", "credentials.json"))
					if err != nil {
						params.Logger.Error("Unable to read client secret file", zap.Error(err))
					}
					// If modifying these scopes, delete your previously saved token.json.
					return google.ConfigFromJSON(
						b,
						"https://www.googleapis.com/auth/spreadsheets.readonly",
						"https://www.googleapis.com/auth/drive.file")
				},
			},
		),
		ProvideNamedObjects(),
		fx.Invoke(
			func(
				params struct {
					fx.In
					PubSub                 *pubsub.PubSub    `name:"Application"`
					ApplicationContext     context.Context   `name:"Application"`
					Oauth2Config           []*oauth2.Config  `group:"MarketPricesApplication"`
					NamedHttpClient        []namedHttpClient `group:"MarketPricesClients"`
					Lifecycle              fx.Lifecycle
					Logger                 *zap.Logger
					UniqueReferenceService interfaces.IUniqueReferenceService
					UniqueSessionNumber    interfaces.IUniqueSessionNumber
					GoFunctionCounter      GoFunctionCounter.IService
					FmdServiceHelper       fullMarketDataHelper.IFullMarketDataHelper
					FmdManagerService      fullMarketDataManagerService.IFmdManagerService
					FxManagerService       service2.IFxManagerService
					HttpClientsMap         map[string]namedHttpClient
					GoogleDrivesMap        map[string]namedGoogleDrive
					GoogleSheetsMap        map[string]namedGoogleSheet
					ClientDetailsMap       map[string]clientDetails
				},
			) error {
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							for k, clientDetailsInstance := range params.ClientDetailsMap {
								err := params.FxManagerService.Add(k,
									func() (messages.IApp, goConn.ICancellationContext, error) {
										onData := func() (IGoogleSheetsData, error) {
											return newData(
												params.FmdServiceHelper,
												params.PubSub,
												clientDetailsInstance.Client,
												clientDetailsInstance.Drive,
												clientDetailsInstance.Sheet,
												params.Logger,
											)
										}
										namedLogger := params.Logger.Named(k)
										ctx, cancelFunc := context.WithCancel(params.ApplicationContext)
										cancellationContext, err := goConn.NewCancellationContextNoCloser(
											k,
											cancelFunc,
											ctx, namedLogger,
										)
										if err != nil {
											return nil, nil, err
										}
										googleSheetService, err := newService(
											cancellationContext,
											onData,
											params.Logger,
											params.PubSub,
											params.GoFunctionCounter,
											params.FmdServiceHelper,
											params.FmdManagerService,
											clientDetailsInstance,
										)
										if err != nil {
											return nil, nil, err
										}
										return newAppWrapper(googleSheetService), cancellationContext, nil
									},
								)
								if err != nil {
									return err
								}
							}
							return nil
						},
						OnStop: func(ctx context.Context) error {
							return nil
						},
					},
				)
				return nil
			},
		),
	)
}

func ProvideNamedObjects() fx.Option {
	getClient := func(config *oauth2.Config, userInformation *UserInformation) (*http.Client, error) {
		return config.Client(context.Background(), userInformation.OauthToken), nil
	}
	var options []fx.Option
	homePath, _ := os.UserHomeDir()
	path := filepath.Join(homePath, ".MarketPrices", "clients")
	err := filepath.Walk(
		path,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			fn := path
			provide := fx.Provide(
				fx.Annotated{
					Group: "MarketPricesClients",
					Target: func(
						params struct {
							fx.In
							AppContext context.Context `name:"Application"`

							MarketPricesConfig *oauth2.Config `name:"MarketPricesApplication"`
						},
					) (namedHttpClient, namedGoogleDrive, namedGoogleSheet, clientDetails, error) {
						userInformation, err := UserInformationFrom(fn)
						client, err := getClient(params.MarketPricesConfig, userInformation)
						if err != nil {
							return namedHttpClient{}, namedGoogleDrive{}, namedGoogleSheet{}, clientDetails{}, err
						}
						namedGoogleDriveInstance, err := drive.NewService(params.AppContext, option.WithHTTPClient(client))
						if err != nil {
							return namedHttpClient{}, namedGoogleDrive{}, namedGoogleSheet{}, clientDetails{}, err
						}
						namedGoogleSheetsInstance, err := sheets.NewService(params.AppContext, option.WithHTTPClient(client))
						if err != nil {
							return namedHttpClient{}, namedGoogleDrive{}, namedGoogleSheet{}, clientDetails{}, err
						}
						return newNamedHttpClient(fn, client),
							newNamedGoogleDrive(fn, namedGoogleDriveInstance),
							newNamedGoogleSheet(fn, namedGoogleSheetsInstance),
							newClientDetails(fn, client, namedGoogleDriveInstance, namedGoogleSheetsInstance),
							nil
					},
				},
			)
			options = append(options, provide)
			return nil
		},
	)
	if err != nil {
		return fx.Error(err)
	}
	return fx.Options(options...)
}
