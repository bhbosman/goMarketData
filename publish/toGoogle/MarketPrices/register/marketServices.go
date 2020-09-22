package register

import (
	"context"
	"github.com/bhbosman/goMarketData/publish/toGoogle/MarketPrices/service"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func ProvideMarketPricesServices() fx.Option {
	return fx.Options(
		fx.Provide(
			func(params struct {
				fx.In
				MarketPricesConfig *oauth2.Config     `name:"MarketPricesApplication"`
				NamedHttpClient    []*NamedHttpClient `group:"MarketPricesClients"`
			}) ([]*service.Service, error) {
				var services []*service.Service
				for _, namedHttpClient := range params.NamedHttpClient {
					drive, err := drive.NewService(context.Background(), option.WithHTTPClient(namedHttpClient.Client))
					if err != nil {
						return nil, err
					}
					sheets, err := sheets.NewService(context.Background(), option.WithHTTPClient(namedHttpClient.Client))
					if err != nil {
						return nil, err
					}
					service := service.NewService(
						namedHttpClient.Name,
						drive,
						sheets)
					services = append(services, service)
				}
				return services, nil
			}),
		fx.Invoke(
			func(params struct {
				fx.In
				Lifecycle fx.Lifecycle
				Services  []*service.Service
			}) {
				for _, s := range params.Services {
					service := s
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return service.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							return service.Stop(ctx)
						},
					})
				}

			}))
}
