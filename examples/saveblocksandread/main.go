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
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const (
	CHANNEL = "cc"
	USER    = "User1"
	ORG     = "atomyze"
)

func main() {
	home := os.Getenv("HOME")
	stor, err := storage.NewBadger(path.Join(home, ".crawler-storage"))
	if err != nil {
		logrus.Fatal(err)
	}

	engine, err := crawler.New("connection.yaml", crawler.WithAutoConnect(USER, ORG), crawler.WithStorage(stor))
	if err != nil {
		logrus.Error(err)
	}

	err = engine.Listen(crawler.FromBlock(), crawler.WithBlockNum(0))
	if err != nil {
		logrus.Error(err)
	}

	go engine.Run()

	// wait for pulling blocks to storage
	time.Sleep(1 * time.Second)

	readBlock(engine, 2)
	engine.StopListenAll()
}

func readBlock(engine *crawler.Crawler, num int) {
	data, err := engine.GetBlock(num)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("block %d with hash %s and prebious hash %s\n\nOrderers signed:\n", data.BlockNumber,
		hex.EncodeToString(data.Datahash),
		hex.EncodeToString(data.Prevhash))
	for _, signature := range data.BlockSignatures {
		fmt.Printf("MSP ID: %s\nSignature: %s\nCertificate:\n%s\n", signature.MSPID, hex.EncodeToString(signature.Signature), string(signature.Cert))
	}
	fmt.Println("Transactions")
	for _, tx := range data.Txs {
		t, err := tx.Timestamp()
		if err != nil {
			logrus.Error(err)
		}
		txid, err := tx.TxId()
		if err != nil {
			logrus.Error(err)
		}

		fmt.Printf("Tx ID: %s\nCreation time: %s\n", txid, t.String())
	}
}
