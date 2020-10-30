package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/newity/crawler"
	"github.com/newity/crawler/blocklib"
	"github.com/sirupsen/logrus"
)

const (
	CHANNEL = "mychannel"
	USER    = "User1"
	ORG     = "Org1"
)

func main() {
	engine, err := crawler.New("connection.yaml", crawler.WithAutoConnect(USER, ORG))
	if err != nil {
		logrus.Fatal(err)
	}

	err = engine.Listen(crawler.FromBlock(), crawler.WithBlockNum(0))
	if err != nil {
		logrus.Fatal(err)
	}

	listener := engine.ListenerForChannel(CHANNEL)

	for blockevent := range listener {
		block, err := blocklib.FromFabricBlock(blockevent.Block)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.WithField("Block", block.Number()).Info("Orderer")
		for _, ordererSignature := range block.OrderersSignatures() {
			logrus.Infof("Nonce: %v\nMSP ID: %s\nHex representation of the signature hash: %s\nCertificate: %s\n",
				ordererSignature.Nonce,
				ordererSignature.MSPID,
				GetHashedHex(ordererSignature.Signature),
				string(ordererSignature.Cert))
		}

		txs, err := block.Txs()
		if err != nil {
			logrus.Fatal(err)
		}

		for index, tx := range txs {
			logrus.Info("tx ", index)
			mspidCreator, certCreator, err := tx.Creator()
			if err != nil {
				logrus.Fatal(err)
			}

			actions, err := tx.Actions()
			if err != nil {
				logrus.Fatal(err)
			}
			for i, action := range actions {
				logrus.Info("action ", i)
				endorsements := action.Endorsements()
				for _, endorsement := range endorsements {
					endorser := &msp.SerializedIdentity{}
					err = proto.Unmarshal(endorsement.Endorser, endorser)
					if err != nil {
						logrus.Fatal(err)
					}

					logrus.Infof("Creator MSP ID: %s\nCreator cert: %s\nEndorser MSP ID: %s\nEndorser cert: %s\nHex representation of the endorser's signature hash: %s\n",
						mspidCreator,
						certCreator,
						endorser.Mspid,
						string(endorser.IdBytes),
						GetHashedHex(endorsement.Signature))
				}
			}
		}
	}

	engine.StopListenAll()
}

func GetHashedHex(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
