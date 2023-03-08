package incomingMarketData

import (
	"context"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	stream3 "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type reactor struct {
	common.BaseConnectionReactor
	MessageRouter                     *messageRouter.MessageRouter
	FullMarketDataHelper              fullMarketDataHelper.IFullMarketDataHelper
	FmdService                        fullMarketDataManagerService.IFmdManagerService
	externalFullMarketDataInstruments []string
}

func (self *reactor) Close() error {
	var err error
	fmdMessages := make([]interface{}, len(self.externalFullMarketDataInstruments))
	for i, instrument := range self.externalFullMarketDataInstruments {
		fmdMessages[i] = &stream2.FullMarketData_RemoveInstrumentInstruction{
			FeedName:   self.UniqueReference,
			Instrument: instrument,
		}
	}
	self.FmdService.MultiSend(fmdMessages...)
	return multierr.Append(err, self.BaseConnectionReactor.Close())
}

func (self *reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}

	self.OnSendToConnection(&stream2.FullMarketData_InstrumentList_Subscribe{})
	self.OnSendToConnection(&stream2.FullMarketData_InstrumentList_Request{})

	return nil
}

func (self *reactor) Init(
	onSendToReactor rxgo.NextFunc,
	onSendToConnection rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(
		onSendToReactor,
		onSendToConnection,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return func(i interface{}) {
			self.MessageRouter.Route(i)
		},
		func(err error) {
			self.MessageRouter.Route(err)
		},
		func() {

		},
		nil
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_InstrumentList_ResponseWrapper(incomingMessage *stream2.FullMarketData_InstrumentList_ResponseWrapper) {
	incomingMessage.Data.FeedName = self.UniqueReference
	registeredSources := make([]string, len(self.externalFullMarketDataInstruments), len(self.externalFullMarketDataInstruments))
	for i, instrument := range self.externalFullMarketDataInstruments {
		registeredSources[i] = self.FullMarketDataHelper.RegisteredSource(instrument)
	}

	// must check len here, otherwise we will deregister the whole self.OnSendToConnectionPubSubBag
	if len(registeredSources) > 0 {
		self.PubSub.Unsub(self.OnSendToConnectionPubSubBag, registeredSources...)
	}

	self.externalFullMarketDataInstruments = make([]string, len(incomingMessage.Data.Instruments), len(incomingMessage.Data.Instruments))
	for i, s := range incomingMessage.Data.Instruments {
		self.externalFullMarketDataInstruments[i] = s.Instrument
	}
	registeredSources = make([]string, len(incomingMessage.Data.Instruments), len(incomingMessage.Data.Instruments))
	for i, instrument := range self.externalFullMarketDataInstruments {
		registeredSources[i] = self.FullMarketDataHelper.RegisteredSource(instrument)
	}
	self.PubSub.AddSub(self.OnSendToConnectionPubSubBag, registeredSources...)
	_ = self.FmdService.Send(incomingMessage)
}

func (self *reactor) handlePublishRxHandlerCounters(incomingMessage *model.PublishRxHandlerCounters) {
	// not used. Swallowing message
}

func (self *reactor) handleEmptyQueue(incomingMessage *messages.EmptyQueue) {
	// not used. Swallowing message
}

func (self *reactor) handlePingWrapper(incomingMessage *stream3.PingWrapper) {
	// not used. Swallowing message
}

func (self *reactor) handlePongWrapper(incomingMessage *stream3.PongWrapper) {
}

func (self *reactor) OnUnknown(i interface{}) {
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_AddOrderInstructionWrapper(msg *stream2.FullMarketData_AddOrderInstructionWrapper) {
	msg.Data.FeedName = self.UniqueReference
	_ = self.FmdService.Send(msg)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_ClearWrapper(msg *stream2.FullMarketData_ClearWrapper) {
	msg.Data.FeedName = self.UniqueReference
	_ = self.FmdService.Send(msg)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_ReduceVolumeInstructionWrapper(msg *stream2.FullMarketData_ReduceVolumeInstructionWrapper) {
	msg.Data.FeedName = self.UniqueReference
	_ = self.FmdService.Send(msg)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_DeleteOrderInstructionWrapper(msg *stream2.FullMarketData_DeleteOrderInstructionWrapper) {
	msg.Data.FeedName = self.UniqueReference
	_ = self.FmdService.Send(msg)
}

func NewConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	PubSub *pubsub.PubSub,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	GoFunctionCounter GoFunctionCounter.IService,
	FmdService fullMarketDataManagerService.IFmdManagerService,
) *reactor {
	result := &reactor{
		BaseConnectionReactor: common.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			UniqueReferenceService.Next("ConnectionReactor"),
			PubSub,
			GoFunctionCounter,
		),
		MessageRouter:        messageRouter.NewMessageRouter(),
		FullMarketDataHelper: FullMarketDataHelper,
		FmdService:           FmdService,
	}
	result.MessageRouter.RegisterUnknown(result.OnUnknown)
	_ = result.MessageRouter.Add(result.handlePongWrapper)
	_ = result.MessageRouter.Add(result.handlePingWrapper)
	_ = result.MessageRouter.Add(result.handleEmptyQueue)
	_ = result.MessageRouter.Add(result.handlePublishRxHandlerCounters)

	_ = result.MessageRouter.Add(result.handleFullMarketData_InstrumentList_ResponseWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_AddOrderInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ClearWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ReduceVolumeInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_DeleteOrderInstructionWrapper)

	return result
}
