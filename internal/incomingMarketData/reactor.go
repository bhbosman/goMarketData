package incomingMarketData

import (
	"context"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type reactor struct {
	common.BaseConnectionReactor
	MessageRouter                     messageRouter.IMessageRouter
	FullMarketDataHelper              fullMarketDataHelper.IFullMarketDataHelper
	FmdService                        fullMarketDataManagerService.IFmdManagerService
	externalFullMarketDataInstruments map[string]bool
}

func (self *reactor) Close() error {
	var err error

	fmdMessages := make([]interface{}, 0, len(self.externalFullMarketDataInstruments))
	for key := range self.externalFullMarketDataInstruments {
		fmdMessages = append(
			fmdMessages,

			&stream2.FullMarketData_RemoveInstrumentInstruction{
				FeedName:   self.UniqueReference,
				Instrument: key,
			},
		)

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

func (self *reactor) Init(params intf.IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(params)
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
	removeInstruments := make(map[string]bool)
	for s := range self.externalFullMarketDataInstruments {
		removeInstruments[s] = true
	}

	addInstruments := make(map[string]bool)
	for _, s := range incomingMessage.Data.Instruments {
		if _, ok := self.externalFullMarketDataInstruments[s.Instrument]; ok {
			delete(removeInstruments, s.Instrument)
		} else {
			addInstruments[s.Instrument] = true
		}
	}
	if len(removeInstruments) > 0 {
		l := make([]string, 0, len(removeInstruments))
		for s := range removeInstruments {
			l = append(l, self.FullMarketDataHelper.RegisteredSource(s))
			self.PubSub.Unsub(self.OnSendToConnectionPubSubBag, l...)
		}
	}
	if len(addInstruments) > 0 {
		l := make([]string, 0, len(addInstruments))
		for s := range addInstruments {
			l = append(l, self.FullMarketDataHelper.RegisteredSource(s))
			self.externalFullMarketDataInstruments[s] = true
		}
		self.PubSub.AddSub(self.OnSendToConnectionPubSubBag, l...)
	}

	_ = self.FmdService.Send(incomingMessage)
}

func (self *reactor) OnUnknown(_ interface{}) {
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
) intf.IConnectionReactor {
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
		MessageRouter:                     messageRouter.NewMessageRouter(),
		FullMarketDataHelper:              FullMarketDataHelper,
		FmdService:                        FmdService,
		externalFullMarketDataInstruments: make(map[string]bool),
	}
	result.MessageRouter.RegisterUnknown(result.OnUnknown)
	_ = result.MessageRouter.Add(result.handleFullMarketData_InstrumentList_ResponseWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_AddOrderInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ClearWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ReduceVolumeInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_DeleteOrderInstructionWrapper)

	return result
}
