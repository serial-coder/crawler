package storageadapter

import (
	"context"
	"github.com/newity/crawler/parser"
	"github.com/newity/crawler/storage"
	"sync"
)

// QueueAdapter is a general adapter for message brokers
type QueueAdapter struct {
	storage storage.Storage
}

func NewQueueAdapter(stor storage.Storage) *QueueAdapter {
	return &QueueAdapter{stor}
}

func (s *QueueAdapter) Inject(data *parser.Data) error {
	encoded, err := Encode(data)
	if err != nil {
		return err
	}
	return s.storage.Put(data.Channel, encoded)
}

func (s *QueueAdapter) Retrieve(topic string) (*parser.Data, error) {
	value, err := s.storage.Get(topic)
	if err != nil {
		return nil, err
	}
	return Decode(value)
}

func (s *QueueAdapter) ReadStream(topic string) (<-chan *parser.Data, <-chan error, context.CancelFunc) {
	stream, errChan := s.storage.GetStream(topic)
	var out, errOutChan = make(chan *parser.Data), make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		for {
			select {
			case msg := <-stream:
				decodedMsg, err := Decode(msg)
				if err != nil {
					errOutChan <- err
				}
				out <- decodedMsg
			case err := <-errChan:
				errOutChan <- err
			}
		}
	}()
	wg.Wait()
	return out, errOutChan, cancel
}
