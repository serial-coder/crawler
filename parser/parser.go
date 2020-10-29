/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package parser

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/sirupsen/logrus"
	"gitlab.n-t.io/atmz/crawler/parser/blocklib"
	"io/ioutil"
	"strconv"
)

type ParserImpl struct {
}

func New() *ParserImpl {
	return &ParserImpl{}
}

func (p *ParserImpl) Parse(block *common.Block) error {
	blockBytes, err := proto.Marshal(block)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(strconv.Itoa(int(block.Header.Number)), blockBytes, 0644)

	b, err := blocklib.FromFabricBlock(block)
	if err != nil {
		return err
	}
	txs, err := b.GetTxs()
	if err != nil {
		return err
	}
	for _, tx := range txs {
		chdr, err := tx.ChannelHeader()
		if err != nil {
			logrus.Errorf("tx parser error: %s", err)
			continue
		}

		if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

			actions, err := tx.Actions()
			if err != nil {
				logrus.Errorf("failed to actions from transaction: %s", err)
				continue
			}
			for _, action := range actions {

				ccActionPayload := action.ChaincodeActionPayload()
				if ccActionPayload.Action == nil || ccActionPayload.Action.ProposalResponsePayload == nil {
					logrus.Debug("no payload in ChaincodeActionPayload")
					continue
				}

				ccAction, err := action.ChaincodeAction()
				if err != nil {
					logrus.Errorf("failed to get to ChaincodeAction: %s", err)
					continue
				}
				_ = ccAction
				rwsets, err := action.RWSets()
				if err != nil {
					logrus.Errorf("failed to extract rwsets: %+v", err)
					continue
				}

				for _, rw := range rwsets {
					for _, write := range rw.KVRWSet.Writes {
						_ = write
					}
				}

				ccEvent, err := action.ChaincodeEvent()
				if err != nil {
					logrus.Errorf("failed to extract chaincode events: %s", err)
					continue
				}
				_ = ccEvent

			}
		}
	}
	return nil
}
