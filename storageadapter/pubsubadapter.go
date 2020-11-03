package storageadapter

import (
	"github.com/newity/crawler/parser"
	"github.com/newity/crawler/storage"
)

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
