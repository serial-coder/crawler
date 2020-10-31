package storage

import "github.com/newity/crawler/parser"

// Storage interface is a contract for storage implementations
type Storage interface {
	// init storage (initial setup of storage and connection create operations)
	InitChannelStorage(channel string) error
	// put parser.Data to storage
	Put(data *parser.Data) error
	// get parser.Data from storage by specified key
	Get(key string) (*parser.Data, error)
	// remove parser.Data from storage
	Delete(key string) error
	// close connection to storage (network connection or file descriptor)
	Close() error
}
