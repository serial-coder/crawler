package storageadapter

import "github.com/newity/crawler/parser"

type StorageAdapter interface {
	Inject(data *parser.Data) error
	Retrieve(blocknum int) (*parser.Data, error)
}
