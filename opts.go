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
			for ch, _ := range channelsMap {
				crawler.chCli[ch], err = channel.New(crawler.sdk.ChannelContext(ch, fabsdk.WithUser(username), fabsdk.WithOrg(org)))
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

type ListenOpt func() int

const (
	LISTEN_FROM = iota
	LISTEN_NEWEST
	LISTEN_OLDEST
)

func FromBlock() ListenOpt {
	return func() int {
		return LISTEN_FROM
	}
}

func WithBlockNum(block uint64) ListenOpt {
	return func() int {
		return int(block)
	}
}

func Newest() ListenOpt {
	return func() int {
		return LISTEN_NEWEST
	}
}

func Oldest() ListenOpt {
	return func() int {
		return LISTEN_OLDEST
	}
}
