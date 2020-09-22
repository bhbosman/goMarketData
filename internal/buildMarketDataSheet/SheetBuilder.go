package buildMarketDataSheet

import (
	"context"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"time"
)

type SheetBuilder struct {
	PubSub        *pubsub.PubSub
	cancelContext context.Context
	cancelFunc    context.CancelFunc
	MessageRouter *messageRouter.MessageRouter
}

func (self *SheetBuilder) Start(ctx context.Context) error {
	channel := self.PubSub.Sub("Top5Data")
	go self.start(channel)
	return nil
}

func (self *SheetBuilder) Stop(ctx context.Context) error {
	self.cancelFunc()
	return nil
}

func (self *SheetBuilder) start(channel chan interface{}) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
loop:
	for true {
		select {
		case <-self.cancelContext.Done():
			break loop
		case v, ok := <-channel:
			if ok {
				_, _ = self.MessageRouter.Route(v)
				if len(channel) == 0 {
					_, _ = self.MessageRouter.Route(&rxgo.EmptyQueue{})
				}
				continue loop
			}
			break loop
		}
	}
	// unsubscribe
	self.PubSub.Unsub(channel)
	// flush to avoid deadlocks
	for range channel {
	}
}

func (self *SheetBuilder) handlePublishTop(incomingMessage *marketDataStream.PublishTop5) error {
	return nil
}

func (self *SheetBuilder) handleEmptyQueue(incomingMessage *rxgo.EmptyQueue) error {
	return nil
}

func NewSheetBuilder(
	params struct {
		fx.In
		PubSub *pubsub.PubSub `name:"Application"`
	}) (*SheetBuilder, error) {
	cancel, cancelFunc := context.WithCancel(context.Background())
	result := &SheetBuilder{
		PubSub:        params.PubSub,
		cancelContext: cancel,
		cancelFunc:    cancelFunc,
		MessageRouter: messageRouter.NewMessageRouter(),
	}
	_ = result.MessageRouter.Add(result.handlePublishTop)
	_ = result.MessageRouter.Add(result.handleEmptyQueue)
	return result, nil
}
