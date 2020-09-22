package main

import (
	"context"
	"github.com/bhbosman/goMarketData/internal/buildMarketDataSheet"
	"github.com/bhbosman/goMarketData/internal/incomingMarketData"
	"github.com/bhbosman/goMarketData/publish/toGoogle/MarketPrices/register"
	"github.com/bhbosman/goMarketData/publish/toGoogle/MarketPrices/service"
	app2 "github.com/bhbosman/gocommon/app"
	"os"

	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"log"
	"time"
)

type LunoStreamAppSettings struct {
	logger *log.Logger
}

func main() {
	settings := &LunoStreamAppSettings{
		//logger: log.New(&stream.NullWriter{}, "", log.LstdFlags),
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	var services []*service.Service

	app := fx.New(
		fx.Logger(settings.logger),
		gologging.ProvideLogFactory(settings.logger, nil),
		app2.RegisterRootContext(),
		connectionManager.RegisterDefaultConnectionManager(),
		impl.RegisterAllConnectionRelatedServices(),
		register.ProvideMarketPriceGoogleApplication(),
		register.ProvideMarketPricesClients(),
		register.ProvideMarketPricesServices(),
		fx.Provide(fx.Annotated{Target: buildMarketDataSheet.NewSheetBuilder}),
		incomingMarketData.ProvideLunoMarketDataDialer(1, "tcp4://127.0.0.1:3001"),
		fx.Invoke(func(params struct {
			fx.In
			SheetBuilder *buildMarketDataSheet.SheetBuilder
			Lifecycle    fx.Lifecycle
		}) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: params.SheetBuilder.Start,
					OnStop: params.SheetBuilder.Stop,
				})
		}),
		fx.Invoke(
			func(params struct {
				fx.In
				Lifecycle fx.Lifecycle
				Apps      []*fx.App `group:"Apps"`
			}) {
				for _, item := range params.Apps {
					localApp := item
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return localApp.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							return localApp.Stop(ctx)
						},
					})
				}
			}),
		fx.Populate(&services),
	)
	err := app.Err()
	if err != nil {
		return
	}
	app.Run()

	// allow shutdown to complete
	time.Sleep(time.Second)

	//values := make([][]interface{}, 2)
	//values[0] = make([]interface{}, 2)
	//values[1] = make([]interface{}, 2)
	//values[0][0] = 1
	//values[0][1] = 2
	//values[1][0] = 3
	//values[1][1] = 4
	//for _, service := range services {
	//	updateCall := service.SheetService.Spreadsheets.Values.Update(
	//		service.Id,
	//		"A1:B2",
	//		&sheets.ValueRange{
	//			MajorDimension:  "ROWS",
	//			Range:           "",
	//			Values:          values,
	//			ServerResponse:  googleapi.ServerResponse{},
	//			ForceSendFields: nil,
	//			NullFields:      nil,
	//		})
	//	updateCall.ValueInputOption("RAW")
	//	valuesResponse, err := updateCall.Do()
	//	if err != nil {
	//		println(err.Error())
	//	}
	//	if valuesResponse != nil {
	//	}
	//	println("good")
	//}
}
