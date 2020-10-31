package parser

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/newity/crawler/blocklib"
)

type Data struct {
	Block           blocklib.Block
	BlockSignatures []blocklib.BlockSignature
	Txs             []blocklib.Tx
	Events          []*peer.ChaincodeEvent
}
