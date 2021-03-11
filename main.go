package main

import (
	"context"
	"github.com/bhbosman/goMarketData/internal/incomingMarketData"
	"github.com/bhbosman/goMarketData/internal/publish/toGoogle/MarketPrices/register"
	"github.com/bhbosman/goMarketData/internal/publish/toGoogle/MarketPrices/service"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/cskr/pubsub"
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
	var services []*service.MarketPricesLatestService

	pubSub := pubsub.New(32)

	app := fx.New(
		fx.StartTimeout(time.Hour),
		fx.Logger(settings.logger),
		gologging.ProvideLogFactory(settings.logger, nil),
		app2.RegisterRootContext(pubSub),
		connectionManager.RegisterDefaultConnectionManager(),
		impl.RegisterAllConnectionRelatedServices(),
		register.ProvideMarketPriceGoogleApplication(),
		register.ProvideMarketPricesClients(),
		register.ProvideMarketPricesServices(),
		incomingMarketData.ProvideLunoMarketDataDialer(1, "tcp4://127.0.0.1:3001", pubSub),
		incomingMarketData.ProvideKrakenMarketDataDialer(pubSub, 1, "tcp4://127.0.0.1:3011"),
		//fx.Invoke(func(params struct {
		//	fx.In
		//	SheetBuilder *buildMarketDataSheet.SheetBuilder
		//	Lifecycle    fx.Lifecycle
		//}) {
		//	params.Lifecycle.Append(
		//		fx.Hook{
		//			OnStart: params.SheetBuilder.Start,
		//			OnStop: params.SheetBuilder.Stop,
		//		})
		//}),
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

	time.Sleep(time.Second)

}
