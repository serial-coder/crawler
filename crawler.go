/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Please note, we don't return any errors from constructor (New()) or initialization funcs (Connect(...), Listen(...)).
// This is needed for a smooth developer experience.
//
// This is what the initialization of the crawler might look like:
//	if err := crawler.New("connection.yaml").Connect("User1", "Org1").Listen(crawler.FromBlock(), crawler.WithBlockNum(0)); err != nil {
//		fmt.Errorf(err.Error())
//	}
//
// For connection to all channels from connection profile:
//	if err := crawler.New("connection.yaml", crawler.WithAutoConnect("User1", "Org1")).Listen(crawler.FromBlock(), crawler.WithBlockNum(0)); err != nil {
//		fmt.Errorf(err.Error())
//	}
//
// If something bad happens during initialization, the Crawler will panic and recover.
package crawler

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/newity/crawler/parser"
	"github.com/sirupsen/logrus"
	"reflect"
)

// Crawler is responsible for fetching info from blockchain
type Crawler struct {
	sdk             *fabsdk.FabricSDK
	chCli           map[string]*channel.Client
	eventCli        map[string]*event.Client
	channelProvider contextApi.ChannelProvider
	notifiers       map[string]<-chan *fab.BlockEvent
	parser          parser.Parser
}

// New creates Crawler instance from HLF connection profile and returns pointer to it.
// "connectionProfile" is a path to HLF connection profile
func New(connectionProfile string, opts ...Option) *Crawler {
	go func() {
		if r := recover(); r != nil {
			logrus.Error("Recovered", r)
		}
	}()

	sdk, err := fabsdk.New(config.FromFile(connectionProfile))
	if err != nil {
		panic(err)
	}

	crawl := &Crawler{
		sdk:       sdk,
		chCli:     make(map[string]*channel.Client),
		eventCli:  make(map[string]*event.Client),
		notifiers: make(map[string]<-chan *fab.BlockEvent),
	}

	for _, opt := range opts {
		if err = opt(crawl); err != nil {
			panic(err)
		}
	}

	// if no parser is specified, use the default parser ParserImpl
	if crawl.parser == nil {
		crawl.parser = parser.New()
	}

	return crawl
}

// Connect connects crawler to channel 'ch' as identity specified in 'username' from organization with name 'org'
func (c *Crawler) Connect(ch, username, org string) *Crawler {
	var err error
	c.channelProvider = c.sdk.ChannelContext(ch, fabsdk.WithUser(username), fabsdk.WithOrg(org))
	c.chCli[ch], err = channel.New(c.channelProvider)
	if err != nil {
		panic(err)
	}
	return c
}

// Listen starts blocks listener starting from block with num 'from'
func (c *Crawler) Listen(opts ...ListenOpt) error {
	var (
		err        error
		listenType int
		fromBlock  uint64
		clientOpts []event.ClientOption
	)

	for _, opt := range opts {
		if reflect.TypeOf(opt).Name() == "WithBlockNum" {
			fromBlock = uint64(opt())
		} else {
			listenType = opt()
		}
	}

	switch listenType {
	case LISTEN_FROM:
		clientOpts = append(clientOpts,
			event.WithBlockEvents(),
			event.WithSeekType(seek.FromBlock),
			event.WithBlockNum(fromBlock),
		)
	case LISTEN_NEWEST:
		clientOpts = append(clientOpts,
			event.WithBlockEvents(),
			event.WithSeekType(seek.Newest),
		)
	case LISTEN_OLDEST:
		clientOpts = append(clientOpts,
			event.WithBlockEvents(),
			event.WithSeekType(seek.Oldest),
		)
	}

	for ch := range c.chCli {
		c.eventCli[ch], err = event.New(
			c.channelProvider,
			clientOpts...,
		)
		if err != nil {
			return err
		}
		_, c.notifiers[ch], err = c.eventCli[ch].RegisterBlockEvent()
		return err
	}
	return err
}

//func (c *Crawler) (){}
//func (c *Crawler) (){}
//func (c *Crawler) (){}
//func (c *Crawler) (){}
//func (c *Crawler) (){}
//func (c *Crawler) (){}
//func (c *Crawler) (){}
