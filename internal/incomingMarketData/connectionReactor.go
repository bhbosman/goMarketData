package incomingMarketData

import (
	"context"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"net"
	"net/url"
)

type ConnectionReactor struct {
	impl.BaseConnectionReactor
	MessageRouter *messageRouter.MessageRouter
	PubSub *pubsub.PubSub
}

func NewConnectionReactor(
	logger *gologging.SubSystemLogger,
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	PubSub *pubsub.PubSub,
	userContext interface{}) *ConnectionReactor {
	connectionReactor := &ConnectionReactor{
		BaseConnectionReactor: impl.NewBaseConnectionReactor(logger, name, cancelCtx, cancelFunc, userContext),
		MessageRouter:         messageRouter.NewMessageRouter(),
		PubSub:                PubSub,
	}
	_ = connectionReactor.MessageRouter.Add(connectionReactor.handleRws)
	_ = connectionReactor.MessageRouter.Add(connectionReactor.handlePublishTop5Wrapper)

	return connectionReactor
}

func (self *ConnectionReactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, err := self.BaseConnectionReactor.Init(conn, url, connectionId, connectionManager, toConnectionFunc, toConnectionReactor)
	if err != nil {
		return nil, err
	}
	return self.doNext, nil
}

func (self *ConnectionReactor) doNext(external bool, i interface{}) {
	_, _ = self.MessageRouter.Route(i)
}

func (self ConnectionReactor) handlePublishTop5Wrapper(incomingMessage *marketDataStream.PublishTop5Wrapper) error {
	self.PubSub.TryPub(incomingMessage.Data, "Top5Data")
	return nil
}

func (self ConnectionReactor) handleRws(incomingMessage *gomessageblock.ReaderWriter) error {
	msg, err := stream.UnMarshal(incomingMessage, self.CancelCtx, self.CancelFunc, self.ToReactor, self.ToConnection)
	if err != nil {
		return err
	}
	_, _ = self.MessageRouter.Route(msg)
	return nil
}
