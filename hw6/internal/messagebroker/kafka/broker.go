package kafka

import (
	"Go/hw6/internal/messagebroker"
	"context"

	lru "github.com/hashicorp/golang-lru"
)

type Broker struct {
	brokers  []string 
	clientID string   
	
	cacheBroker messagebroker.CacheBroker 
	cache       *lru.TwoQueueCache        
}

func NewBroker(brokers []string, cache *lru.TwoQueueCache, clientID string) messagebroker.MessageBroker {
	return &Broker{brokers: brokers, cache: cache, clientID: clientID}
}

func (b *Broker) Connect(ctx context.Context) error {
	brokers := []messagebroker.BrokerWithClient{b.Cache()}

	for _, broker := range brokers {
		if err := broker.Connect(ctx, b.brokers); err != nil {
			return err
		}
	}

	return nil
}

func (b *Broker) Close() error {
	brokers := []messagebroker.BrokerWithClient{b.Cache()}
	for _, broker := range brokers {
		if err := broker.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Broker) Cache() messagebroker.CacheBroker {
	if b.cacheBroker == nil {
		b.cacheBroker = NewCacheBroker(b.cache, b.clientID)
	}
	return b.cacheBroker
}