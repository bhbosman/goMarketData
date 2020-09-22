package incomingMarketData

import (
	"context"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
)

type ConnectionReactorFactory struct {
	name string
	PubSub *pubsub.PubSub
}

func NewConnectionReactorFactory(name string, PubSub *pubsub.PubSub) *ConnectionReactorFactory {
	return &ConnectionReactorFactory{
		name:   name,
		PubSub: PubSub,
	}
}

func (self *ConnectionReactorFactory) Name() string {
	return self.name
}

func (self *ConnectionReactorFactory) Create(name string, cancelCtx context.Context, cancelFunc context.CancelFunc, logger *gologging.SubSystemLogger, userContext interface{}) intf.IConnectionReactor {
	return NewConnectionReactor(logger, name, cancelCtx, cancelFunc, self.PubSub, userContext)
}
