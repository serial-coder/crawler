/*
Copyright LLC Newity. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package blocklib

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

// Block contains all the necessary information about the blockchain block
type Block struct {
	Data       [][]byte
	Signatures []BlockSignature
	Metadata   [][]byte
	txsFilter  []uint8
}

// BlockSignature contains nonce, cert, MSP ID and signature of the orderer which signed the block
type BlockSignature struct {
	Cert      []byte // pem-encoded
	MSPID     string
	Signature []byte
	Nonce     []byte
}

// FromFabricBlock converts common.Block to blocklib.Block.
// Such conversion is necessary for further comfortable work with information from the block.
func FromFabricBlock(block *common.Block) (*Block, error) {
	metadata := &common.Metadata{}
	err := proto.Unmarshal(block.Metadata.Metadata[common.BlockMetadataIndex_SIGNATURES], metadata)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling metadata from block at index [%s]", common.BlockMetadataIndex_SIGNATURES)
	}

	var blockSignatures []BlockSignature
	for _, metadataSignature := range metadata.Signatures {
		sigHdr := &common.SignatureHeader{}
		if err := proto.Unmarshal(metadataSignature.SignatureHeader, sigHdr); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling SignatureHeader")
		}
		creator := &msp.SerializedIdentity{}
		if err = proto.Unmarshal(sigHdr.Creator, creator); err != nil {
			return nil, err
		}
		blockSignatures = append(blockSignatures,
			BlockSignature{
				Cert:      creator.IdBytes,
				MSPID:     creator.Mspid,
				Signature: metadataSignature.Signature,
				Nonce:     sigHdr.Nonce,
			},
		)
	}

	return &Block{
		Data:       block.Data.Data,
		Signatures: blockSignatures,
		Metadata:   block.Metadata.Metadata,
		txsFilter:  block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER],
	}, nil
}

func (b *Block) BlockSignatures() ([]BlockSignature, error) {
	metadata := &common.Metadata{}
	err := proto.Unmarshal(b.Metadata[common.BlockMetadataIndex_SIGNATURES], metadata)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling metadata from block at index [%s]", common.BlockMetadataIndex_SIGNATURES)
	}

	var blockSignatures []BlockSignature
	for _, metadataSignature := range metadata.Signatures {
		sigHdr := &common.SignatureHeader{}
		if err := proto.Unmarshal(metadataSignature.SignatureHeader, sigHdr); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling SignatureHeader")
		}
		creator := &msp.SerializedIdentity{}
		if err = proto.Unmarshal(sigHdr.Creator, creator); err != nil {
			return nil, err
		}
		blockSignatures = append(blockSignatures,
			BlockSignature{
				Cert:      creator.IdBytes,
				MSPID:     creator.Mspid,
				Signature: metadataSignature.Signature,
				Nonce:     sigHdr.Nonce,
			},
		)
	}
	return blockSignatures, nil
}

func (b *Block) Txs() ([]Tx, error) {
	var (
		validationCode   int32  = 255
		validationStatus string = "INVALID_OTHER_REASON"
	)

	var txs []Tx
	for txNumber, data := range b.Data {
		for _, code := range peer.TxValidationCode_value {
			if b.txsFilter[txNumber] == uint8(code) {
				validationCode = code
				validationStatus = peer.TxValidationCode_name[code]
			}
		}

		txs = append(txs, Tx{Data: data, validationCode: validationCode, validationStatus: validationStatus})
	}
	return txs, nil
}

// LastConfig returns last configuration block index for provided block.
func (b *Block) LastConfig() (uint64, error) {
	lastConfig := &common.LastConfig{}
	err := proto.Unmarshal(b.Metadata[common.BlockMetadataIndex_LAST_CONFIG], lastConfig)
	return lastConfig.Index, err
}

func GetTx(block *common.Block, txNumber int) *Tx {
	txsFilter := TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	var (
		validationCode   int32  = 255
		validationStatus string = "INVALID_OTHER_REASON"
	)

	for _, code := range peer.TxValidationCode_value {
		if txsFilter[txNumber] == uint8(code) {
			validationCode = code
			validationStatus = peer.TxValidationCode_name[code]
		}
	}

	return &Tx{Data: block.Data.Data[txNumber], validationCode: validationCode, validationStatus: validationStatus}
}
