package incomingMarketData

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker"
	"github.com/bhbosman/goCommsStacks/messageCompressor"
	"github.com/bhbosman/goCommsStacks/messageNumber"
	"github.com/bhbosman/goCommsStacks/pingPong"
	"github.com/bhbosman/goCommsStacks/protoBuf"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"net/url"
)

func ProvideKrakenMarketDataDialer(
	maxConnections int,
	urlAsText string) fx.Option {
	const LunoMarketData = "KrakenMarketData"
	var opt []fx.Option
	opt = append(
		opt,
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(
					params struct {
						fx.In
						PubSub               *pubsub.PubSub `name:"Application"`
						NetAppFuncInParams   common.NetAppFuncInParams
						FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
						FmdService           fullMarketDataManagerService.IFmdManagerService
					},
				) (messages.CreateAppCallback, error) {

					dialerUrl, err := url.Parse(urlAsText)
					if err != nil {
						return messages.CreateAppCallback{}, err
					}

					f := goCommsNetDialer.NewSingleNetDialApp(
						LunoMarketData,
						common.MaxConnectionsSetting(maxConnections),
						common.MoreOptions(
							goCommsDefinitions.ProvideUrl("ConnectionUrl", dialerUrl),
							goCommsDefinitions.ProvideUrl("ProxyUrl", nil),
							goCommsDefinitions.ProvideBool("UseProxy", false),
						),
						common.NewConnectionInstanceOptions(
							PubSub.ProvidePubSubInstance("Application", params.PubSub),
							ProvideConnectionReactor(),
							goCommsDefinitions.ProvideTransportFactoryForCompressedName(
								topStack.ProvideTopStack(),
								pingPong.ProvidePingPongStacks(),
								protoBuf.ProvideStack(),
								messageCompressor.Provide(),
								messageNumber.ProvideMessageNumberStack(),
								bvisMessageBreaker.Provide(),
								bottom.Provide(),
							),
							fx.Provide(
								fx.Annotated{
									Target: func() (fullMarketDataHelper.IFullMarketDataHelper, fullMarketDataManagerService.IFmdManagerService) {
										return params.FullMarketDataHelper, params.FmdService
									},
								},
							),
						),
					)
					return f(params.NetAppFuncInParams), nil
				},
			},
		),
	)

	return fx.Options(opt...)
}
