package injector

import (
	"github.com/newity/crawler/parser"
	"github.com/newity/crawler/storage"
)

type SimpleInjector struct {
	storage storage.Storage
}

func NewSimpleInjector(stor storage.Storage) *SimpleInjector {
	return &SimpleInjector{stor}
}

func (s *SimpleInjector) Inject(data *parser.Data) error {
	return s.storage.Put(data)
}
