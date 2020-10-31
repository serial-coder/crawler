/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package storage

import (
	"bytes"
	"encoding/gob"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/newity/crawler/parser"
	"os"
	"strconv"
	"strings"
)

type Badger struct {
	db *badger.DB
}

func NewBadger(path string) (Storage, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		if strings.Contains(err.Error(), "Another process is using this Badger database") {
			db, err = badger.Open(badger.DefaultOptions(path + strconv.Itoa(os.Getpid())))
		} else {
			return nil, err
		}
	}

	return &Badger{db}, nil
}

func (b *Badger) InitChannelStorage(channel string) error {
	return nil
}

// Put saves parser.Data using block number as key
func (b *Badger) Put(data *parser.Data) error {
	encoded, err := Encode(data)
	if err != nil {
		return err
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(strconv.Itoa(int(data.BlockNumber))), encoded)
	})
}

// Get retrieves parser.Data from BadgerDB using block number as key
func (b *Badger) Get(key string) (*parser.Data, error) {
	var value []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		return err
	})
	decoded, err := Decode(value)
	if err != nil {
		return nil, err
	}
	return decoded, err
}

func (b *Badger) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *Badger) Close() error {
	return b.db.Close()
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
