package storageadapter

import "github.com/newity/crawler/parser"

type StorageAdapter interface {
	Inject(data *parser.Data) error
	Retrieve(blocknum string) (*parser.Data, error)
}
