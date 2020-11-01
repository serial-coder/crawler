package storageadapter

import (
	"bytes"
	"encoding/gob"
	"github.com/newity/crawler/parser"
	"github.com/newity/crawler/storage"
	"strconv"
)

type SimpleAdapter struct {
	storage storage.Storage
}

func NewSimpleAdapter(stor storage.Storage) *SimpleAdapter {
	return &SimpleAdapter{stor}
}

func (s *SimpleAdapter) Inject(data *parser.Data) error {
	encoded, err := Encode(data)
	if err != nil {
		return err
	}
	return s.storage.Put(strconv.Itoa(int(data.BlockNumber)), encoded)
}

func (s *SimpleAdapter) Retrieve(blocknum int) (*parser.Data, error) {
	value, err := s.storage.Get(strconv.Itoa(blocknum))
	if err != nil {
		return nil, err
	}
	return Decode(value)
}

func Encode(data *parser.Data) ([]byte, error) {
	var bytebuffer bytes.Buffer
	e := gob.NewEncoder(&bytebuffer)
	if err := e.Encode(data); err != nil {
		return nil, err
	}
	return bytebuffer.Bytes(), nil
}

func Decode(data []byte) (*parser.Data, error) {
	decoded := &parser.Data{}
	bytebuffer := bytes.NewBuffer(data)
	d := gob.NewDecoder(bytebuffer)
	if err := d.Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}
