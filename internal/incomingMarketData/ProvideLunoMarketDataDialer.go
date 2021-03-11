package incomingMarketData

import (
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvideLunoMarketDataDialer(maxConnections int, url string, PubSub *pubsub.PubSub) fx.Option {
	const LunoMarketData = "LunoMarketData"
	cfr := NewConnectionReactorFactory(LunoMarketData, PubSub)
	var opt []fx.Option
	opt = append(
		opt,
		fx.Provide(fx.Annotated{
			Group: impl.ConnectionReactorFactoryConst,
			Target: func(
				params struct {
					fx.In
					PubSub *pubsub.PubSub `name:"Application"`
				}) (intf.IConnectionReactorFactory, error) {

				return cfr, nil

			},
		}))

	opt = append(
		opt,
		fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: netDial.NewNetDialApp(
				LunoMarketData,
				url,
				impl.CreateCompressedStack,
				LunoMarketData,
				cfr,
				netDial.MaxConnectionsSetting(maxConnections)),
		}))

	return fx.Options(opt...)
}
