package message-broker

type CacheBroker interface {
	BrokerWithClient
	Remove(key interface{}) error
	Purge() error
}