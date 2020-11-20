/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/newity/crawler"
	"github.com/newity/crawler/storage"
	"github.com/newity/crawler/storageadapter"
	"github.com/sirupsen/logrus"
)

const (
	CHANNEL    = "mychannel"
	USER       = "User1"
	ORG        = "Org1"
	CLUSTER_ID = "testcluster"
	NATS_URL   = "nats://0.0.0.0:4222"
)

func main() {
	natsStorage, err := storage.NewNats(CLUSTER_ID, USER, NATS_URL)
	if err != nil {
		logrus.Fatal(err)
	}

	err = natsStorage.InitChannelsStorage([]string{CHANNEL})
	if err != nil {
		logrus.Fatal(err)
	}

	engine, err := crawler.New("connection.yaml", crawler.WithStorage(natsStorage), crawler.WithStorageAdapter(storageadapter.NewQueueAdapter(natsStorage)))
	if err != nil {
		logrus.Fatal(err)
	}

	err = engine.Connect(CHANNEL, USER, ORG)
	if err != nil {
		logrus.Fatal(err)
	}

	err = engine.Listen(crawler.FromBlock(), crawler.WithBlockNum(0))
	if err != nil {
		logrus.Fatal(err)
	}

	go engine.Run()

	readFromQueue(engine, CHANNEL)
	engine.StopListenAll()
}

func readFromQueue(engine *crawler.Crawler, topic string) {
	dataChan, errChan := engine.ReadStreamFromStorage(topic)
	for {
		select {
		case data := <-dataChan:
			logrus.Infof("block %d with hash %s and previous hash %s\n\nOrderers signed:\n", data.BlockNumber,
				hex.EncodeToString(data.Datahash),
				hex.EncodeToString(data.Prevhash))
			for _, signature := range data.BlockSignatures {
				fmt.Printf("MSP ID: %s\nSignature: %s\nCertificate:\n%s\n", signature.MSPID, hex.EncodeToString(signature.Signature), string(signature.Cert))
			}
			for _, tx := range data.Txs {
				t, err := tx.Timestamp()
				if err != nil {
					logrus.Error("failed to get timestamp", err)
				}
				txid, err := tx.TxId()
				if err != nil {
					logrus.Error("failed to get tx ID", err)
				}

				fmt.Printf("Tx ID: %s\nCreation time: %s\n", txid, t.String())
			}
		case err := <-errChan:
			logrus.Error(err)
		}
	}
}
