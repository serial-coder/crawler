/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package crawler

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/newity/crawler/parser"
	"github.com/newity/crawler/storage"
	"github.com/newity/crawler/storageadapter"
)

type Option func(crawler *Crawler) error

// WithAutoConnect connects crawler to all channels specified in connection profile
// 'username' is a Fabric identity name and 'org' is a Fabric organization ti which the identity belongs
func WithAutoConnect(username, org string) Option {
	return func(crawler *Crawler) error {
		configBackend, err := crawler.sdk.Config()
		if err != nil {
			return err
		}
		if channelsIface, ok := configBackend.Lookup("channels"); !ok {
			return fmt.Errorf("failed to find channels in connection profile")
		} else {
			channelsMap, ok := channelsIface.(map[string]interface{})
			if !ok {
				return fmt.Errorf("failed to parse connection profile")
			}

			var err error
			for ch := range channelsMap {
				crawler.channelProvider = crawler.sdk.ChannelContext(ch, fabsdk.WithUser(username), fabsdk.WithOrg(org))
				crawler.chCli[ch], err = channel.New(crawler.channelProvider)
				return err
			}
		}
		return nil
	}
}

// WithParser injects a specific parser that satisfies the Parser interface to the Crawler instance.
// If no parser is specified, the default parser ParserImpl will be used.
func WithParser(p parser.Parser) Option {
	return func(crawler *Crawler) error {
		crawler.parser = p
		return nil
	}
}

// WithStorage adds a specific storage that satisfies the Storage interface to the Crawler instance.
// If no storage is specified, the default storage Badger (BadgerDB) will be used.
func WithStorage(s storage.Storage) Option {
	return func(crawler *Crawler) error {
		crawler.storage = s
		return nil
	}
}

// WithStorageAdapter adds storage adapter that satisfies the StorageAdapter interface to the Crawler instance.
// If no storage adapter is specified, the default SimpleAdapter will be used.
func WithStorageAdapter(a storageadapter.StorageAdapter) Option {
	return func(crawler *Crawler) error {
		crawler.adapter = a
		return nil
	}
}

type ListenOpt func() interface{}

const (
	LISTEN_FROM   = "from"
	LISTEN_NEWEST = "newest"
	LISTEN_OLDEST = "oldest"
)

func FromBlock() ListenOpt {
	return func() interface{} {
		return LISTEN_FROM
	}
}

func WithBlockNum(block uint64) ListenOpt {
	return func() interface{} {
		return int(block)
	}
}

func Newest() ListenOpt {
	return func() interface{} {
		return LISTEN_NEWEST
	}
}

func Oldest() ListenOpt {
	return func() interface{} {
		return LISTEN_OLDEST
	}
}
