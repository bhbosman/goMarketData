package googleSheets

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
	"time"
)

type service struct {
	ctx               context.Context
	cancelFunc        context.CancelFunc
	cmdChannel        chan interface{}
	onData            func() (IGoogleSheetsData, error)
	Logger            *zap.Logger
	state             IFxService.State
	pubSub            *pubsub.PubSub
	goFunctionCounter GoFunctionCounter.IService
	subscribeChannel  *pubsub.NextFuncSubscription
	fmdServiceHelper  fullMarketDataHelper.IFullMarketDataHelper
	FmdManagerService fullMarketDataManagerService.IFmdManagerService
	clientDetails     clientDetails
}

func (self *service) MultiSend(messages ...interface{}) {
	_, err := CallIGoogleSheetsMultiSend(self.ctx, self.cmdChannel, false, messages...)
	if err != nil {
		return
	}
}

func (self *service) Send(message interface{}) error {
	send, err := CallIGoogleSheetsSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) start(_ context.Context) error {
	instanceData, err := self.onData()
	if err != nil {
		return err
	}

	return self.goFunctionCounter.GoRun(
		"FMD Service",
		func() {
			self.goStart(instanceData)
		},
	)
}

func (self *service) goStart(instanceData IGoogleSheetsData) {
	self.subscribeChannel = pubsub.NewNextFuncSubscription(goCommsDefinitions.CreateNextFunc(self.cmdChannel))
	self.pubSub.AddSub(self.subscribeChannel, self.fmdServiceHelper.InstrumentListChannelName())
	_ = self.FmdManagerService.Send(fullMarketDataManagerService.NewRequestAllInstruments(self.subscribeChannel.Add))
	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IGoogleSheets); ok {
						return ChannelEventsForIGoogleSheets(unk, message)
					}
					return false, nil

				},
			},
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return true, unk.Send(message)
					}
					return false, nil
				},
			},
			// TODO: add handlers here
		},
		func() int {
			return len(self.cmdChannel)
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
	interval := time.Second
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()

	_ = instanceData.Start(context.Background())
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.Logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		case publishTime, ok := <-ticker.C:
			if !ok {
				break loop
			}
			ticker.Stop()
			breakLoop, err := channelHandlerCallback(
				&publishData{
					publishTime: publishTime,
				},
			)
			if err != nil || breakLoop {
				break loop
			}
			ticker.Reset(interval)
		}
	}
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self service) ServiceName() string {
	return "GoogleSheets"
}

func newService(
	parentContext context.Context,
	onData func() (IGoogleSheetsData, error),
	logger *zap.Logger,
	pubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	fmdServiceHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdManagerService fullMarketDataManagerService.IFmdManagerService,
	clientDetails clientDetails,
) (IGoogleSheetsService, error) {
	localCtx, localCancelFunc := context.WithCancel(parentContext)
	return &service{
		ctx:               localCtx,
		cancelFunc:        localCancelFunc,
		cmdChannel:        make(chan interface{}, 32),
		onData:            onData,
		Logger:            logger,
		pubSub:            pubSub,
		goFunctionCounter: goFunctionCounter,
		fmdServiceHelper:  fmdServiceHelper,
		FmdManagerService: FmdManagerService,
		clientDetails:     clientDetails,
	}, nil
}

type appWrapper struct {
	err                 error
	googleSheetsService IGoogleSheetsService
}

func newAppWrapper(googleSheetsService IGoogleSheetsService) *appWrapper {
	return &appWrapper{
		googleSheetsService: googleSheetsService,
	}
}

func (self *appWrapper) Start(ctx context.Context) error {
	self.err = self.googleSheetsService.OnStart(ctx)
	return self.err
}

func (self *appWrapper) Stop(ctx context.Context) error {
	return self.googleSheetsService.OnStop(ctx)
}

func (self *appWrapper) Err() error {
	return self.err
}
