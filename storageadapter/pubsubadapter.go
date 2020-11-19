/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package storageadapter

import (
	"github.com/newity/crawler/parser"
	"github.com/newity/crawler/storage"
	"sync"
)

// PubSubAdapter is a storage adapter for the Google Pub/Sub
type PubSubAdapter struct {
	storage storage.Storage
}

func NewPubSubAdapter(stor storage.Storage) *PubSubAdapter {
	return &PubSubAdapter{stor}
}

func (s *PubSubAdapter) Inject(data *parser.Data) error {
	encoded, err := Encode(data)
	if err != nil {
		return err
	}
	return s.storage.Put(data.Channel, encoded)
}

func (s *PubSubAdapter) Retrieve(topic string) (*parser.Data, error) {
	value, err := s.storage.Get(topic)
	if err != nil {
		return nil, err
	}
	return Decode(value)
}

func (s *PubSubAdapter) ReadStream(topic string) (<-chan *parser.Data, <-chan error) {
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
	return out, errOutChan
}
